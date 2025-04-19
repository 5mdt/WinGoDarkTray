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

	for i := 0; i < 3; i++ {
		if len(latestParts) <= i || len(currentParts) <= i {
			return false
		}

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
