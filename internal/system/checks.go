package system

import (
	"fmt"
	"os/exec"
)

// CheckDependencies verifies that required external tools are available in the PATH.
func CheckDependencies() error {
	requiredTools := []string{"ffmpeg", "magick", "pdftotext"}
	missingTools := []string{}

	for _, tool := range requiredTools {
		_, err := exec.LookPath(tool)
		if err != nil {
			// Special case for ImageMagick: older versions use 'convert'
			if tool == "magick" {
				if _, err := exec.LookPath("convert"); err == nil {
					continue
				}
			}
			missingTools = append(missingTools, tool)
		}
	}

	if len(missingTools) > 0 {
		return fmt.Errorf("missing required tools: %v. Please install them and ensure they are in your PATH", missingTools)
	}

	return nil
}
