package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/alexdzyoba/bin/target"
)

type Archive struct{}

func (a *Archive) Extract(source *os.File, matcher target.Matcher) (*bytes.Reader, error) {
	gr, err := gzip.NewReader(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gzip reader")
	}

	tr := tar.NewReader(gr)

	var buf bytes.Buffer
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return nil, errors.Wrap(err, "failed to iterate tar")
		}

		buf.Reset()
		_, err = io.CopyN(&buf, tr, hdr.Size)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read tar")
		}

		r := bytes.NewReader(buf.Bytes())
		if matcher.Match(r) {
			return r, nil
		}
	}

	return nil, fmt.Errorf("no file matched")
}
