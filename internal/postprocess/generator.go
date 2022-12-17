package postprocess

import (
	"strings"
)

func genImage(arg, rawUrl, vaultAsset, baseUrl string) string {
	imageHtml := `<img src="` + getImageUrl(rawUrl, vaultAsset, baseUrl) + `"`

	if arg != "" {
		el := strings.Split(arg, "|")

		for _, v := range el {
			const altPrefix = "alt "
			if strings.HasPrefix(v, altPrefix) {
				// specify alt
				imageHtml = imageHtml + ` alt="` + v[len(altPrefix):] + `"`
			} else {
				// specify size
				// possible: `<Width>x<Height>`, `<Width>`

				size := strings.Split(v, "x")

				switch len(size) {
				case 1:
					// width only
					imageHtml = imageHtml + ` width="` + size[0] + `"`
				case 2:
					// both width and height
					imageHtml = imageHtml + ` width="` + size[0] + `" height="` + size[1] + `"`
				default:
					// unexpected
					return errStr
				}
			}
		}
	}

	imageHtml = imageHtml + ">"

	return imageHtml
}

func genPostLink(arg, rawUrl, vaultPost string) string {
	return `[` + arg + `]({{< relref "` + urlToFilename(rawUrl, vaultPost, urlTypePost) + `" >}})`
}

func genLink(arg, rawUrl string) string {
	return `<a href="` + rawUrl + `" target="_blank">` + arg + `</a>`
}
