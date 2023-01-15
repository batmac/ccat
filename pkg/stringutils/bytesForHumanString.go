package stringutils

import (
	"github.com/docker/go-units"
)

func BytesForHumanString(b uint64) string {
	return units.HumanSize(float64(b))
}
