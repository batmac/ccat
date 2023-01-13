package mutators

import (
	"strconv"

	"github.com/docker/go-units"
)

func stdConfigHumanSizeAsInt64(args []string) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs(1, len(args))
	}

	n, err := units.FromHumanSize(args[0])
	if err != nil {
		return nil, err
	}
	return n, nil
}

func stdConfigUint64WithDefault(defaultValue uint64) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		if len(args) != 1 {
			return nil, ErrWrongNumberOfArgs(1, len(args))
		}

		n, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return nil, err
		}
		return n, nil
	}
}
