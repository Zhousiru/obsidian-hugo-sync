package mapping

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
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
	m    map[string]string
	path string
}

// SetType sets the type of mapping.
// SetType must be called to init.
func (mp *Mapping) SetType(t string) {
	mp.m = make(map[string]string)
	mp.path = t
}

// Load loads mapping data from file.
// If file isn't exist, it will do nothing.
func (mp *Mapping) Load() error {
	if !util.IsExist("data/post_mapping") || !util.IsExist("data/asset_mapping") {
		// keep `mp.m` empty
		return nil
	}

	data, err := os.ReadFile(mp.path)
	if err != nil {
		return err
	}

	for _, v := range strings.Split(string(data), "\n") {
		splited := strings.Split(v, "|")
		if len(splited) != 2 {
			continue
		}
		mp.m[splited[0]] = splited[1]
	}

	return nil
}

// Add adds a new mapping.
// Add evaluates the hash by filename and eTag.
func (mp *Mapping) Add(filename, eTag, v string) {
	hash := md5.Sum([]byte(filename + eTag))
	mp.m[hex.EncodeToString(hash[:])] = v
}

// Diff finds out which mapping is created or deleted, and call back.
func (mp *Mapping) Diff(newMapping *Mapping, delCallback, newCallback func(hash, v string)) {
	tmpNew := make(map[string]string)

	for k, v := range newMapping.m {
		tmpNew[k] = v
	}

	for k, v := range mp.m {
		_, ok := tmpNew[k]
		if !ok {
			delCallback(k, v)
		} else {
			delete(tmpNew, k)
		}
	}

	for k, v := range tmpNew {
		newCallback(k, v)
	}
}

// Save saves mapping to file
func (mp *Mapping) Save() error {
	data := ""
	for k, v := range mp.m {
		data += k + "|" + v + "\n"
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
			if assetconv.CanToWebP(filepath.Ext(rawFilename)) {
				// can be converted to WebP
				// then filename is `xxx.webp`
				filename = util.TrimExt(rawFilename) + ".webp"
			} else {
				// can not be converted to WebP
				// then `filename` = `rawFilename`
				filename = rawFilename
			}
		}
		m.Add(rawFilename, obj.ETag, filename)
	}

	return m
}
