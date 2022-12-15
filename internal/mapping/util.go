package mapping

import (
	"crypto/md5"
	"encoding/hex"
)

func genHash(rawFilename, eTag string) string {
	hash := md5.Sum([]byte(rawFilename + eTag))
	return hex.EncodeToString(hash[:])
}
