package validators

import (
	"errors"
	"os"
	"strings"
)

// IsValidFolderPath checks to validate the string parsed is a valid folder path
func isValidFolderPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
// ValidateRequestInput validates the input decoded from the request body 
// to ensure its a valid chart patch
func ValidateRequestInput(chart_url string) error {

	newURL := strings.TrimSpace(chart_url)
	if newURL == "" {
		err := "chart path is required/cannot be blank"

		return errors.New(err)

	}

	if !isValidFolderPath(newURL) {
		err := "invalid chart Path provided"

		return errors.New(err)

	}

	return nil

}
