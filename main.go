package triangulate

import (
	"github.com/gorilla/mux"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func Run() {
	closeHandler()
	initSettings()
	initQueue()
	initDatabase()
	initSessions()
	httpRateLimiter := initRatelimiter()

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
	apiRouter.HandleFunc("/api/v1/portal", portal)
	apiRouter.HandleFunc("/api/v1/stripe/webhook", handleWebhook)
	apiRouter.HandleFunc("/api/v1/generate", generate)
	apiRouter.HandleFunc("/api/v1/generate/{id}", generatePoll)
	apiRouter.HandleFunc("/api/v1/img/{id}.png", img)
	apiRouter.Use(AuthMiddleware)
	apiRouter.Use(SensitiveHeadersMiddleware)

	r := mux.NewRouter()
	r.PathPrefix("/api/v1/").Handler(httpRateLimiter.RateLimit(apiRouter))
	r.PathPrefix("/static/").Handler(staticRouter)

	fs2 := http.FileServer(http.Dir("./assets/build"))
	r.Path("/robots.txt").Handler(fs2)
	r.Path("/asset-manifest.json").Handler(fs2)

	pages := mux.NewRouter()
	pages.PathPrefix("/").HandlerFunc(index)
	pages.Use(StripeSessionMiddleware)

	r.PathPrefix("/").Handler(pages)
	r.Use(GeneralHeadersMiddleware)

	log.Printf("Listening on %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
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
}
