package art

import (
	"encoding/gob"
	"fmt"
	"github.com/esimov/triangle"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	store            *sessions.CookieStore
	sourceDir        = "resources/source"
	outDir           = "resources/out"
	queue            []string
	images           map[string]Image
	mutex            = &sync.Mutex{}
	jobChan          = make(chan Image, 999999)
	currentJob       Image
	db               *gorm.DB
	priceID          string
	stripePrivateKey string
	stripePublicKey  string
	successUrl       string
	cancelUrl        string
	returnURL        string
	webhookSecret    string
	cookieName       string
	sessionIDParam   string
)

func worker(jobChan <-chan Image) {
	for job := range jobChan {
		callGenerator(job)
	}
}

func Run() {
	closeHandler()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	stripePrivateKey = os.Getenv("STRIPE_PRIVATE_KEY")
	stripePublicKey = os.Getenv("STRIPE_PUBLIC_KEY")
	successUrl = os.Getenv("SUCCESS_URL")
	cancelUrl = os.Getenv("CANCEL_URL")
	returnURL = os.Getenv("RETURN_URL")
	webhookSecret = os.Getenv("WEBHOOK_SECRET")
	priceID = os.Getenv("PRICE_ID")
	cookieName = os.Getenv("COOKIE_NAME")
	sessionIDParam = os.Getenv("SESSION_ID_PARAM")

	stripe.Key = stripePrivateKey

	db, err = gorm.Open(sqlite.Open("triangulate.db"), &gorm.Config{})
	if err != nil {
		log.Println(err)
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Image{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&Stat{})
	if err != nil {
		log.Fatal(err)
	}

	stat := Stat{}
	if res := db.First(&stat, "key = ?", "total_generated"); res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Fatal(res.Error)
	}
	if stat.ID == 0 {
		stat.Key = "total_generated"
		stat.Value = 0
		res := db.Create(&stat)
		if res.Error != nil {
			log.Fatal(res.Error)
		}
	}

	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		HttpOnly: true,
	}

	gob.Register(User{})

	store, err := memstore.New(65536)
	if err != nil {
		log.Fatal(err)
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(100),
		MaxBurst: 20,
	}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{RemoteAddr: true},
	}

	staticRouter := mux.NewRouter()
	fs := http.FileServer(http.Dir("./assets/build/static"))
	staticRouter.PathPrefix("/").Handler(http.StripPrefix("/static/", fs))

	apiRouter := mux.NewRouter()
	apiRouter.HandleFunc("/api/v1/login", login)
	apiRouter.HandleFunc("/api/v1/logout", logout)
	apiRouter.HandleFunc("/api/v1/register", register)
	apiRouter.HandleFunc("/api/v1/forgot-password", forgotPassword)
	apiRouter.HandleFunc("/api/v1/reset-password/{code}", resetPassword)
	apiRouter.HandleFunc("/api/v1/settings", settings)
	apiRouter.HandleFunc("/api/v1/stripe/webhook", handleWebhook)
	apiRouter.HandleFunc("/api/v1/auth/settings", authSettings)
	apiRouter.HandleFunc("/api/v1/auth/generate", generate)
	apiRouter.HandleFunc("/api/v1/auth/generate/{id}", generatePoll)
	apiRouter.HandleFunc("/api/v1/auth/img/{id}.png", img)
	apiRouter.HandleFunc("/api/v1/generate", generate)
	apiRouter.HandleFunc("/api/v1/generate/{id}", generatePoll)
	apiRouter.HandleFunc("/api/v1/img/{id}.png", img)
	// asset-manifest.json
	// robots.txt

	r := mux.NewRouter()
	r.PathPrefix("/api/v1/").Handler(httpRateLimiter.RateLimit(apiRouter))
	r.PathPrefix("/static/").Handler(staticRouter)

	fs2 := http.FileServer(http.Dir("./assets/build"))
	r.Path("/robots.txt").Handler(fs2)
	r.Path("/asset-manifest.json").Handler(fs2)
	r.PathPrefix("/").HandlerFunc(index)

	// TODO add
	// X-Content-Type-Options: nosniff
	// X-XSS-Protection: 1; mode=block
	// Strict-Transport-Security: max-age=<seconds>[; includeSubDomains]
	// Cache-control: no-store
	// Pragma: no-cache
	// X-Frame-Options: DENY
	// generate does not appear to contain an anti-CSRF token

	log.Printf("Listening on %s\n", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

func closeHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Terminating via terminal")
		d, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		err = d.Close()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func callGenerator(job Image) {
	var err error
	if images == nil {
		images = map[string]Image{}
	}
	imgName := fmt.Sprintf("%d_%s.png", time.Now().UnixNano(), RandStringRunes(10))
	mutex.Lock()
	i := indexOf(job.Identifier, queue)
	if i > -1 {
		queue = append(queue[:i], queue[i+1:]...)
		currentJob = job
	}
	mutex.Unlock()
	if i > -1 {
		log.Println("image generation started")
		wireframe := 0
		if job.TriangulateWireframe {
			wireframe = 1
		}
		noise := 0
		if job.TriangulateNoise {

		}
		p := &triangle.Processor{
			BlurRadius:      int(math.Round(float64(job.ComplexityAmount/10))) + 1,
			SobelThreshold:  job.SobelThreshold,
			PointsThreshold: job.PointsThreshold,
			MaxPoints:       job.MaxPoints,
			Wireframe:       wireframe,
			Noise:           noise,
			StrokeWidth:     float64(job.StrokeThickness),
			Grayscale:       job.TriangulateGrayscale,
		}
		tri := triangle.Image{Processor: *p}
		img := job.Image
		if img == nil {
			var source UnsplashRandomImageResponse
			img, source, err = loadRandomUnsplashImage(job.Width, job.Height)
			if err != nil {
				log.Println(err)
				return
			}
			if currentJob.Identifier == job.Identifier {
				mutex.Lock()
				currentJob.RandomImage = true
				currentJob.Thumbnail = source.Urls.Thumb
				currentJob.Description = source.Description
				currentJob.UserName = source.User.Name
				currentJob.UserLocation = source.User.Location
				currentJob.UserLink = source.User.Links.HTML
				currentJob.ThumbnailLink = source.Links.HTML
				mutex.Unlock()
			}
		}

		if img != nil {
			err = saveOutput(img, fmt.Sprintf("%s/%s", sourceDir, imgName))
			if err != nil {
				log.Println(err)
				return
			}
		}

		if job.Triangulate && job.TriangulateBefore {
			img, _, _, err = tri.Draw(img, nil, triangulate)
			if err != nil {
				log.Println(err)
				return
			}
		}
		if job.Shapes || (!job.Shapes && !job.Triangulate) {
			img = GenerateImage(img, job.Width, job.Height, job.ShapesStroke, job.StrokeThickness, job.ComplexityAmount, job.Min, job.Max)
		}
		if job.Triangulate && !job.TriangulateBefore {
			img, _, _, err = tri.Draw(img, nil, triangulate)
			if err != nil {
				log.Println(err)
				return
			}
		}

		if img.Bounds().Max.X > 200 && img.Bounds().Max.Y > 200 {
			b := img.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
			addLabel(m, img.Bounds().Max.X-125, img.Bounds().Max.Y-5, "Triangulate.xyz")

			err = saveOutput(m, fmt.Sprintf("%s/%s", outDir, imgName))
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			err = saveOutput(img, fmt.Sprintf("%s/%s", outDir, imgName))
			if err != nil {
				log.Println(err)
				return
			}
		}

		log.Println("image generated")
		mutex.Lock()
		stat := Stat{}
		if res := db.First(&stat, "key = ?", "total_generated"); res.Error == nil {
			stat.Value = stat.Value + 1
			db.Save(&stat)
		}
		job.FileName = imgName
		images[job.Identifier] = job
		mutex.Unlock()
	}

}

func triangulate() {
	log.Println("triangulate done")
}

func saveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		err = os.MkdirAll(sourceDir, 0744)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0744)
		if err != nil {
			panic(err)
		}
	}
	go worker(jobChan)
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

func generateUniqueId(data []string, len int) string {
	id := RandStringRunes(len)
	if indexOf(id, data) == -1 {
		return id
	}
	return generateUniqueId(data, len+1)
}

func RandStringRunes(n int) string {
	letterRunes := []rune("bcdfghjlmnpqrstvwxz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
