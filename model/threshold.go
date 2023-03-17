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
	runesLen := len(runes)

	var thresholdType ThresholdType
	var compareType CompareType
	var startIdx, endIdx int

	// 判断阈值类型
	if utils.Contains(runes, '%') {
		thresholdType = ThresholdTypePercent
		endIdx = runesLen - 1
	} else if utils.Contains(sizeSuffix, string(runes[len(runes)-1])) {
		if strings.HasSuffix(s, "/s") {
			thresholdType = ThresholdTypeSpeed
			// 10m/s -> m/s -> 3
			endIdx = runesLen - 3
		} else {
			// 10m -> m -> 1
			thresholdType = ThresholdTypeSize
			endIdx = runesLen - 1
		}
	} else if strings.HasSuffix(s, "c") {
		thresholdType = ThresholdTypeTemperature
		// 32c -> c -> 1
		endIdx = runesLen - 1
	}

	runes = runes[:len(runes)-1]
	if runes[0] == '<' {
		if runes[1] == '=' {
			compareType = CompareTypeLessOrEqual
			startIdx = 2
		} else {
			compareType = CompareTypeLess
			startIdx = 1
		}
	} else if runes[0] == '>' {
		if runes[1] == '=' {
			compareType = CompareTypeGreaterOrEqual
			startIdx = 2
		} else {
			compareType = CompareTypeGreater
			startIdx = 1
		}
	} else if runes[0] == '=' {
		compareType = CompareTypeEqual
		startIdx = 0
	}

	value, err := strconv.ParseFloat(string(runes[startIdx:endIdx]), 64)
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
	// eg: 32c
	ThresholdTypeTemperature
)
