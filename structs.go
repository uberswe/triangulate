package art

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
	EmailHash        string
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
	TriangulateGrayscale bool
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
	PriceId   string `json:"price_id"`
	StripeKey string `json:"stripe_key"`
}
