package endpoints

import (
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
	"net/http"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"

	"backend/runtime-data"
)

type UploadRequest struct {
	Content string `json:"content"` // base64-encoded file content
	Ext     string `json:"ext"`
}

type UploadResponse struct {
	Filename string `json:"filename"`
}

// POST /upload
func Uploader(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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

	OK(w, UploadResponse{Filename: filename})
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

	filename := generateFilename(runtime_data.GetConfig().GetKey(), time.Now(), req.Ext)

	if err = ioutil.WriteFile(runtime_data.GetConfig().DataDir+filename, data, 0644); err != nil {
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
