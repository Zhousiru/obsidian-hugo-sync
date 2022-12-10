package postprocess

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/assetconv"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

// ConvertImageMark converts Markdown image mark to HTML image tag.
// And modify its URL.
func ConvertImageMark(post string, baseUrl string, vaultAsset string) string {
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
	src := ""

	if strings.HasPrefix(rawUrl, "https://") || strings.HasPrefix(rawUrl, "http://") {
		// remote `rawUrl`: `http(s)://...`
		src = rawUrl
	} else {
		// local `rawUrl`: `<Vault Asset>/xxx.jpg`, `xxx.jpg`
		imageFilename := ""

		if strings.HasPrefix(rawUrl, vaultAsset+"/") {
			// `<Vault Asset>/xxx.jpg`
			imageFilename = rawUrl[len(vaultAsset+"/"):]
		} else {
			// `xxx.jpg`
			imageFilename = rawUrl
		}

		if assetconv.CanToWebP(filepath.Ext(imageFilename)) {
			src = baseUrl + util.TrimExt(imageFilename) + ".webp"
		} else {
			src = baseUrl + imageFilename
		}
	}

	imageHtmlTag := `<img src="` + src + `"`

	if rawAlt != "" {
		// use `alt xxx` to specify HTML alt

		const altPrefix = "alt "
		if strings.HasPrefix(rawAlt, altPrefix) {
			// has specific alt

			imageHtmlTag = imageHtmlTag + ` alt="` + rawAlt[len(altPrefix):] + `"`
		} else {
			// no specific alt, add width and height
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
