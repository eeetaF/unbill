package endpoints

import (
	"backend/config"
	"backend/models"
	"backend/splitter"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type productsPutterRequest struct {
	UpdatedProducts []models.ProductUnit `json:"products"`
}
type productsPutterResponse struct {
	ProductUnits []models.ProductUnit `json:"product_units"`
	TotalPrice   int64                `json:"total_price"`
}

// ProductsPutter PUT /api/products/{filename}
func ProductsPutter(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")

	req := &productsPutterRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteJSON(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	price, err := putProducts(filename, &req.UpdatedProducts)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	OK(w, productsPutterResponse{ProductUnits: req.UpdatedProducts, TotalPrice: price})
}

func putProducts(filename string, updatedProducts *[]models.ProductUnit) (int64, error) {
	filePath := filepath.Join(config.GetConfig().DataSplitDir, fmt.Sprintf("%s.json", filename))
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	file.Close()

	splitter.SaveInitialSplitData(filename, updatedProducts)

	var totalPrice int64
	for _, product := range *updatedProducts {
		totalPrice += product.Price
	}
	return totalPrice, nil
}
