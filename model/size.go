package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const programKilo float64 = 1024

var (
	sizeSuffix = []string{"b", "k", "m", "g", "t"}

	zeroSpeed = Speed{0, time.Second}
)

// Size is a type that represents a size in bytes.
type Size uint64

func (s Size) String() string {
	nth := 0
	temp := float64(s)
	for {
		if temp < programKilo || nth == len(sizeSuffix)-1 {
			return fmt.Sprintf("%.1f %s", temp, sizeSuffix[nth])
		}
		temp = temp / programKilo
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

type Speed struct {
	Size
	Time time.Duration
}
func (s *Speed) String() (string, error) {
	if s.Time == 0 {
		return "", fmt.Errorf("time equals zero: %#v", s)
	}
	return fmt.Sprintf("%s/s", Size(float64(s.Size)/s.Time.Seconds()).String()), nil
}
func (s *Speed) Compare(other *Speed) (int, error) {
	if s.Time == 0 || other.Time == 0 {
		return 0, fmt.Errorf("time equals zero: %#v, %#v", s, other)
	}
	return int(float64(s.Size)/s.Time.Seconds() - float64(other.Size)/other.Time.Seconds()), nil
}

func ParseToSpeed(s string) (*Speed, error) {
	s = strings.ToLower(s)
	if s == "0" {
		return &zeroSpeed, nil
	}
	splited := strings.Split(s, "/")
	if len(splited) != 2 {
		return nil, fmt.Errorf("invalid speed: %s", s)
	}

	size, err := ParseToSize(splited[0])
	if err != nil {
		return nil, err
	}
	time, err := time.ParseDuration(splited[1])
	if err != nil {
		return nil, err
	}
	return &Speed{size, time}, nil
}

