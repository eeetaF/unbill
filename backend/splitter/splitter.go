package splitter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"backend/config"
	"backend/models"
)

type Person struct {
	ID                 int
	Name               string
	AssignedProductIDs []int
}
type SplitData struct {
	Products []models.ProductUnit `json:"products"`
	Persons  []Person             `json:"persons"`
}

func EvenSplit(filename string, numPeople int) (int64, error) {
	filePath := filepath.Join(config.GetConfig().DataSplitDir, fmt.Sprintf("%s.json", filename))
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var splitData SplitData
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&splitData); err != nil {
		return 0, err
	}
	var totalPrice int64
	for i := 0; i < len(splitData.Products); i++ {
		totalPrice += splitData.Products[i].Price
	}
	return totalPrice / int64(numPeople), nil
}

func SaveInitialSplitData(filename string, products *[]models.ProductUnit) {
	splitData := SplitData{
		Products: *products,
	}

	filePath := filepath.Join(config.GetConfig().DataSplitDir, fmt.Sprintf("%s.json", filename))
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("SaveInitialSplitData: Failed to create file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&splitData); err != nil {
		log.Println("SaveInitialSplitData: Failed to create file:", err)
	}
}
