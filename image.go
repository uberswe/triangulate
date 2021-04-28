package art

import (
	"fmt"
	"github.com/uberswe/art/generator"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func GenerateImage(img image.Image, src string, out string, width int, height int, shapes bool, stroke bool, triangulate bool, triangulateBefore bool, StrokeThickness int, blurAmount int, shapeMin int, shapeMax int) string {
	var err error
	imgName := fmt.Sprintf("%d_%s.png", time.Now().UnixNano(), RandStringRunes(10))

	if img == nil {
		img, err = loadRandomUnsplashImage(width, height)
		if err != nil {
			log.Println(err)
			return ""
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

	err = saveOutput(s.Output(), fmt.Sprintf("%s/%s", out, imgName))
	if err != nil {
		log.Println(err)
		return "#"
	}

	err = saveOutput(img, fmt.Sprintf("%s/%s", src, imgName))
	if err != nil {
		log.Println(err)
		return "#"
	}
	return imgName
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

func saveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func RandStringRunes(n int) string {
	letterRunes := []rune("bcdfghjlmnpqrstvwxz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
