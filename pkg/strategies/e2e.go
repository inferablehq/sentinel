package strategies

import (
	"regexp"

	"github.com/inferablehq/sentinel/pkg/handlers"
)

func E2EConfig() Config {
	return Config{
		Routes: []Route{
			{
				Method:              "GET",
				PathMatcher:         func(path string) bool { return path == "/live" },
				RequestBodyHandler:  handlers.DefaultPassthroughHandler,
				ResponseBodyHandler: handlers.DefaultPassthroughHandler,
			},
			{
				Method:              "POST",
				PathMatcher:         func(path string) bool { return path == "/machines" },
				RequestBodyHandler:  handlers.DefaultPassthroughHandler,
				ResponseBodyHandler: handlers.DefaultPassthroughHandler,
			},
			{
				Method: "GET",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/calls$`).MatchString(path)
				},
				RequestBodyHandler:  handlers.DefaultPassthroughHandler,
				ResponseBodyHandler: handlers.DetokenizableDataHandler,
			},
			{
				Method: "POST",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/calls/\w+/result$`).MatchString(path)
				},
				RequestBodyHandler:  func(in string) (string, error) { return handlers.TokenizableDataHandler(in, []string{".resultType"}) },
				ResponseBodyHandler: handlers.DefaultPassthroughHandler,
			},
			{
				Method: "POST",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/runs$`).MatchString(path)
				},
				RequestBodyHandler: func(in string) (string, error) {
					return handlers.TokenizableDataHandler(in, []string{".result.schema", ".message"})
				},
				ResponseBodyHandler: handlers.DefaultPassthroughHandler,
			},
			{
				Method:              "GET",
				PathMatcher:         func(path string) bool { return regexp.MustCompile(`^/clusters/\w+/runs/\w+$`).MatchString(path) },
				RequestBodyHandler:  handlers.DefaultPassthroughHandler,
				ResponseBodyHandler: handlers.DetokenizableDataHandler,
			},
		},
	}
}
