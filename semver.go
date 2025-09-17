package main

import (
	"strconv"
	"strings"
)

func isVersionNewer(latest, current string) bool {
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	// Normalize to same length by padding with zeros
	maxLen := max(len(latestParts), len(currentParts))
	latestParts = padVersionParts(latestParts, maxLen)
	currentParts = padVersionParts(currentParts, maxLen)

	for i := 0; i < maxLen; i++ {
		latestNum, err1 := strconv.Atoi(latestParts[i])
		currentNum, err2 := strconv.Atoi(currentParts[i])
		if err1 != nil || err2 != nil {
			return false
		}

		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	return false
}

func padVersionParts(parts []string, targetLen int) []string {
	for len(parts) < targetLen {
		parts = append(parts, "0")
	}
	return parts
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
