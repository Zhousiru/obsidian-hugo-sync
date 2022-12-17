package postprocess

import "regexp"

// convertImageMark converts Markdown image mark to HTML image tag.
// And modify its URL.
func convertCommonMarkImage(post, baseUrl, vaultAsset string) string {
	reImageMark := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)

	ret := reImageMark.ReplaceAllStringFunc(post, func(imageMark string) string {
		submatch := reImageMark.FindStringSubmatch(imageMark)
		rawAlt := submatch[1]
		rawUrl := submatch[2]

		return genImage(rawAlt, rawUrl, vaultAsset, baseUrl)
	})

	return ret
}
