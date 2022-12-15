package mapping

import (
	"os"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/assetconv"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/s3"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

const (
	PostMapping  = "data/post_mapping"
	AssetMapping = "data/asset_mapping"
)

// Mapping saves the relation between hash and value.
type Mapping struct {
	m    entrySlice
	path string
}

// SetType sets the type of mapping.
// SetType must be called to init.
func (mp *Mapping) SetType(t string) {
	mp.path = t
}

// Load loads mapping data from file.
// If file isn't exist, it will do nothing.
func (mp *Mapping) Load() error {
	if !util.IsExist(mp.path) {
		// keep `mp.m` empty
		return nil
	}

	data, err := os.ReadFile(mp.path)
	if err != nil {
		return err
	}

	for _, v := range strings.Split(string(data), "\n") {
		splited := strings.Split(v, "|")
		if len(splited) != 3 {
			// skip invalid lines
			continue
		}
		ent := new(entry)
		ent.FromString(v)
		mp.m = append(mp.m, ent)
	}

	return nil
}

// Add adds a new mapping.
// Add evaluates the hash by filename and eTag.
func (mp *Mapping) Add(eTag, rawFilename, processedFilename string) {
	ent := new(entry)
	ent.Hash = genHash(rawFilename, eTag)
	ent.rawFilename = rawFilename
	ent.processedFilename = processedFilename
	mp.m = append(mp.m, ent)
}

// Diff finds out which mapping is added or deleted.
func (mp *Mapping) Diff(newMapping *Mapping) (add entrySlice, del entrySlice) {
	common := entrySlice{}
	for _, ent := range newMapping.m {
		if mp.m.haveHash(ent.Hash) {
			// common ent
			common = append(common, ent)
		} else {
			// new ent
			add = append(add, ent)
		}
	}

	for _, ent := range mp.m {
		if !common.haveHash(ent.Hash) {
			// del ent
			del = append(del, ent)
		}
	}

	return
}

// Save saves mapping to file
func (mp *Mapping) Save() error {
	data := ""
	for _, ent := range mp.m {
		data += ent.ToString() + "\n"
	}

	if !util.IsExist("data") {
		os.Mkdir("data", 0644)
	}

	return os.WriteFile(mp.path, []byte(data), 0644)
}

// VaultToMapping dumps Obsidian vault post/asset object slice to post/asset mapping.
func VaultToMapping(objSlice []*s3.Object, prefix, mappingType string) *Mapping {
	m := new(Mapping)
	m.SetType(mappingType)
	for _, obj := range objSlice {
		rawFilename := util.TrimPrefix(obj.Key, prefix)
		filename := ""
		if mappingType == PostMapping {
			// when it's `PostMapping`, `rawFilename` == `filename`
			filename = rawFilename
		} else {
			// when it's `AssetMapping`
			if assetconv.CanToWebP(util.GetExt(rawFilename)) {
				// can be converted to WebP
				// then filename is `xxx.webp`
				filename = util.TrimExt(rawFilename) + ".webp"
			} else {
				// can not be converted to WebP
				// then `filename` = `rawFilename`
				filename = rawFilename
			}
		}
		m.Add(obj.ETag, rawFilename, filename)
	}

	return m
}
