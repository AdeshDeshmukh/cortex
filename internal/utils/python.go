package utils

import (
	"os/exec"
	"runtime"
)

func DetectPython() string {
	candidates := []string{"python3", "python"}
	
	if runtime.GOOS == "windows" {
		candidates = []string{"python", "python3", "py"}
	}
	
	for _, cmd := range candidates {
		if path, err := exec.LookPath(cmd); err == nil {
			return path
		}
	}
	
	return "python3"
}
