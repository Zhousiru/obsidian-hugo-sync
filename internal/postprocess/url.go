package postprocess

import (
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/assetconv"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

const (
	urlTypePost = iota
	urlTypeAsset
)

func urlToFilename(url, prefix string, urlType int) string {
	var filename string
	if strings.HasPrefix(url, prefix) {
		// `<Prefix>xxx.xx`, `<Prefix>xxx`
		filename = url[len(prefix):]
	} else {
		// `xxx.xx`, `xxx`
		filename = url
	}

	if urlType == urlTypePost {
		return filename + ".md"
	}
	return filename
}

func getImageUrl(rawUrl, vaultAsset, baseUrl string) string {
	if strings.HasPrefix(rawUrl, "https://") || strings.HasPrefix(rawUrl, "http://") {
		// remote: `http(s)://...`
		return rawUrl
	}
	// local
	filename := urlToFilename(rawUrl, vaultAsset, urlTypeAsset)

	if assetconv.CanToWebP(util.GetExt(filename)) {
		return baseUrl + util.TrimExt(filename) + ".webp"
	}
	return baseUrl + filename
}
