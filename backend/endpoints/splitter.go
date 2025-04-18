package endpoints

import (
	"backend/splitter"
	"net/http"
	"strconv"
)

type evenSplitterResponse struct {
	TotalPerPerson int64 `json:"total_per_person"`
}

// EvenSplitter GET /api/split_equally/{filename}/{num_people}
func EvenSplitter(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	numPeople, _ := strconv.Atoi(r.PathValue("num_people"))
	if numPeople == 0 {
		WriteJSON(w, http.StatusBadRequest, "Invalid num_people")
		return
	}

	totalPerPerson, err := splitter.EvenSplit(filename, numPeople)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	OK(w, evenSplitterResponse{TotalPerPerson: totalPerPerson})
}
