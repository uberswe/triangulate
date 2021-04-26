package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func main() {
	// https://api.unsplash.com/photos/random
	// PkygelfQYXYPxyYvzEbj5CWs9keFFdbZqjaavRsbT78
	req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Client-ID PkygelfQYXYPxyYvzEbj5CWs9keFFdbZqjaavRsbT78")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	r := UnsplashRandomImageResponse{}
	json.NewDecoder(resp.Body).Decode(&r)

	log.Println(r.Height)
	log.Println(r.Width)
	log.Println(r.Urls.Full)
	log.Println(r.Urls.Thumb)
	log.Println(r.Description)
	for name, values := range resp.Header {
		for _, value := range values {
			log.Println(name, value)
		}
	}

}
