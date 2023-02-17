package mutators

import (
	"strconv"

	"github.com/batmac/ccat/pkg/stringutils"
)

func stdConfigHumanSizeAsInt64(args []string) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs(1, len(args))
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
			return nil, ErrWrongNumberOfArgs(1, len(args))
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
		return nil, ErrWrongNumberOfArgs(1, len(args))
	}

	return args[0], nil
}

func stdConfigStringsAtLeast(n int) configBuilder {
	return func(args []string) (any, error) {
		if len(args) < n {
			return nil, ErrWrongNumberOfArgs(n, len(args))
		}

		return args, nil
	}
}
