package art

import (
	"github.com/gorilla/mux"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
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
	r.HandleFunc("/api/v1/Img/{id}", Img)
	r.HandleFunc("/", index)

	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", httpRateLimiter.RateLimit(r))
	if err != nil {
		log.Fatal(err)
	}
}

func callGenerator(job Image) {
	if images == nil {
		images = map[string]Image{}
	}
	mutex.Lock()
	i := indexOf(job.Identifier, queue)
	if i > -1 {
		queue = append(queue[:i], queue[i+1:]...)
	}
	mutex.Unlock()
	if i > -1 {
		res := GenerateImage(job.Image, sourceDir, outDir, job.Width, job.Height, job.Shapes, job.ShapesStroke, job.Triangulate, job.TriangulateBefore, job.StrokeThickness, job.BlurAmount, job.Min, job.Max)
		mutex.Lock()
		job.FileName = res
		images[job.Identifier] = job
		mutex.Unlock()
	}
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
