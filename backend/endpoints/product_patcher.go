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
	"strconv"
)

type productsPatcherRequest struct {
	UpdatedProducts []models.ProductUnit `json:"products"`
}
type productsPatcherResponse struct {
	ProductUnits []models.ProductUnit `json:"product_units"`
	TotalPrice   int64                `json:"total_price"`
}

// ProductsPatcher PATCH /api/products/{filename}
func ProductsPatcher(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")

	req := &productsPatcherRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteJSON(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	products, price, err := patchProducts(filename, &req.UpdatedProducts)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	OK(w, productsPatcherResponse{ProductUnits: products, TotalPrice: price})
}

func patchProducts(filename string, updatedProducts *[]models.ProductUnit) ([]models.ProductUnit, int64, error) {
	filePath := filepath.Join(config.GetConfig().DataSplitDir, fmt.Sprintf("%s.json", filename))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var splitData splitter.SplitData
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&splitData); err != nil {
		return nil, 0, err
	}

	products := splitData.Products
	for _, updatedProduct := range *updatedProducts {
		// this is not the best implementation. Change data struct to make this operation in O(1) instead O(n)
		for i, product := range products {
			if product.ID == updatedProduct.ID {
				products[i] = updatedProduct
				break
			}
		}
	}
	return products
}
