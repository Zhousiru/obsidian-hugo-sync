package mapping

import "strings"

type entry struct {
	Hash              string
	rawFilename       string
	processedFilename string
}

func (e *entry) FromString(str string) {
	split := strings.Split(str, "|")
	if len(split) == 3 {
		e.Hash = split[0]
		e.rawFilename = split[1]
		e.processedFilename = split[2]
	}
}

func (e *entry) ToString() string {
	return e.Hash + "|" + e.rawFilename + "|" + e.processedFilename
}

type entrySlice []*entry

func (e entrySlice) haveHash(hash string) bool {
	for _, ent := range e {
		if ent.Hash == hash {
			return true
		}
	}

	return false
}
