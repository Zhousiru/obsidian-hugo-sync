package postprocess

import (
	"regexp"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

func modifyFrontMatter(post string, modifier ...func(string) string) string {
	reFrontMatter := regexp.MustCompile(`^---\n[\S\s]*?\n---`)
	ret := reFrontMatter.ReplaceAllStringFunc(post, func(frontMatter string) string {
		modified := frontMatter
		for _, m := range modifier {
			modified = m(modified)
		}
		return modified
	})

	return ret
}

// titleModifier replaces `title: {auto}` to `title: <Post Filename Without Ext>`
func titleModifier(frontMatter, filename string) string {
	title := util.TrimExt(filename)

	reTitle := regexp.MustCompile(`title: .*?\n`)

	return reTitle.ReplaceAllStringFunc(frontMatter, func(titleField string) string {
		return strings.ReplaceAll(titleField, "{{auto}}", title)
	})
}

// assetRefModifier replaces `{{asset image.jpg}}` to `https://<Asset URL>/image.webp`
func assetRefModifier(frontMatter, vaultAsset, assetUrl string) string {
	reAsset := regexp.MustCompile(`{{asset (.*?)}}`)

	return reAsset.ReplaceAllStringFunc(frontMatter, func(assetFiedld string) string {
		url := assetFiedld[8 : len(assetFiedld)-2]
		return getImageUrl(url, vaultAsset, assetUrl)
	})
}

// DetectBypassFlag returns the existence of bypass flag in front matter.
func DetectBypassFlag(frontMatter string) bool {
	return strings.Contains(frontMatter, "_ohs_bypass: true")
}
