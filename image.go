package art

import (
	"fmt"
	"github.com/uberswe/art/generator"
	"image"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func GenerateImage(img image.Image, width int, height int, stroke bool, StrokeThickness int, blurAmount int, shapeMin int, shapeMax int) image.Image {
	var err error

	if img == nil {
		img, err = loadRandomUnsplashImage(width, height)
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	totalCycleCount := 500 * blurAmount

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

func loadRandomUnsplashImage(width, height int) (image.Image, error) {
	url := fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	return img, err
}
