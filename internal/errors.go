package internal

import "gox/internal/scanning"

type RuntimeError struct {
	Error error
	Token *scanning.Token
}
