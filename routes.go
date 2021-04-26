package art

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func index(w http.ResponseWriter, r *http.Request) {
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

	gr := GenerateRequest{}
	err := json.NewDecoder(r.Body).Decode(&gr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if gr.Width > 1200 || gr.Height > 1200 {
		http.Error(w, "max size is 1200x1200", 500)
		return
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	mutex.Lock()
	id := generateUniqueId(queue, 10)
	queue = append(queue, id)
	resp := GeneratePollResponse{
		Queue:      len(queue),
		Link:       "",
		Identifier: id,
	}
	job := Image{
		Identifier: id,
		Timestamp:  time.Now(),
		// TODO if we run this behind a load balancer the IP will be local so we have to adapt
		RequestIP: ip,
		Width:     gr.Width,
		Height:    gr.Height,
		ImageType: gr.ImageType,
		Shapes:    gr.Shapes,
	}
	jobChan <- job
	mutex.Unlock()
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func image(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res := ""
	mutex.Lock()
	if val, ok := images[vars["id"]]; ok {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && ip == val.RequestIP {
			res = fmt.Sprintf("%s/%s", outDir, val.FileName)
		} else {
			log.Println(errors.New(fmt.Sprintf("%s did not match %s", r.RemoteAddr, val.RequestIP)))
			http.Error(w, http.StatusText(500), 500)
			return
		}
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
