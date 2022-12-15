package mapping

import "strings"

type Entry struct {
	Hash              string
	RawFilename       string
	ProcessedFilename string
}

func (e *Entry) FromString(str string) {
	split := strings.Split(str, "|")
	if len(split) == 3 {
		e.Hash = split[0]
		e.RawFilename = split[1]
		e.ProcessedFilename = split[2]
	}
}

func (e *Entry) ToString() string {
	return e.Hash + "|" + e.RawFilename + "|" + e.ProcessedFilename
}

type entrySlice []*Entry

func (e entrySlice) haveHash(hash string) bool {
	for _, ent := range e {
		if ent.Hash == hash {
			return true
		}
	}

	return false
}
