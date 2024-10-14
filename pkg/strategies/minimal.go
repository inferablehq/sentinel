package strategies

import (
	"regexp"

	"github.com/inferablehq/sentinel/pkg/handlers"
)

func MinimalConfig() Config {
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
		},
	}
}
