package triangulate

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func generate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO refactor this function

	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	loggedIn := false
	user := r.Context().Value(ContextUserKey)
	if user != nil && user.(uint) > 0 {
		loggedIn = true
	}

	var img image.Image

	imageType := r.FormValue("type")
	if imageType == "upload" {

		uploaded, _, err := r.FormFile("fileUpload")
		if err != nil {
			log.Println(err)
			writeJSONError(w, "", http.StatusInternalServerError)
			return
		}
		defer uploaded.Close()
		buffer := make([]byte, 512)
		_, err = uploaded.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		_, err = uploaded.Seek(0, 0)
		if err != nil {
			fmt.Println(err)
		}

		contentType := http.DetectContentType(buffer)

		if contentType != "image/jpeg" && contentType != "image/png" {
			log.Println(contentType)
			writeJSONError(w, "", http.StatusUnprocessableEntity)
			return
		}
		if contentType == "image/jpeg" {
			img, err = jpeg.Decode(uploaded)
			if err != nil {
				log.Println(err)
				writeJSONError(w, "", http.StatusInternalServerError)
				return
			}
		}

		if contentType == "image/png" {
			img, err = png.Decode(uploaded)
			if err != nil {
				log.Println(err)
				writeJSONError(w, "", http.StatusInternalServerError)
				return
			}
		}

	}

	width := r.FormValue("width")
	height := r.FormValue("height")
	shapes := r.FormValue("shapes")
	shapeStroke := r.FormValue("shapeStroke")
	triangulate := r.FormValue("triangulate")
	triangulateBefore := r.FormValue("triangulateBefore")
	strokeThickness := r.FormValue("strokeThickness")
	complexityAmount := r.FormValue("complexityAmount")
	min := r.FormValue("min")
	max := r.FormValue("max")
	maxPoints := r.FormValue("maxPoints")
	pointsThreshold := r.FormValue("pointsThreshold")
	sobelThreshold := r.FormValue("sobelThreshold")
	triangulateWireframe := r.FormValue("triangulateWireframe")
	triangulateNoise := r.FormValue("triangulateNoise")
	triangulateGrayscale := r.FormValue("triangulateGrayscale")
	text := r.FormValue("text")

	wi, err := strconv.Atoi(width)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	hi, err := strconv.Atoi(height)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}
	// Hash the IP
	hash := sha256.Sum256([]byte(ip))
	ip = fmt.Sprintf("%x", hash)

	shapesBool, err := strconv.ParseBool(shapes)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	shapesStrokeBool, err := strconv.ParseBool(shapeStroke)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	triangulateBool, err := strconv.ParseBool(triangulate)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	triangulateBeforeBool, err := strconv.ParseBool(triangulateBefore)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	triangulateNoiseBool, err := strconv.ParseBool(triangulateNoise)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	triangulateWireframeBool, err := strconv.ParseBool(triangulateWireframe)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	triangulateGrayscaleBool, err := strconv.ParseBool(triangulateGrayscale)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	strokeThicknessInt, err := strconv.Atoi(strokeThickness)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	complexityAmountInt, err := strconv.Atoi(complexityAmount)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	minInt, err := strconv.Atoi(min)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	maxInt, err := strconv.Atoi(max)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	maxPointsInt, err := strconv.Atoi(maxPoints)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	pointsThresholdInt, err := strconv.Atoi(pointsThreshold)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	sobelThresholdInt, err := strconv.Atoi(sobelThreshold)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	maxSize := 2000
	if loggedIn {
		maxSize = 10000
	}

	if wi > maxSize || hi > maxSize || wi < 0 || hi < 0 {
		writeJSONError(w, fmt.Sprintf("max size is %dx%d", maxSize, maxSize), http.StatusInternalServerError)
		return
	}

	if minInt > 10 || minInt < 3 {
		writeJSONError(w, "min vertices is invalid", http.StatusInternalServerError)
		return
	}

	if maxInt > 10 || maxInt < 3 {
		writeJSONError(w, "max vertices is invalid", http.StatusInternalServerError)
		return
	}

	if complexityAmountInt > 100 || complexityAmountInt < 1 {
		writeJSONError(w, "complexity is invalid", http.StatusInternalServerError)
		return
	}

	if strokeThicknessInt > 10 || strokeThicknessInt < 1 {
		log.Println(strokeThickness)
		writeJSONError(w, "stroke is invalid", http.StatusInternalServerError)
		return
	}

	if maxPointsInt > 5000 || maxPointsInt < 500 {
		log.Println(maxPoints)
		writeJSONError(w, "max points is invalid", http.StatusInternalServerError)
		return
	}

	if pointsThresholdInt > 30 || pointsThresholdInt < 10 {
		log.Println(pointsThreshold)
		writeJSONError(w, "point threshold is invalid", http.StatusInternalServerError)
		return
	}

	if sobelThresholdInt > 20 || sobelThresholdInt < 5 {
		log.Println(sobelThresholdInt)
		writeJSONError(w, "sobel threshold is invalid", http.StatusInternalServerError)
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
		RequestIP:            ip,
		Width:                wi,
		Height:               hi,
		ImageType:            imageType,
		Shapes:               shapesBool,
		Max:                  maxInt,
		Min:                  minInt,
		ComplexityAmount:     complexityAmountInt,
		StrokeThickness:      strokeThicknessInt,
		Triangulate:          triangulateBool,
		TriangulateBefore:    triangulateBeforeBool,
		ShapesStroke:         shapesStrokeBool,
		Image:                img,
		MaxPoints:            maxPointsInt,
		SobelThreshold:       sobelThresholdInt,
		PointsThreshold:      pointsThresholdInt,
		TriangulateWireframe: triangulateWireframeBool,
		TriangulateGrayscale: triangulateGrayscaleBool,
		TriangulateNoise:     triangulateNoiseBool,
		Text:                 text,
		AuthenticatedUser:    loggedIn,
	}
	if loggedIn {
		premiumJobChan <- job
	} else {
		jobChan <- job
	}
	mutex.Unlock()
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}
}

func img(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res := ""
	mutex.Lock()
	if val, ok := images[vars["id"]]; ok {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && ip == val.RequestIP {
			res = fmt.Sprintf("%s/%s", outDir, val.FileName)
		} else {
			log.Println(errors.New(fmt.Sprintf("%s did not match %s", r.RemoteAddr, val.RequestIP)))
			writeJSONError(w, "", http.StatusInternalServerError)
			return
		}
	}
	mutex.Unlock()
	img, err := os.Open(res)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/png")
	_, _ = io.Copy(w, img)
}

func generatePoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	res := ""
	id := vars["id"]
	i := indexOf(id, queue)
	resp := GeneratePollResponse{
		Queue: i + 1,
	}
	if i == -1 {
		if _, ok := images[id]; ok {
			res = fmt.Sprintf("/api/v1/img/%s.png", id)
			resp.Link = res
		}
		if currentJob.Identifier == id && currentJob.RandomImage {
			resp.Thumbnail = currentJob.Thumbnail
			resp.Description = currentJob.Description
			resp.RandomImage = currentJob.RandomImage
			resp.UserLocation = currentJob.UserLocation
			resp.UserName = currentJob.UserName
			resp.UserLink = currentJob.UserLink
			resp.ThumbnailLink = currentJob.ThumbnailLink
		}
	}
	mutex.Unlock()
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}
}
