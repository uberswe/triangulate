package art

import (
	"image"
	"time"
)

type GeneratePollResponse struct {
	Queue      int    `json:"queue"`
	Link       string `json:"link"`
	Identifier string `json:"identifier"`
}

type Image struct {
	FileName          string    `json:"file_name"`
	Identifier        string    `json:"identifier"`
	Timestamp         time.Time `json:"timestamp"`
	RequestIP         string    `json:"request_ip"`
	Width             int
	Height            int
	ImageType         string
	Shapes            bool
	Max               int
	Min               int
	BlurAmount        int
	ShapesStroke      bool
	StrokeThickness   int
	Triangulate       bool
	TriangulateBefore bool
	Image             image.Image
}
