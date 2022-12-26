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
		ent := new(Entry)
		ent.FromString(v)
		mp.m = append(mp.m, ent)
	}

	return nil
}

// Add adds a new mapping.
// Add evaluates the hash by filename and eTag.
func (mp *Mapping) Add(eTag, rawFilename, processedFilename string) {
	ent := new(Entry)
	ent.Hash = genHash(rawFilename, eTag)
	ent.RawFilename = rawFilename
	ent.ProcessedFilename = processedFilename
	mp.m = append(mp.m, ent)
}

// Append adds the specified mapping entry to `Mapping`.
func (mp *Mapping) Append(ent *Entry) {
	mp.m = append(mp.m, ent)
}

// Remove removes a mapping entry with specified hash.
func (mp *Mapping) Remove(hash string) {
	for index, ent := range mp.m {
		if ent.Hash == hash {
			mp.m[index] = mp.m[len(mp.m)-1]
			mp.m = mp.m[:len(mp.m)-1]

			return
		}
	}
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
		os.Mkdir("data", 0664)
	}

	return os.WriteFile(mp.path, []byte(data), 0664)
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
