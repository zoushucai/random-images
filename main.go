package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ImageInfo struct {
	Sub  string `json:"sub"`
	File string `json:"file"`
}

var imagesData []ImageInfo
var imagesDir = "./images"

func loadImagesData() error {
	data, err := ioutil.ReadFile("images_info.json")
	if err != nil {
		return fmt.Errorf("failed to read images_info.json: %v", err)
	}
	if err := json.Unmarshal(data, &imagesData); err != nil {
		return fmt.Errorf("failed to parse images_info.json: %v", err)
	}
	return nil
}

func main() {
	port := 2113

	// 检查端口是否被占用
	if isPortInUse(port) {
		fmt.Printf("Port %d is already in use. Attempting to kill the process...\n", port)
		err := killProcessUsingPort(port)
		if err != nil {
			fmt.Printf("Failed to kill process using port %d: %v\n", port, err)
			return
		}
		fmt.Printf("Process using port %d killed successfully.\n", port)
	}

	// 启动应用程序
	if err := loadImagesData(); err != nil {
		log.Fatalf("Failed to load images data: %v", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(imagesDir))))
	r.HandleFunc("/random", handleRandomImage).Methods("GET")

	handler := cors.Default().Handler(r)
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func handleRandomImage(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	sub := params.Get("sub")
	widthStr := params.Get("width")
	imageType := params.Get("type")
	contains := params.Get("contains")
	indexStr := params.Get("index")
	device := params.Get("device")
	jsonStr := params.Get("json")

	width, err := strconv.Atoi(widthStr)
	if err != nil || width <= 0 {
		width = 1920
	}

	if imageType == "" {
		imageType = "jpeg"
	}

	jsonOutput, err := strconv.Atoi(jsonStr)
	if err != nil {
		jsonOutput = 0
	}

	var filteredImages []ImageInfo
	for _, image := range imagesData {
		if (sub == "" || image.Sub == sub) && (contains == "" || strings.Contains(image.File, contains)) {
			filteredImages = append(filteredImages, image)
		}
	}

	if len(filteredImages) == 0 {
		http.Error(w, "No images found.", http.StatusNotFound)
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(filteredImages) {
		rand.Seed(time.Now().UnixNano())
		index = rand.Intn(len(filteredImages))
	}

	selectedImage := filteredImages[index]
	imagePath := filepath.Join(imagesDir, selectedImage.Sub, selectedImage.File)

	img, err := imaging.Open(imagePath)
	if err != nil {
		http.Error(w, "Error opening image.", http.StatusInternalServerError)
		log.Printf("Error opening image '%s': %v", imagePath, err)
		return
	}

	img = imaging.Resize(img, width, 0, imaging.Lanczos)

	if device == "mobile" || device == "tablet" || device == "phone" {
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		cropSize_min := min(width, height)
		img = imaging.CropCenter(img, cropSize_min*height/width, cropSize_min)
	}

	var buffer []byte
	var contentType string

	switch imageType {
		case "jpeg", "jpg":
			buffer, err = encodeJPEG(img)
			contentType = "image/jpeg"
		case "png":
			buffer, err = encodePNG(img)
			contentType = "image/png"
		case "webp":
			buffer, err = webpEncode(img)
			contentType = "image/webp"
		default:
			ext := filepath.Ext(selectedImage.File)
			switch ext {
			case ".jpg", ".jpeg":
				buffer, err = encodeJPEG(img)
				contentType = "image/jpeg"
			case ".png":
				buffer, err = encodePNG(img)
				contentType = "image/png"
			case ".webp":
				buffer, err = webpEncode(img)
				contentType = "image/webp"
			default:
				buffer, err = encodeJPEG(img)
				contentType = "image/jpeg"
			}
	}

	if err != nil {
		http.Error(w, "Error encoding image.", http.StatusInternalServerError)
		log.Printf("Error encoding image: %v", err)
		return
	}

	if jsonOutput == 1 {
		imageData := map[string]interface{}{
			"width":    img.Bounds().Dx(),
			"height":   img.Bounds().Dy(),
			"type":     imageType,
			"sub":      selectedImage.Sub,
			"imageurl": selectedImage.File,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(imageData)
	} else {
		w.Header().Set("Content-Type", contentType)
		w.Write(buffer)
	}
}

func encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, fmt.Errorf("JPEG encoding error: %v", err)
	}
	return buf.Bytes(), nil
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("PNG encoding error: %v", err)
	}
	return buf.Bytes(), nil
}

func webpEncode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, &webp.Options{Quality: 90})
	if err != nil {
		return nil, fmt.Errorf("WebP encoding error: %v", err)
	}
	return buf.Bytes(), nil
}

// 检查端口是否被占用
func isPortInUse(port int) bool {
	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	listener.Close()
	return false
}

// 强制关闭占用端口的进程
func killProcessUsingPort(port int) error {
	out, err := exec.Command("lsof", "-i", ":"+strconv.Itoa(port)).Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			pid := fields[1]
			pidInt, err := strconv.Atoi(pid)
			if err != nil {
				continue
			}
			process, err := os.FindProcess(pidInt)
			if err != nil {
				continue
			}
			err = process.Signal(os.Kill)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}