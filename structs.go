package art

import (
	"image"
	"time"
)

type GeneratePollResponse struct {
	Queue        int    `json:"queue"`
	Link         string `json:"link"`
	Identifier   string `json:"identifier"`
	RandomImage  bool   `json:"randomImage"`
	Thumbnail    string `json:"thumbnail"`
	Description  string `json:"description"`
	UserName     string `json:"user_name"`
	UserLocation string `json:"user_location"`
	UserLink     string `json:"user_link"`
}

type Image struct {
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
	Image                image.Image
	RandomImage          bool   `json:"randomImage"`
	Thumbnail            string `json:"thumbnail,omitempty"`
	Description          string `json:"description,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	UserLocation         string `json:"user_location,omitempty"`
	UserLink             string `json:"user_link,omitempty"`
}
