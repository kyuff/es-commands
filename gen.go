package commands

import "github.com/kyuff/es"

//go:generate go tool moq -skip-ensure -pkg commands_test -rm -out mocks_test.go . Store State esStream:StreamMock esContent:ContentMock Middleware

type esStream es.Stream
type esContent es.Content
