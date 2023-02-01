package target

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/alexdzyoba/bin/target/linux"
)

const (
	PlatformLinuxAMD64 string = "linux-amd64"
)

var (
	installers = map[string]Matcher{
		PlatformLinuxAMD64: new(linux.AMD64),
	}
)

type Matcher interface {
	Match(r *bytes.Reader) bool
}

func Discover() (Matcher, error) {
	platform := runtime.GOOS + "-" + runtime.GOARCH
	ins, ok := installers[platform]
	if !ok {
		return nil, fmt.Errorf("unsupported platform %v", platform)
	}

	return ins, nil
}
