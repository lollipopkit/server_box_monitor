package utils

import (
	"fmt"
	"time"
)

const (
	red       = "\033[91m"
	green     = "\033[32m"
	yellow    = "\033[93m"
	cyan      = "\033[96m"
	noColor   = "\033[0m"
	noColorLn = "\033[0m\n"
)

func print(s string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s]\n%s", t, s)
}

func Error(s string, f ...any) {
	print(fmt.Sprintf(red+s+noColorLn, f...))
}

func ErrorNln(s string, f ...any) {
	print(fmt.Sprintf(red+s+noColor, f...))
}

func Success(s string, f ...any) {
	print(fmt.Sprintf(green+s+noColorLn, f...))
}

func SuccessNln(s string, f ...any) {
	print(fmt.Sprintf(green+s+noColor, f...))
}

func Warn(s string, f ...any) {
	print(fmt.Sprintf(yellow+s+noColorLn, f...))
}

func WarnNln(s string, f ...any) {
	print(fmt.Sprintf(yellow+s+noColor, f...))
}

func Info(s string, f ...any) {
	print(fmt.Sprintf(cyan+s+noColorLn, f...))
}

func InfoNln(s string, f ...any) {
	print(fmt.Sprintf(cyan+s+noColor, f...))
}

func NoColor(s string, f ...any) {
	print(fmt.Sprintf(s+noColorLn, f...))
}

func NoColorNln(s string, f ...any) {
	print(fmt.Sprintf(s+noColor, f...))
}
