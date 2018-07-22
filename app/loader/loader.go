package loader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vbogretsov/sendmail/app"
	"github.com/vbogretsov/sendmail/app/loader/fs"
)

const (
	errorRootFormat          = "templates path must start with protocol://"
	errorUnsupportedProtocol = "protocol %s is unsupported"
	protocolDelimiter        = "://"
)

type factory func(string) (app.Loader, error)

var loaders = map[string]factory{
	"fs": fs.New,
}

// New creates a new loader, the exact loader type is being detected by prefix.
func New(root string) (app.Loader, error) {
	n := strings.Index(root, protocolDelimiter)
	if n == -1 {
		return nil, errors.New(errorRootFormat)
	}

	prefix := root[0:n]

	fn, ok := loaders[prefix]
	if !ok {
		return nil, fmt.Errorf(errorUnsupportedProtocol, prefix)
	}

	return fn(root[n+len(protocolDelimiter):])
}
