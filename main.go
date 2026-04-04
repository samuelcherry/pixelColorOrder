package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"os"
)


type PixelMap struct {
	R,G,B float64 `json:"h"`
	Count int	  `json:"count"`
}

func CORSmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("image")
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if filename == "" {
		http.Error(w, "Image parameter is required", http.StatusBadRequest)
		return
	}

	img, err := importImage(filename)
	if err != nil {
		http.Error(w, err.Error(),500)
		return
	}

	pixelMap := createHashMap(img)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pixelMap)
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

func roundTo25(value float64) float64 {
	return math.Round(value/25)*25
}

func createHashMap(img image.Image) map[string]PixelMap {
	
	bounds := img.Bounds()
	pixelMap := make(map[string]PixelMap)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _:= img.At(x, y).RGBA()

			r8 := roundTo25(float64(r >> 8))
			g8 := roundTo25(float64(g >> 8))
			b8 := roundTo25(float64(b >> 8))


			key := fmt.Sprintf("%.0f|%.2f|%.2f ", r8,g8,b8)
			if pm, ok := pixelMap[key]; ok {
				pm.Count++
				pixelMap[key] = pm
			}else{
				pixelMap[key] = PixelMap{R: r8, G:g8, B:b8, Count:1}
			}
		}
	}
	return pixelMap
}


func handleProcess(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("image")
	if filename == "" {
		http.Error(w, "Image parameter is required", http.StatusBadRequest)
		return
	}

	img, err := importImage(filename)
	if err != nil {
		http.Error(w, err.Error(),500)
		return
	}

	pixelMap := createHashMap(img)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pixelMap)

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/process", handleProcess)
	
	handler :=CORSmiddleware(mux)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)


}