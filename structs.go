package triangulate

import (
	"gorm.io/gorm"
	"image"
	"time"
)

type GeneratePollResponse struct {
	Queue         int    `json:"queue"`
	Link          string `json:"link"`
	Identifier    string `json:"identifier"`
	RandomImage   bool   `json:"randomImage"`
	Thumbnail     string `json:"thumbnail"`
	Description   string `json:"description"`
	UserName      string `json:"user_name"`
	UserLocation  string `json:"user_location"`
	UserLink      string `json:"user_link"`
	ThumbnailLink string `json:"image_link"`
}

type User struct {
	gorm.Model
	EmailHash        string `gorm:"unique"`
	PasswordHash     string
	StripeCustomerID string
}

type Stat struct {
	gorm.Model
	Key   string
	Value int
}

type Image struct {
	gorm.Model
	FileName             string    `json:"file_name"`
	Identifier           string    `json:"identifier"`
	Timestamp            time.Time `json:"timestamp"`
	RequestIP            string    `json:"request_ip"`
	Width                int
	Height               int
	ImageType            string
	Shapes               bool
	Max                  int
	Min                  int
	ComplexityAmount     int
	ShapesStroke         bool
	StrokeThickness      int
	Triangulate          bool
	TriangulateBefore    bool
	MaxPoints            int
	PointsThreshold      int
	SobelThreshold       int
	TriangulateWireframe bool
	TriangulateNoise     bool
	TriangulateGrayscale bool        `json:"triangulate_grayscale"`
	Image                image.Image `gorm:"-"`
	RandomImage          bool        `json:"randomImage"`
	Thumbnail            string      `json:"thumbnail"`
	Description          string      `json:"description"`
	UserName             string      `json:"user_name"`
	UserLocation         string      `json:"user_location"`
	UserLink             string      `json:"user_link"`
	ThumbnailLink        string      `json:"image_link"`
}

type Settings struct {
	LoggedIn  bool   `json:"logged_in"`
	PriceId   string `json:"price_id"`
	StripeKey string `json:"stripe_key"`
}

type AuthSession struct {
	gorm.Model
	UserID        uint
	AuthSessionID string
}

type Session struct {
	TempSessionID   string
	StripeSessionID string
	AuthSessionID   string
}

type TempSession struct {
	gorm.Model
	SessionString   string `gorm:"unique"`
	StripeSessionID string
	Email           string
	Password        string
}

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
