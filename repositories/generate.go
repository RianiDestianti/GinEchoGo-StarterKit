package repositories

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

func GenerateNumber(length int) string {
	characters := "1234567890"
	result := make([]byte, length)
	for index := range result {
		result[index] = characters[rand.IntN(len(characters))]
	}
	return string(result)
}

func FormatDurationSmart(duration time.Duration) string {
	days := duration.Hours() / 24
	hours := duration.Hours()
	minutes := duration.Minutes()
	seconds := duration.Seconds()

	if days >= 1 {
		return fmt.Sprintf("%d Hari", int(math.Floor(days)))
	} else if hours >= 1 {
		return fmt.Sprintf("%d Jam", int(math.Floor(hours)))
	} else if minutes >= 1 {
		return fmt.Sprintf("%d Menit", int(math.Floor(minutes)))
	}
	return fmt.Sprintf("%d Detik", int(math.Floor(seconds)))
}

func StrWords(text string, words int, end string) string {
	if text == "" {
		return ""
	}

	splitText := strings.Fields(text)
	if len(splitText) <= words {
		return text
	}

	return strings.Join(splitText[:words], " ") + end
}

func ToInt(value string, fallback int) int {
	number, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return number
}
