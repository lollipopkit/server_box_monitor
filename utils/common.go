package utils

import (
	"bytes"
	"os/exec"
)

func Contains[T string|int|float64|rune](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Execute(bin string, args ...string) (string, string, error) {
	cmd := exec.Command(bin, args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    return string(stdout.Bytes()), string(stderr.Bytes()), err
}