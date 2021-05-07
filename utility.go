package triangulate

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(http.StatusText(500))
		return
	}
}

func base64Encode(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}

func mergeMaps(left, right map[string]string) map[string]string {
	for key, rightVal := range right {
		if _, present := left[key]; !present {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
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
