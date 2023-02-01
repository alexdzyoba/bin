package extractor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/pkg/errors"

	"github.com/alexdzyoba/bin/extractor/direct"
	"github.com/alexdzyoba/bin/extractor/targz"
	"github.com/alexdzyoba/bin/extractor/zip"
	"github.com/alexdzyoba/bin/target"
)

type Extractor interface {
	Extract(*os.File, target.Matcher) (*bytes.Reader, error)
}

// Discover determine extractor by source file type
func Discover(f *os.File) (Extractor, error) {
	typ, err := filetype.MatchReader(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to match file type")
	}
	// Rewind after filetype matching
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, errors.Wrap(err, "failed to rewind temp file")
	}

	switch typ {
	case matchers.TypeZip:
		return new(zip.Archive), nil

	case matchers.TypeGz:
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create gzip reader")
		}
		defer gr.Close()

		t, err := filetype.MatchReader(gr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to detect type inside gzip")
		}

		if t != matchers.TypeTar {
			return nil, errors.Wrap(err, "unsupported archive inside gzip, expected tar.gz")
		}

		return new(targz.Archive), nil
	case matchers.TypeElf: // matchers.TypeMachO
		return new(direct.File), nil
	default:
		return nil, fmt.Errorf("unsupported filetype %v", typ)
	}
}
