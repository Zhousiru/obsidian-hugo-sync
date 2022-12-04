package postprocess

import (
	"regexp"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

func ReplaceImageUrl(post string, baseUrl string, vaultAsset string) string {
	reImageMark := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)

	ret := reImageMark.ReplaceAllStringFunc(post, func(imageMark string) string {
		submatch := reImageMark.FindStringSubmatch(imageMark)
		rawAlt := submatch[1]
		rawUrl := submatch[2]

		return genImageHtmlTag(rawAlt, rawUrl, vaultAsset, baseUrl)
	})

	return ret
}

func genImageHtmlTag(rawAlt string, rawUrl string, vaultAsset string, baseUrl string) string {
	// acceptable `rawUrl`: `<Vault Asset>/xxx.jpg`, `xxx.jpg`
	imageFilename := ""

	if strings.HasPrefix(rawUrl, vaultAsset+"/") {
		imageFilename = util.TrimExt(rawUrl[len(vaultAsset+"/"):])
	} else {
		imageFilename = util.TrimExt(rawUrl)
	}

	imageHtmlTag := `<img src="` + baseUrl + imageFilename + `.webp"`

	if rawAlt != "" {
		// use `alt xxx` to specify HTML alt

		const altPrefix = "alt "
		if strings.HasPrefix(rawAlt, altPrefix) {
			// has specific alt

			imageHtmlTag = imageHtmlTag + ` alt="` + rawAlt[len(altPrefix):] + `"`
		} else {
			// not specific alt, add width and height
			// acceptable: `100x200`, `100`

			size := strings.Split(rawAlt, "x")

			switch len(size) {
			case 1:
				// width only
				imageHtmlTag = imageHtmlTag + ` width="` + size[0] + `"`
			case 2:
				// both width and height
				imageHtmlTag = imageHtmlTag + ` width="` + size[0] + `" height="` + size[1] + `"`
			default:
				// unsupported format, ignore
			}
		}
	}

	imageHtmlTag = imageHtmlTag + ">"

	return imageHtmlTag
}
