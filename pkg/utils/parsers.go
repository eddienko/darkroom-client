package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// parseMemory converts values like "512Mi" or "2Gi" into MiB.
func ParseMemory(mem string) (int, error) {
	mem = strings.TrimSpace(strings.ToLower(mem))
	if strings.HasSuffix(mem, "gi") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(mem, "gi"), 64)
		if err != nil {
			return 0, err
		}
		return int(v * 1024), nil
	}
	if strings.HasSuffix(mem, "mi") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(mem, "mi"), 64)
		if err != nil {
			return 0, err
		}
		return int(v), nil
	}
	return 0, fmt.Errorf("unknown memory unit (must end with Mi or Gi)")
}
