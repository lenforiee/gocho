package logger

import (
	"fmt"
	"time"
)

func formatDate() (dateFormat string) {
	currentTime := time.Now()
	dateFormat = currentTime.Format("2 Jan 2006 15:04:05")
	return dateFormat
}

func Info(message string) { // sooo cursed lmfao.
	fmt.Println("\033[37m\033[42m[INFO]\033[49m - [", formatDate(), "]", message, "\033[39m")
}

func Error(message string) { // sooo cursed lmfao.
	fmt.Println("\033[37m\033[41m[ERROR]\033[49m - [", formatDate(), "]", message, "\033[39m")
}

func Warning(message string) { // sooo cursed lmfao.
	fmt.Println("\033[37m\033[44m[WARNING]\033[49m - [", formatDate(), "]", message, "\033[39m")
}
