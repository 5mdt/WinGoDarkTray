package main

import (
	"strconv"
	"strings"
)

const versionParts = 3 // major.minor.patch

func isVersionNewer(latest, current string) bool {
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	maxParts := versionParts
	if len(latestParts) < maxParts || len(currentParts) < maxParts {
		return false
	}

	for i := 0; i < maxParts; i++ {
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
