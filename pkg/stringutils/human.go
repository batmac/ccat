package stringutils

import (
	"braces.dev/errtrace"
	"github.com/docker/go-units"
	"golang.org/x/exp/constraints"
)

func HumanSize[T constraints.Integer](n T) string {
	return units.CustomSize("%.2f%s", float64(n), 1000, []string{"", "K", "M", "G", "T", "P", "E", "Z", "Y"})
}

func FromHumanSize[T constraints.Integer](size string) (T, error) {
	n, err := units.FromHumanSize(size)
	return T(n), errtrace.Wrap(err)
}
