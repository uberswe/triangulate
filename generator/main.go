package generator

// The original code can be found at https://github.com/preslavrachev/generative-art-in-go
//
// The following code has been modified by Markus Tenghamn (https://github.com/uberswe)
//
// MIT License
//
// Copyright (c) 2021 Preslav Rachev
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software
// and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
// LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"embed"
	"fmt"
	"github.com/JoshVarga/svgparser"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type UserParams struct {
	StrokeRatio              float64
	DestWidth                int
	DestHeight               int
	InitialAlpha             float64
	StrokeReduction          float64
	AlphaIncrease            float64
	StrokeInversionThreshold float64
	StrokeJitter             int
	MinEdgeCount             int
	MaxEdgeCount             int
	RotationSeed             float64
	RandomRotation           bool
	Stroke                   bool
}

type Image struct {
	UserParams
	iconPaths         [][]Point
	icons             []svgparser.Element
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
}

type Point struct {
	X float64
	Y float64
}

func Generate(source image.Image, userParams UserParams, svgs embed.FS) *Image {
	s := &Image{
		UserParams: userParams,
	}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.StrokeRatio * float64(s.DestWidth)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()

	s.source = source
	s.dc = canvas

	if s.icons == nil {
		s.icons = []svgparser.Element{}
		files, err := svgs.ReadDir("svgs")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			fmt.Println(f.Name())
			in, err := svgs.Open("svgs/" + f.Name())
			if err != nil {
				log.Println(err)
				return s
			}
			icon, err := svgparser.Parse(in, false)
			if err != nil {
				log.Println(err)
				return s
			}
			s.icons = append(s.icons, *icon)
			err = in.Close()
			if err != nil {
				log.Println(err)
				return nil
			}
		}
	}

	if s.icons != nil {
		for i, ic := range s.icons {
			s.iconPaths = append(s.iconPaths, []Point{})
			for _, element := range ic.Children {
				if element.Name == "path" {
					path, exists := element.Attributes["d"]
					if exists {
						reg, err := regexp.Compile("-?\\d+")
						if err != nil {
							log.Println(err)
							break
						}
						matches := reg.FindAllStringSubmatch(path, -1)
						first := ""
						for index, match := range matches {
							if index == 1 {
								first = match[0]
							}
							if index%2 == 0 {
								var f float64
								var f2 float64
								f, err = strconv.ParseFloat(match[0], 64)
								if err != nil {
									log.Println(err)
									break
								}
								floatString := first
								if index+1 < len(matches) {
									floatString = matches[index+1][0]
								}
								f2, err = strconv.ParseFloat(floatString, 64)
								if err != nil {
									log.Println(err)
									break
								}
								s.iconPaths[i] = append(s.iconPaths[i], Point{
									X: f,
									Y: f2,
								})
							}
						}
						break
					}
				}
			}
		}
	}

	return s
}

func (s *Image) Update() {
	rand.Seed(int64(time.Now().Nanosecond()))
	n := rand.Intn(len(s.icons) + 1)
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := rgb255(s.source.At(int(rndX), int(rndY)))

	destX := rndX * float64(s.DestWidth) / float64(s.sourceWidth)
	destX += float64(randRange(s.StrokeJitter))
	destY := rndY * float64(s.DestHeight) / float64(s.sourceHeight)
	destY += float64(randRange(s.StrokeJitter))
	edges := s.MinEdgeCount + rand.Intn(s.MaxEdgeCount-s.MinEdgeCount+1)

	s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
	rotation := s.RotationSeed
	if s.RandomRotation {
		rotation = rotation + rand.Float64()
	}

	if n < len(s.icons) {
		for index, point := range s.iconPaths[n] {
			// TODO this is still too slow
			// TODO this always starts at 0 0 of the entire image, it needs to be offset by dest x and y
			if index == 0 {
				s.dc.MoveTo(point.X, point.Y)
			} else {
				s.dc.LineTo(point.X, point.Y)
			}
		}
	} else {
		s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rotation)
	}

	s.dc.FillPreserve()

	if s.strokeSize <= s.StrokeInversionThreshold*s.initialStrokeSize {
		if (r+g+b)/3 < 128 {
			s.dc.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
		} else {
			s.dc.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
		}
	}
	if s.Stroke {
		s.dc.Stroke()
	} else {
		s.dc.ClearPath()
	}

	s.strokeSize -= s.StrokeReduction * s.strokeSize
	s.InitialAlpha += s.AlphaIncrease

	log.Println("updated")
}

func (s *Image) Output() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 257), int(g0 / 257), int(b0 / 257)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}
