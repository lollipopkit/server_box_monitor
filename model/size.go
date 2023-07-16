package model

import (
	"fmt"
	"strconv"
	"strings"
)

const programKilo float64 = 1024

var (
	sizeSuffix = []string{"b", "k", "m", "g", "t"}
)

// Size is a type that represents a size in bytes.
type Size uint64

func (s Size) String() string {
	nth := 0
	temp := float64(s)
	for {
		if temp < programKilo || nth == len(sizeSuffix)-1 {
			return fmt.Sprintf("%.1f%s", temp, sizeSuffix[nth])
		}
		temp = temp / programKilo
		nth++
	}
}
func ParseToSize(s string) (Size, error) {
	s = strings.ToLower(s)
	if s == "0" {
		return 0, nil
	}
	nth := 0
	for _, v := range sizeSuffix {
		if strings.Contains(s, v) {
			break
		}
		nth++
	}
	if nth == len(sizeSuffix) {
		return 0, fmt.Errorf("invalid size: %s", s)
	}
	temp, err := strconv.ParseFloat(strings.ReplaceAll(s, sizeSuffix[nth], ""), 64)
	if err != nil {
		return 0, err
	}
	for i := 0; i < nth; i++ {
		temp = temp * programKilo
	}
	return Size(temp), nil
}
