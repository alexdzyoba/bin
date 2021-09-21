package direct

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/alexdzyoba/bin/matcher"
)

type File struct{}

func (f *File) Extract(source *os.File, matcher matcher.Matcher) (*bytes.Reader, error) {
	b, err := io.ReadAll(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	r := bytes.NewReader(b)

	if matcher.Match(r) {
		return r, nil
	}

	return nil, fmt.Errorf("no file matched")
}
