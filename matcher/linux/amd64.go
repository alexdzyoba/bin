package linux

import (
	"bytes"
	"debug/elf"
	"log"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
)

type AMD64 struct{}

func (a *AMD64) Match(r *bytes.Reader) bool {
	t, err := filetype.MatchReader(r)
	if err != nil {
		log.Println(err)
		return false
	}

	if t == matchers.TypeElf {
		e, err := elf.NewFile(r)
		if err != nil {
			log.Println(err)
			return false
		}

		if e.Machine == elf.EM_X86_64 {
			return true
		}
	}

	return false
}
