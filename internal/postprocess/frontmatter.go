package postprocess

import (
	"regexp"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

// modifyTitle replaces `title: {auto}` to `title: <Post Filename Without Ext>`
func modifyTitle(post, filename string) string {
	title := util.TrimExt(filename)

	reFrontMatter := regexp.MustCompile(`^---\n[\S\s]*?\n---`)
	ret := reFrontMatter.ReplaceAllStringFunc(post, func(frontMatter string) string {
		reTitle := regexp.MustCompile(`title: .*?\n`)
		return reTitle.ReplaceAllStringFunc(frontMatter, func(titleField string) string {
			return strings.ReplaceAll(titleField, "{{auto}}", title)
		})
	})

	return ret
}

// modifyAssetRef replaces `{{asset image.jpg}}` to `https://<Asset URL>/image.webp`
func modifyAssetRef(post, vaultAsset, assetUrl string) string {
	reFrontMatter := regexp.MustCompile(`^---\n[\S\s]*?\n---`)
	ret := reFrontMatter.ReplaceAllStringFunc(post, func(frontMatter string) string {
		reAsset := regexp.MustCompile(`{{asset (.*?)}}`)
		return reAsset.ReplaceAllStringFunc(frontMatter, func(assetFiedld string) string {
			url := assetFiedld[8 : len(assetFiedld)-2]
			return getImageUrl(url, vaultAsset, assetUrl)
		})
	})

	return ret
}
