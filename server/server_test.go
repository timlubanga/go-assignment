package main

import (
	"bytes"
	"encoding/json"
	"go-api-assignment/handler"
	"go-api-assignment/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandleHelmImages(t *testing.T) {

	router := mux.NewRouter()
	router.Path("/api/helm-images").Methods(http.MethodPost).HandlerFunc(handler.HelmImagesHandler)
	chartURL, _ := utils.DownloadHelmRepoHelper("https://github.com/helm/examples.git", "charts/hello-world")

	tests := []struct {
		name         string
		inputPayload handler.ChartPath
		expectedCode int
	}{
		{
			name: "Valid request",
			inputPayload: handler.ChartPath{
				ChartURL: *chartURL,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing ChartURL",
			inputPayload: handler.ChartPath{
				ChartURL: "",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid path",
			inputPayload: handler.ChartPath{
				ChartURL: "invalid-url",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.inputPayload)

			req := httptest.NewRequest(http.MethodPost, "/api/helm-images", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}
			var response utils.ChartResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("error sending response")
			}
		})
	}
}
