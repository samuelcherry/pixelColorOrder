package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"time"
)


type PixelMap struct {
	H,S,V float64
	Count int
}

func rgbToHSV(r,g,b uint8) (float64, float64, float64) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) /255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min


	var h float64
	switch {
	case delta == 0:
		h = 0
	case max == rf:
		h = math.Mod((gf-bf)/delta, 6)
	case max == gf:
		h = (bf-rf)/delta + 2
	case max == bf:
		h = (rf-gf)/delta + 4
	}

	h *= 60
	if h < 0 {
		h+= 360
	}

	var s float64
	if max == 0 {
		s = 0
	}else {
		s = delta/max
	}

	v := max

	return h,s,v
}

func importImage(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil, err
	} 

		return img, nil
}

func createHashMap(img image.Image) map[string]PixelMap {
	
	bounds := img.Bounds()
	pixelMap := make(map[string]PixelMap)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _:= img.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			h, s, v := rgbToHSV(r8,g8,b8)

			h = math.Ceil(h)
			s = math.Round(s*100)/100
			v = math.Round(v*100)/100

			key := fmt.Sprintf("%.0f|%.2f|%.2f ", h,s,v)
			if pm, ok := pixelMap[key]; ok {
				pm.Count++
				pixelMap[key] = pm
			}else{
				pixelMap[key] = PixelMap{H: h, S:s, V:v, Count:1}
			}
		}
	}

	return pixelMap
}


func writeToFile(pixelMap map[string]PixelMap, outputFile string) error {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	for key, value := range pixelMap {
		line := fmt.Sprintf("%s: %+v\n", key, value)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image>")
		return
	}

	fileName := os.Args[1]

	img, err := importImage(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	start := time.Now()
	pixelMap := createHashMap(img)
	fmt.Println("Processing took:", time.Since(start))


	err = writeToFile(pixelMap, "output.txt")
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}