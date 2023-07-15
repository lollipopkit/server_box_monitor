package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lollipopkit/gommon/util"
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

	if strings.Contains(s, "%") { // 10%
		thresholdType = ThresholdTypePercent
		endIdx = runesLen - 1

	} else if strings.HasSuffix(s, "/s") { // 10m/s
		thresholdType = ThresholdTypeSpeed
		// 10m/s -> "/s" -> 2
		endIdx = runesLen - 2

	} else if util.Contains(sizeSuffix, string(runes[runesLen-1])) { // 10m
		// 10m -> "" -> 0
		thresholdType = ThresholdTypeSize
		endIdx = runesLen

	} else if strings.HasSuffix(s, "c") { // 32c
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

	var value float64
	var err error

	switch thresholdType {
	case ThresholdTypeSize, ThresholdTypeSpeed:
		size, err := ParseToSize(string(runes[startIdx:endIdx]))
		if err != nil {
			return nil, err
		}
		value = float64(size)
	default:
		value, err = strconv.ParseFloat(string(runes[startIdx:endIdx]), 64)
		if err != nil {
			return nil, err
		}
	}

	return &Threshold{
		ThresholdType: thresholdType,
		Value:         value,
		CompareType:   compareType,
	}, nil
}

func (t *Threshold) True(now any) (bool, error) {
	switch t.ThresholdType {
	case ThresholdTypePercent, ThresholdTypeTemperature:
		var nowValue float64
		switch now.(type) {
		case float64:
			nowValue = now.(float64)
		case int:
			nowValue = float64(now.(int))
		case int64:
			nowValue = float64(now.(int64))
		default:
			return false, errors.Join(ErrInvalidRule, fmt.Errorf("%v is %T", now, now))
		}
		switch t.CompareType {
		case CompareTypeLess:
			return nowValue < t.Value, nil
		case CompareTypeLessOrEqual:
			return nowValue <= t.Value, nil
		case CompareTypeEqual:
			return nowValue == t.Value, nil
		case CompareTypeGreaterOrEqual:
			return nowValue >= t.Value, nil
		case CompareTypeGreater:
			return nowValue > t.Value, nil
		}
	case ThresholdTypeSize, ThresholdTypeSpeed:
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
	return false, fmt.Errorf("not support %#v", t)
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

func (tt *ThresholdType) Name() string {
	switch *tt {
	case ThresholdTypePercent:
		return "percent"
	case ThresholdTypeSize:
		return "size"
	case ThresholdTypeSpeed:
		return "speed"
	case ThresholdTypeTemperature:
		return "temperature"
	}
	return "unknown"
}
