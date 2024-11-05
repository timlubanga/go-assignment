package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonErrorResponse(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		statusCode   int
		expectedBody string
		expectedCode int
	}{
		{
			name:         "Test Not Found Error",
			message:      "Resource not found",
			statusCode:   http.StatusNotFound,
			expectedBody: `{"message":"Resource not found"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Test Bad Request Error",
			message:      "Invalid request parameters",
			statusCode:   http.StatusBadRequest,
			expectedBody: `{"message":"Invalid request parameters"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Test Internal Server Error",
			message:      "Unexpected error occurred",
			statusCode:   http.StatusInternalServerError,
			expectedBody: `{"message":"Unexpected error occurred"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			JsonErrorResponse(recorder, tt.message, tt.statusCode)
			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
		})
	}
}

func TestFetchImagesFromHelmChart(t *testing.T) {
	repoURL := "https://github.com/helm/examples.git"
	chartPath := "charts/hello-world"
	validChartPath, err := DownloadHelmRepoHelper(repoURL, chartPath)
	assert.Nil(t, err)
	tests := []struct {
		Name    string
		Input   string
		Wanterr bool
	}{

		{
			Name:    "Test can fetch images successfully",
			Input:   *validChartPath,
			Wanterr: false,
		},
		{
			Name:    "Test throws error with invalid chart patch",
			Input:   "/temp/wrong-path/",
			Wanterr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			images, err := FetchImagesFromHelmChart(tt.Input)
			if tt.Wanterr {
				assert.NotNil(t, err)
				assert.Nil(t, images)
				assert.Contains(t, err.Error(), "failed to render Helm chart")
			}

			if !tt.Wanterr {
				assert.NotNil(t, images)
				assert.Nil(t, err)
				assert.True(t, len(images) > 0)

				for _, value := range images {
					assert.Contains(t, value, strings.TrimSpace("nginx:1.16.0"))

				}

			}

		})
	}

}

func TestFetchImageDetails(t *testing.T) {
	type images []string
	tests := []struct {
		Name    string
		Input   images
		Wanterr bool
	}{

		{
			Name:    "Test can fetch image details successfully",
			Input:   images{"nginx:1.16.0"},
			Wanterr: false,
		},
		{
			Name:    "Test throws error with invalid images",
			Input:   images{"invalidimage:latest"},
			Wanterr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			imagesInfo, err := FetchImageDetails(tt.Input)
			if tt.Wanterr {
				assert.NotNil(t, err)
				assert.Nil(t, imagesInfo)
				assert.Contains(t, err.Error(), "failed to pull image")
			}

			if !tt.Wanterr {
				assert.NotNil(t, imagesInfo)
				assert.Nil(t, err)
				assert.True(t, len(imagesInfo) > 0)

				for _, image := range imagesInfo {
					assert.Contains(t, image.Name, "nginx:1.16.0")
					if image.Size <= 0 {
						t.Errorf("Expected image size to be greater than 0, but got %d", image.Size)
					}

					if image.Layers < 0 {
						t.Errorf("Expected image layers to be greater than or equal to 0, but got %d", image.Layers)
					}

				}

			}

		})
	}

}
