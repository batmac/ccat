package color

type Color interface {
	Sprint(s string) string
	Next() Color
}
