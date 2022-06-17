package utils

import (
	"github.com/batmac/ccat/log"
	"github.com/hbollon/go-edlib"
)

func FuzzySearch(str string, strList []string, threshold float32) (string, error) {
	res, err := edlib.FuzzySearchThreshold(str, strList, threshold, edlib.JaroWinkler)
	if err != nil {
		log.Debugln(err)
	} else {
		log.Debugf("Result for '%s': %s", str, res)
	}
	return res, err
}
