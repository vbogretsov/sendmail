package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vbogretsov/sendmail/app"
)

type fsloader struct {
	root string
}

// New creates a new loader. The exact loader type is determined by root prefix.
func New(root string) (app.Loader, error) {
	return fsloader{root: root}, nil
}

// Load loads a template with the language and name provided from the local
// file system.
func (ld fsloader) Load(lang, name string) (io.Reader, error) {
	fname := path.Join(ld.root, lang, fmt.Sprintf("%s.msg", name))

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(file), nil
}
