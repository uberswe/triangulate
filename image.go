package triangulate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/esimov/triangle"
	"github.com/uberswe/triangulate/generator"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func callGenerator(job Image) {
	var err error
	if images == nil {
		images = map[string]Image{}
	}
	imgName := fmt.Sprintf("%d_%s.png", time.Now().UnixNano(), RandStringRunes(10))
	mutex.Lock()
	i := indexOf(job.Identifier, queue)
	if i > -1 {
		queue = append(queue[:i], queue[i+1:]...)
		currentJob = job
	}
	mutex.Unlock()
	if i > -1 {
		log.Println("image generation started")
		wireframe := 0
		if job.TriangulateWireframe {
			wireframe = 1
		}
		noise := 0
		if job.TriangulateNoise {

		}
		p := &triangle.Processor{
			BlurRadius:      int(math.Round(float64(job.ComplexityAmount/10))) + 1,
			SobelThreshold:  job.SobelThreshold,
			PointsThreshold: job.PointsThreshold,
			MaxPoints:       job.MaxPoints,
			Wireframe:       wireframe,
			Noise:           noise,
			StrokeWidth:     float64(job.StrokeThickness),
			Grayscale:       job.TriangulateGrayscale,
		}
		tri := triangle.Image{Processor: *p}
		img := job.Image
		if img == nil {
			var source UnsplashRandomImageResponse
			img, source, err = loadRandomUnsplashImage(job.Width, job.Height)
			if err != nil {
				log.Println(err)
				return
			}
			if currentJob.Identifier == job.Identifier {
				mutex.Lock()
				currentJob.RandomImage = true
				currentJob.Thumbnail = source.Urls.Regular
				currentJob.Description = source.Description
				currentJob.UserName = source.User.Name
				currentJob.UserLocation = source.User.Location
				currentJob.UserLink = source.User.Links.HTML
				currentJob.ThumbnailLink = source.Links.HTML
				mutex.Unlock()
			}
		}

		if img != nil {
			err = saveOutput(img, fmt.Sprintf("%s/%s", sourceDir, imgName))
			if err != nil {
				log.Println(err)
				return
			}
		}

		if job.Triangulate && job.TriangulateBefore {
			img, _, _, err = tri.Draw(img, nil, triangulate)
			if err != nil {
				log.Println(err)
				return
			}
		}
		if job.Shapes || (!job.Shapes && !job.Triangulate) {
			img = GenerateImage(img, job.Width, job.Height, job.ShapesStroke, job.StrokeThickness, job.ComplexityAmount, job.Min, job.Max)
		}
		if job.Triangulate && !job.TriangulateBefore {
			img, _, _, err = tri.Draw(img, nil, triangulate)
			if err != nil {
				log.Println(err)
				return
			}
		}

		// Watermark for unauthenticated users
		if !job.AuthenticatedUser && img.Bounds().Max.X > 200 && img.Bounds().Max.Y > 200 {
			b := img.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
			addLabel(m, img.Bounds().Max.X-125, img.Bounds().Max.Y-10, "Triangulate.xyz")

			err = saveOutput(m, fmt.Sprintf("%s/%s", outDir, imgName))
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			err = saveOutput(img, fmt.Sprintf("%s/%s", outDir, imgName))
			if err != nil {
				log.Println(err)
				return
			}
		}

		log.Println("image generated")
		mutex.Lock()
		stat := Stat{}
		if res := db.First(&stat, "key = ?", "total_generated"); res.Error == nil {
			stat.Value = stat.Value + 1
			db.Save(&stat)
		}
		job.FileName = imgName
		images[job.Identifier] = job
		mutex.Unlock()
	}

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

		req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", unsplashAccessKey))
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

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
