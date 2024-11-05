package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type ImageInfo struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Layers int    `json:"layers"`
}

type ChartResponse struct {
	Images []ImageInfo `json:"images"`
}

// ErrorResponse represents the structure of the error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// FetchImagesFromHelmChart runs `helm template` on the local chart path and
// Parses output to find images
func FetchImagesFromHelmChart(chartPath string) ([]string, error) {

	helmCmd := exec.Command("helm", "template", chartPath)
	var stdout, stderr bytes.Buffer
	helmCmd.Stdout = &stdout
	helmCmd.Stderr = &stderr

	if err := helmCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to render Helm chart: %v, stderr: %s", err, stderr.String())
	}

	images, err := ParseImagesFromTemplateOutput(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse images from output: %v", err)
	}

	return images, nil
}

// ParseImagesFromTemplateOutput parses images from the Helm template output to find image names
func ParseImagesFromTemplateOutput(output []byte) ([]string, error) {
	var images []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "image:") {
			parts := strings.Split(line, ": ")
			if len(parts) > 1 {
				image := strings.TrimSpace(parts[1])
				images = append(images, image)
			}
		}
	}
	return images, nil
}

// FetchImageDetails reads the output fully to make sure the image is pulled to fetch image details
// Such as number of layers, size and append image details to response
func FetchImageDetails(images []string) ([]ImageInfo, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	var imageDetails []ImageInfo
	for _, image := range images {
		image = strings.Trim(image, `"`)
		reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to pull image %s: %v", image, err)
		}

		if _, err := io.ReadAll(reader); err != nil {
			return nil, fmt.Errorf("failed to read image pull output for %s: %v", image, err)
		}
		defer reader.Close()

		imgInfo, _, err := cli.ImageInspectWithRaw(context.Background(), image)
		if err != nil {
			return nil, fmt.Errorf("failed to inspect image %s: %v", image, err)
		}

		numLayers := len(imgInfo.RootFS.Layers)

		imageDetails = append(imageDetails, ImageInfo{
			Name:   image,
			Size:   imgInfo.Size,
			Layers: numLayers,
		})
	}
	return imageDetails, nil
}

// JsonErrorResponse sends a JSON-formatted error response
func JsonErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Unable to encode error response", http.StatusInternalServerError)
	}
}

// DownloadHelmRepoHelper clones the helm repository  and returns path to helm chart
func DownloadHelmRepoHelper(repoURL, chartPath string) (*string, error) {

	tmpDir, err := os.MkdirTemp("", "helm-chart-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to clone repository: %v", err)
	}

	chartFullPath := filepath.Join(tmpDir, chartPath)
	return &chartFullPath, nil
}
