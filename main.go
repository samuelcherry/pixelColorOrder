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


func importImage(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error opening file:", err)
	}
	return file, nil
}

func decodeImage(file *os.File)(image.Image, error){
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Error decoding image:", err)
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

			key := fmt.Sprintf("%.0f|%.0f|%.0f", r8,g8,b8)
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