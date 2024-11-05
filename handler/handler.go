package handler

import (
	"encoding/json"
	"fmt"
	"go-api-assignment/utils"
	"go-api-assignment/validators"
	"net/http"
)

type ChartPath struct {
	ChartURL string `json:"chart_url"`
}

// HelmImagesHandler handles request to read the helm chart images and return image details
func HelmImagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reqData ChartPath
	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		errorMsg := fmt.Sprintf("Error decoding request to struct: %v", err.Error())
		utils.JsonErrorResponse(w, errorMsg, http.StatusBadRequest)
		return
	}

	err = validators.ValidateRequestInput(reqData.ChartURL)
	if err != nil {
		utils.JsonErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	images, err := utils.FetchImagesFromHelmChart(reqData.ChartURL)
	if err != nil {
		utils.JsonErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return

	}

	imageDetails, err := utils.FetchImageDetails(images)
	if err != nil {
		utils.JsonErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utils.ChartResponse{Images: imageDetails}
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		fmt.Printf("error encoding JSON response: %v\n", encodeErr)
	}
}



// MethodNotAllowedHandler handles method not allowed statuses
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err:=fmt.Sprintf("method %v not allowed for this endpoint", r.Method)
	utils.JsonErrorResponse(w, err, http.StatusMethodNotAllowed)
	
}
