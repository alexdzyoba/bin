package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/alexdzyoba/bin/extractor"
	"github.com/alexdzyoba/bin/matcher"
)

func install(url string, config Config) error {
	// Get real binary name from url
	filename := name(url)
	if filename == "" {
		return fmt.Errorf("failed to parse filename from URL %v\nTry to override it with --target.filename, -o options", url)
	}

	// Create temp file to download from url
	tmp, err := os.CreateTemp("", filename)
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	err = fetch(url, tmp)
	if err != nil {
		return err
	}
	// Rewind after download so reading file will be from the start
	_, err = tmp.Seek(0, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "failed to rewind temp file")
	}

	// Determine extractor from filetype
	extract, err := extractor.Discover(tmp)
	if err != nil {
		return errors.Wrap(err, "failed to discover source file type")
	}

	// Determine how we match file inside archive from platform
	matcher, err := matcher.Discover()
	if err != nil {
		return errors.Wrap(err, "failed to discover matcher")
	}

	// Extract the binary
	r, err := extract.Extract(tmp, matcher)
	if err != nil {
		log.Fatal(err)
	}

	// Write the binary into the target filepath
	if config.Target.Filename != "" {
		filename = config.Target.Filename
	}

	filepath := path.Join(config.Target.Dir, filename)
	err = write(r, filepath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// name parses url for a name without version and extension suffixes
func name(url string) string {
	filename := path.Base(url)
	if filename == "." || filename == "/" {
		return ""
	}

	delim := '-'
	fields := strings.FieldsFunc(filename, func(c rune) bool {
		if c == '-' || c == '_' {
			delim = c
			return true
		}
		return false
	})

	version := regexp.MustCompile(`^(v\d+\.|\d+)`) // e.g. v1.2.3, v3, 93
	platform := regexp.MustCompile(`(?i)(linux|windows|.*bsd|unknown)`)
	arch := regexp.MustCompile(`(?i)(i.86|amd64|x86|x64|x86[-_]64|arm.*|aarch.*|mips.*|powerpc.*|sparc.*|s390.*)`)

	i := 0
	for i = 0; i < len(fields); i++ {
		if version.MatchString(fields[i]) {
			break
		}
		if platform.MatchString(fields[i]) {
			break
		}
		if arch.MatchString(fields[i]) {
			break
		}
	}

	return strings.Join(fields[:i], string(delim))
}

func fetch(url string, f *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func write(r *bytes.Reader, filepath string) error {
	output, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(output, r)
	if err != nil {
		return err
	}

	return nil
}
