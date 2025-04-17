package endpoints

import (
	"backend/models"
	"backend/splitter"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"

	"backend/config"
)

type UploadRequest struct {
	Content string `json:"content"` // base64-encoded file content
	Ext     string `json:"ext"`
}

type UploadResponse struct {
	Filename     string               `json:"filename"`
	ProductUnits []models.ProductUnit `json:"product_units"`
	TotalPrice   int64                `json:"total_price"`
}

// Uploader POST /api/upload_and_analyze
func Uploader(w http.ResponseWriter, r *http.Request) {
	req := &UploadRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteJSON(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	filename, err := uploadFile(req)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	totalPrice, productUnits, err := getOcrAns(filename)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	splitter.SaveInitialSplitData(filename, &productUnits)

	OK(w, UploadResponse{Filename: filename, ProductUnits: productUnits, TotalPrice: totalPrice})
}

func uploadFile(req *UploadRequest) (string, error) {
	if req.Content == "" {
		return "", fmt.Errorf("missing file content")
	}

	data, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		return "", fmt.Errorf("invalid base64 content")
	}

	if !isImage(req.Ext) {
		return "", fmt.Errorf("uploaded file is not a supported image type")
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width < 400 || height < 400 {
		return "", fmt.Errorf("image resolution too small: %dx%d", width, height)
	}

	if width > 2000 || height > 2000 {
		ratio := float64(2000) / math.Max(float64(width), float64(height))
		newWidth := int(float64(width) * ratio)
		newHeight := int(float64(height) * ratio)

		resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		switch format {
		case "jpeg":
			err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 90})
		case "png":
			err = png.Encode(&buf, resizedImg)
		case "bmp":
			err = bmp.Encode(&buf, resizedImg)
		default:
			return "", fmt.Errorf("unsupported image format: %s", format)
		}
		if err != nil {
			return "", fmt.Errorf("failed to encode resized image: %v", err)
		}
		data = buf.Bytes()
	}

	filename := generateFilename(config.GetConfig().GetKey(), time.Now(), req.Ext)

	if err = ioutil.WriteFile(config.GetConfig().DataSharedDir+filename, data, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func isImage(contentType string) bool {
	switch contentType {
	case "jpeg", "png", "bmp", "jpg":
		return true
	default:
		return false
	}
}

func generateFilename(secret []byte, t time.Time, ext string) string {
	timestamp := t.UnixNano()
	h := hmac.New(sha256.New, secret)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(timestamp))
	h.Write(b)
	hash := h.Sum(nil)
	encoded := hex.EncodeToString(hash)
	return fmt.Sprintf("%s.%s", encoded, ext)
}

func getOcrAns(filename string) (int64, []models.ProductUnit, error) {
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		return 0, nil, err
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(filename)); err != nil {
		return 0, nil, err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return 0, nil, err
	}
	fmt.Println()
	fmt.Println(string(buf[:n]))
	fmt.Println()

	lines := strings.Split(string(buf[:n]), "\n")
	result := make([]models.ProductUnit, 0, len(lines))
	var totalCost int64

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		fields := strings.Split(line, "|")
		if len(fields) != 3 {
			continue
		}
		quant, _ := strconv.Atoi(fields[1])
		price, _ := strconv.ParseInt(fields[2], 10, 64)
		totalCost += price
		result = append(result, models.ProductUnit{ID: len(result), Name: fields[0], Quantity: quant, Price: price})
	}

	return totalCost, result, nil
}
