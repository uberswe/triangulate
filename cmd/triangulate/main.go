package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/uberswe/art"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type GeneratePollResponse struct {
	Queue      int    `json:"queue"`
	Link       string `json:"link"`
	Identifier string `json:"identifier"`
}

var (
	sourceDir = "resources/source"
	outDir    = "resources/out"
	queue     []string
	images    map[string]string
	mutex     = &sync.Mutex{}
)

func main() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./assets/build/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/api/v1/generate", generate)
	r.HandleFunc("/api/v1/generate/{id}", generatePoll)
	r.HandleFunc("/api/v1/image/{id}", image)
	r.HandleFunc("/", serveTemplate)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	indexFile := filepath.Join("assets", "build", "index.html")

	tmpl, err := template.New("").ParseFiles(indexFile)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// check your err
	err = tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func generate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	id := generateUniqueId(queue, 10)
	queue = append(queue, id)
	resp := GeneratePollResponse{
		Queue:      len(queue),
		Link:       "",
		Identifier: id,
	}
	go callGenerator(id)
	mutex.Unlock()
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func callGenerator(id string) {
	if images == nil {
		images = map[string]string{}
	}
	res := art.GenerateImage(nil, sourceDir, outDir)
	mutex.Lock()
	images[id] = res
	i := indexOf(id, queue)
	if i > -1 {
		queue = append(queue[:i], queue[i+1:]...)
	}
	mutex.Unlock()
}

func image(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res := ""
	mutex.Lock()
	if val, ok := images[vars["id"]]; ok {
		res = fmt.Sprintf("%s/%s", outDir, val)
	}
	mutex.Unlock()
	img, err := os.Open(res)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, img)
}

func generatePoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	res := ""
	id := vars["id"]
	i := indexOf(id, queue)
	if i == -1 {
		if _, ok := images[id]; ok {
			res = fmt.Sprintf("/api/v1/image/%s", id)
		}
	}
	resp := GeneratePollResponse{
		Queue: i + 1,
		Link:  res,
	}
	mutex.Unlock()
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
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
	id := art.RandStringRunes(len)
	if indexOf(id, data) == -1 {
		return id
	}
	return generateUniqueId(data, len+1)
}
