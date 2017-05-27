package app

type TestStorage interface {
	Register(t *NewTest) string
}
