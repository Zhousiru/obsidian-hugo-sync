package mapping

import (
	"os"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/s3"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

const (
	PostMapping  = "data/post_mapping"
	AssetMapping = "data/asset_mapping"
)

// Mapping saves the relation between ETag and value.
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
func (mp *Mapping) Add(eTag string, v string) {
	mp.m[eTag] = v
}

// Diff finds out which mapping is created or deleted, and call back.
func (mp *Mapping) Diff(newMapping *Mapping, delCallback func(eTag string, v string), newCallback func(eTag string, v string)) {
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
func VaultToMapping(objSlice []*s3.Object, prefix string, mappingType string) *Mapping {
	m := new(Mapping)
	m.SetType(PostMapping)
	for _, obj := range objSlice {
		m.Add(obj.ETag, util.TrimPrefix(obj.Key, prefix))
	}

	return m
}
