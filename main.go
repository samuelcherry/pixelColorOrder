package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"sort"
)

type Pixel struct {
	R, G, B uint8
	H, S, V float64
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

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image>")
		return
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	} 

	bounds := img.Bounds()

	var pixels []Pixel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _:= img.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			h, s, v := rgbToHSV(r8,g8,b8)

			pixels = append(pixels, Pixel{
				R: r8,
				G: g8,
				B: b8,
				H: h,
				S: s,
				V: v,
			})
		}
	}

	fmt.Println("Total pixels:", len(pixels))

	sort.Slice(pixels, func(i,j int) bool {

		if pixels[i].V < 0.05 && pixels[j].V >= 0.05 {
			return true
		}
		if pixels[j].V < 0.05 && pixels[i].V >= 0.05 {
			return false
		}

		if pixels[i].S < 0.1 && pixels[j].S >= 0.1 {
			return true
		}
		if pixels[j].S < 0.1 && pixels[i].S >= 0.1 {
			return false
		}

		if pixels[i].H != pixels[j].H{
			return pixels[i].H < pixels[j].H
		}

		if pixels[i].S != pixels[j].S{
			return pixels[i].S < pixels[j].S
		}

		return pixels[i].V < pixels[j].V
	})

	outImg := image.NewRGBA(bounds)

	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p:= pixels[i]

			outImg.Set(x,y, color.RGBA{
				R: p.R,
				G: p.G,
				B: p.B,
				A: 255,
			})
			i++
		}
	}

	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	png.Encode(outFile, outImg)
	fmt.Println("Saved output.png")
}