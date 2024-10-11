package mutators

import (
	"strconv"

	"github.com/batmac/ccat/pkg/stringutils"
)

func stdConfigHumanSizeAsInt64(args []string) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs(1, 1, len(args))
	}

	n, err := stringutils.FromHumanSize[int64](args[0])
	return n, err
}

func stdConfigHumanSizeAsInt64WithDefault(defaultValue int64) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		return stdConfigHumanSizeAsInt64(args)
	}
}

func stdConfigUint64WithDefault(defaultValue uint64) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		if len(args) != 1 {
			return nil, ErrWrongNumberOfArgs(1, 1, len(args))
		}

		if args[0][0] == '-' {
			n, err := strconv.ParseUint(args[0][1:], 10, 64)
			if err != nil {
				return nil, err
			}
			return ^uint64(0) - n + 1, nil
		}

		return strconv.ParseUint(args[0], 10, 64)
	}
}

func stdConfigStringWithDefault(defaultValue string) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		return stdConfigString(args)
	}
}

func stdConfigString(args []string) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs(1, 1, len(args))
	}

	return args[0], nil
}

func stdConfigStrings(amin, amax int) configBuilder {
	return func(args []string) (any, error) {
		if len(args) < amin || len(args) > amax {
			return nil, ErrWrongNumberOfArgs(amin, amax, len(args))
		}

		return args, nil
	}
}

func stdConfigInts(amin, amax int) configBuilder {
	return func(args []string) (any, error) {
		if len(args) < amin || len(args) > amax {
			return nil, ErrWrongNumberOfArgs(amin, amax, len(args))
		}

		var ints []int
		for _, arg := range args {
			n, err := strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
			ints = append(ints, n)
		}

		return ints, nil
	}
}
