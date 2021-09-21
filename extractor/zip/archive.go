package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/alexdzyoba/bin/matcher"
)

type Archive struct{}

func (a *Archive) Extract(source *os.File, matcher matcher.Matcher) (*bytes.Reader, error) {
	info, err := source.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat source file %v", source.Name())
	}

	zr, err := zip.NewReader(source, info.Size())
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			log.Println(err)
			continue
		}

		r := bytes.NewReader(b)
		if matcher.Match(r) {
			return r, nil
		}
	}

	return nil, fmt.Errorf("no file matched")
}
