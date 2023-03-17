package utils

import (
	"os/exec"
)

func Contains[T string | int | float64 | rune](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Execute(bin string, args ...string) (string, error) {
	cmd := exec.Command(bin, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
