package art

import (
	"fmt"
	"github.com/esimov/triangle"
	"github.com/gorilla/mux"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	sourceDir = "resources/source"
	outDir    = "resources/out"
	queue     []string
	images    map[string]Image
	mutex     = &sync.Mutex{}
	jobChan   = make(chan Image, 999999)
)

func worker(jobChan <-chan Image) {
	for job := range jobChan {
		callGenerator(job)
	}
}

func Run() {
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

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./assets/build/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/api/v1/generate", generate)
	r.HandleFunc("/api/v1/generate/{id}", generatePoll)
	r.HandleFunc("/api/v1/img/{id}.png", Img)
	r.HandleFunc("/", index)

	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", httpRateLimiter.RateLimit(r))
	if err != nil {
		log.Fatal(err)
	}
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
	}
	mutex.Unlock()
	if i > -1 {
		p := &triangle.Processor{
			BlurRadius:      job.BlurAmount,
			SobelThreshold:  10,
			PointsThreshold: 20,
			MaxPoints:       2500,
			Wireframe:       0,
			Noise:           0,
			StrokeWidth:     1,
			Grayscale:       false,
		}
		tri := triangle.Image{Processor: *p}
		img := job.Image
		if job.Triangulate && job.TriangulateBefore {
			img, _, _, err = tri.Draw(img, nil, triangulate)
			if err != nil {
				log.Println(err)
				return
			}
		}
		if job.Shapes || (!job.Shapes && !job.Triangulate) {
			img = GenerateImage(img, job.Width, job.Height, job.ShapesStroke, job.StrokeThickness, job.BlurAmount, job.Min, job.Max)
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

		err = saveOutput(job.Image, fmt.Sprintf("%s/%s", sourceDir, imgName))
		if err != nil {
			log.Println(err)
			return
		}

		mutex.Lock()
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
