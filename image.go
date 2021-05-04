package art

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uberswe/art/generator"
	"image"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var randomUnsplashImages []UnsplashRandomImageResponse

type UnsplashRandomImageResponse struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Color       string `json:"color"`
	BlurHash    string `json:"blur_hash"`
	Downloads   int    `json:"downloads"`
	Likes       int    `json:"likes"`
	LikedByUser bool   `json:"liked_by_user"`
	Description string `json:"description"`
	Exif        struct {
		Make         string `json:"make"`
		Model        string `json:"model"`
		ExposureTime string `json:"exposure_time"`
		Aperture     string `json:"aperture"`
		FocalLength  string `json:"focal_length"`
		Iso          int    `json:"iso"`
	} `json:"exif"`
	Location struct {
		Name     string `json:"name"`
		City     string `json:"city"`
		Country  string `json:"country"`
		Position struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"position"`
	} `json:"location"`
	CurrentUserCollections []struct {
		ID              int         `json:"id"`
		Title           string      `json:"title"`
		PublishedAt     string      `json:"published_at"`
		LastCollectedAt string      `json:"last_collected_at"`
		UpdatedAt       string      `json:"updated_at"`
		CoverPhoto      interface{} `json:"cover_photo"`
		User            interface{} `json:"user"`
	} `json:"current_user_collections"`
	Urls struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
	} `json:"urls"`
	Links struct {
		Self             string `json:"self"`
		HTML             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	User struct {
		ID                string `json:"id"`
		UpdatedAt         string `json:"updated_at"`
		Username          string `json:"username"`
		Name              string `json:"name"`
		PortfolioURL      string `json:"portfolio_url"`
		Bio               string `json:"bio"`
		Location          string `json:"location"`
		TotalLikes        int    `json:"total_likes"`
		TotalPhotos       int    `json:"total_photos"`
		TotalCollections  int    `json:"total_collections"`
		InstagramUsername string `json:"instagram_username"`
		TwitterUsername   string `json:"twitter_username"`
		Links             struct {
			Self      string `json:"self"`
			HTML      string `json:"html"`
			Photos    string `json:"photos"`
			Likes     string `json:"likes"`
			Portfolio string `json:"portfolio"`
		} `json:"links"`
	} `json:"user"`
}

func GenerateImage(img image.Image, width int, height int, stroke bool, StrokeThickness int, blurAmount int, shapeMin int, shapeMax int) image.Image {

	totalCycleCount := 20 * ((blurAmount * blurAmount) * 5)

	s := generator.Generate(img, generator.UserParams{
		StrokeRatio:              0.25 * float64(StrokeThickness),
		DestWidth:                width,
		DestHeight:               height,
		InitialAlpha:             0.1,
		StrokeReduction:          0.002,
		AlphaIncrease:            0.06,
		StrokeInversionThreshold: 0.05,
		StrokeJitter:             int(0.1 * float64(img.Bounds().Size().X)),
		MinEdgeCount:             shapeMin,
		MaxEdgeCount:             shapeMax,
		RotationSeed:             0.45,
		RandomRotation:           true,
		Stroke:                   stroke,
	})

	rand.Seed(time.Now().Unix())

	for i := 0; i < totalCycleCount; i++ {
		s.Update()
	}
	return s.Output()
}

func loadRandomUnsplashImage(width int, height int) (image.Image, UnsplashRandomImageResponse, error) {
	queueImages()
	var url UnsplashRandomImageResponse
	if len(randomUnsplashImages) > 0 {
		url, randomUnsplashImages = randomUnsplashImages[len(randomUnsplashImages)-1], randomUnsplashImages[:len(randomUnsplashImages)-1]

		res, err := http.Get(fmt.Sprintf("%s&w=%d&h=%d", url.Urls.Full, width, height))
		if err != nil {
			return nil, url, err
		}
		defer res.Body.Close()

		img, _, err := image.Decode(res.Body)
		return img, url, err
	}
	return nil, url, errors.New("error with image queue")
}

func queueImages() {
	if len(randomUnsplashImages) <= 0 {
		req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random?count=30", nil)
		if err != nil {
			log.Println(err)
			return
		}

		req.Header.Set("Authorization", "Client-ID PkygelfQYXYPxyYvzEbj5CWs9keFFdbZqjaavRsbT78")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		var imgs []UnsplashRandomImageResponse
		err = json.NewDecoder(resp.Body).Decode(&imgs)
		if err != nil {
			log.Println(err)
			return
		}
		for _, r := range imgs {
			randomUnsplashImages = append(randomUnsplashImages, r)
		}
	}
}
