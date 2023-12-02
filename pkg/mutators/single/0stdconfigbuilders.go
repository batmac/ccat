package mutators

import (
	"strconv"

	"braces.dev/errtrace"
	"github.com/batmac/ccat/pkg/stringutils"
)

func stdConfigHumanSizeAsInt64(args []string) (any, error) {
	if len(args) != 1 {
		return nil, errtrace.Wrap(ErrWrongNumberOfArgs(1, 1, len(args)))
	}

	n, err := stringutils.FromHumanSize[int64](args[0])
	return n, errtrace.Wrap(err)
}

func stdConfigHumanSizeAsInt64WithDefault(defaultValue int64) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		return errtrace.Wrap2(stdConfigHumanSizeAsInt64(args))
	}
}

func stdConfigUint64WithDefault(defaultValue uint64) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		if len(args) != 1 {
			return nil, errtrace.Wrap(ErrWrongNumberOfArgs(1, 1, len(args)))
		}

		if args[0][0] == '-' {
			n, err := strconv.ParseUint(args[0][1:], 10, 64)
			if err != nil {
				return nil, errtrace.Wrap(err)
			}
			return ^uint64(0) - n + 1, nil
		}

		return errtrace.Wrap2(strconv.ParseUint(args[0], 10, 64))
	}
}

func stdConfigStringWithDefault(defaultValue string) configBuilder {
	return func(args []string) (any, error) {
		if len(args) == 0 {
			return defaultValue, nil
		}
		return errtrace.Wrap2(stdConfigString(args))
	}
}

func stdConfigString(args []string) (any, error) {
	if len(args) != 1 {
		return nil, errtrace.Wrap(ErrWrongNumberOfArgs(1, 1, len(args)))
	}

	return args[0], nil
}

func stdConfigStrings(min, max int) configBuilder {
	return func(args []string) (any, error) {
		if len(args) < min || len(args) > max {
			return nil, errtrace.Wrap(ErrWrongNumberOfArgs(min, max, len(args)))
		}

		return args, nil
	}
}
