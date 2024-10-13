package strategies

type Route struct {
	Method              string
	PathMatcher         func(path string) bool
	RequestBodyHandler  func(in string) (string, error)
	ResponseBodyHandler func(in string) (string, error)
}

type Config struct {
	Routes []Route
}
