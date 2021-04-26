package art

import "time"

type GenerateRequest struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ImageType string `json:"image_type"`
	Shapes    bool   `json:"shapes"`
}

type GeneratePollResponse struct {
	Queue      int    `json:"queue"`
	Link       string `json:"link"`
	Identifier string `json:"identifier"`
}

type Image struct {
	FileName   string    `json:"file_name"`
	Identifier string    `json:"identifier"`
	Timestamp  time.Time `json:"timestamp"`
	RequestIP  string    `json:"request_ip"`
	Width      int
	Height     int
	ImageType  string
	Shapes     bool
}
