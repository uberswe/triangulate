package main

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

var (
	totalCycleCount = 5000
	sourceDir       = "resources/source"
	outDir          = "resources/out"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		err = os.MkdirAll(sourceDir, 0744)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0744)
		if err != nil {
			panic(err)
		}
	}
}

func main() {

}

func generateImage(img image.Image) {
	var err error
	imgName := fmt.Sprintf("%d_%s.png", time.Now().UnixNano(), randStringRunes(10))

	if img == nil {
		img, err = loadRandomUnsplashImage(2000, 2000)
	}
	if err != nil {
		log.Panicln(err)
	}

	s := generator.Generate(img, generator.UserParams{
		StrokeRatio:              0.75,
		DestWidth:                img.Bounds().Size().X,
		DestHeight:               img.Bounds().Size().Y,
		InitialAlpha:             0.1,
		StrokeReduction:          0.002,
		AlphaIncrease:            0.06,
		StrokeInversionThreshold: 0.05,
		StrokeJitter:             int(0.01 * float64(img.Bounds().Size().X)),
		MinEdgeCount:             4,
		MaxEdgeCount:             4,
		RotationSeed:             0.5,
		RandomRotation:           true,
		Stroke:                   false,
	})

	rand.Seed(time.Now().Unix())

	for i := 0; i < totalCycleCount; i++ {
		s.Update()
	}

	err = saveOutput(s.Output(), fmt.Sprintf("%s/%s", outDir, imgName))
	if err != nil {
		log.Println(err)
		return
	}

	err = saveOutput(img, fmt.Sprintf("%s/%s", sourceDir, imgName))
	if err != nil {
		log.Println(err)
		return
	}
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

func randStringRunes(n int) string {
	letterRunes := []rune("bcdfghjlmnpqrstvwxz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
