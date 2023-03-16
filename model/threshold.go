package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lollipopkit/server_box_monitor/utils"
)

type Threshold struct {
	ThresholdType
	Value float64
	CompareType
}

func ParseToThreshold(s string) (*Threshold, error) {
	s = strings.ToLower(s)
	runes := []rune(s)

	var thresholdType ThresholdType
	var compareType CompareType
	var idx int

	if utils.Contains(runes, '%') {
		thresholdType = ThresholdTypePercent
	} else if utils.Contains(sizeSuffix, string(runes[len(runes)-1])) {
		if strings.HasSuffix(s, "/s") {
			thresholdType = ThresholdTypeSpeed
		} else {
			thresholdType = ThresholdTypeSize
		}
	}
	runes = runes[:len(runes)-1]
	if runes[0] == '<' {
		if runes[1] == '=' {
			compareType = CompareTypeLessOrEqual
			idx = 2
		} else {
			compareType = CompareTypeLess
			idx = 1
		}
	} else if runes[0] == '>' {
		if runes[1] == '=' {
			compareType = CompareTypeGreaterOrEqual
			idx = 2
		} else {
			compareType = CompareTypeGreater
			idx = 1
		}
	} else if runes[0] == '=' {
		compareType = CompareTypeEqual
		idx = 0
	}

	value, err := strconv.ParseFloat(string(runes[idx:]), 64)
	if err != nil {
		return nil, err
	}
	return &Threshold{
		ThresholdType: thresholdType,
		Value:         value,
		CompareType:   compareType,
	}, nil
}

func (t *Threshold) True(now any) (bool, error) {
	switch t.ThresholdType {
	case ThresholdTypePercent:
		now, ok := now.(float64)
		if !ok {
			return false, errors.Join(ErrInvalidRule, fmt.Errorf("%v is not float64", now))
		}
		switch t.CompareType {
		case CompareTypeLess:
			return now < t.Value, nil
		case CompareTypeLessOrEqual:
			return now <= t.Value, nil
		case CompareTypeEqual:
			return now == t.Value, nil
		case CompareTypeGreaterOrEqual:
			return now >= t.Value, nil
		case CompareTypeGreater:
			return now > t.Value, nil
		}
	case ThresholdTypeSize:
		now, ok := now.(Size)
		if !ok {
			return false, errors.Join(ErrInvalidRule, fmt.Errorf("%v is not Size", now))
		}
		nowFloat64 := float64(now)
		switch t.CompareType {
		case CompareTypeLess:
			return nowFloat64 < t.Value, nil
		case CompareTypeLessOrEqual:
			return nowFloat64 <= t.Value, nil
		case CompareTypeEqual:
			return nowFloat64 == t.Value, nil
		case CompareTypeGreaterOrEqual:
			return nowFloat64 >= t.Value, nil
		case CompareTypeGreater:
			return nowFloat64 > t.Value, nil
		}
	}
	return false, ErrCompareFailed
}

type CompareType uint8

const (
	CompareTypeLess           CompareType = iota // <
	CompareTypeLessOrEqual                       // <=
	CompareTypeEqual                             // =
	CompareTypeGreaterOrEqual                    // >=
	CompareTypeGreater                           // >
)

type ThresholdType uint8

const (
	ThresholdTypeUnknown ThresholdType = iota
	// eg: 80% 80.001%
	ThresholdTypePercent
	// eg: 100m 10k 1g
	// Values should be in lower case
	ThresholdTypeSize
	// eg: 10m/s
	ThresholdTypeSpeed
)
