package postprocess

import (
	"regexp"
	"strings"
)

// convertImageMark converts Markdown image mark to HTML image tag.
// And modify its URL.
func convertImageMark(post, baseUrl, vaultAsset string) string {
	reImageMark := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)

	ret := reImageMark.ReplaceAllStringFunc(post, func(imageMark string) string {
		submatch := reImageMark.FindStringSubmatch(imageMark)
		arg := submatch[1]
		rawUrl := submatch[2]

		return genImage(arg, rawUrl, vaultAsset, baseUrl)
	})

	return ret
}

// convertLinkMark converts Markdown link mark to HTML link tag.
func convertLinkMark(post, baseUrl, vaultPost string) string {
	reLinkMark := regexp.MustCompile(`(!?)\[(.*?)\]\((.*?)\)`)

	ret := reLinkMark.ReplaceAllStringFunc(post, func(linkMark string) string {
		submatch := reLinkMark.FindStringSubmatch(linkMark)
		flag := submatch[1]
		if flag == "!" {
			// it's a image mark
			return linkMark
		}

		arg := submatch[2]
		rawUrl := submatch[3]

		if strings.HasPrefix(rawUrl, "https://") || strings.HasPrefix(rawUrl, "http://") {
			// external link
			return genLink(arg, rawUrl)
		} else {
			// post link
			return genPostLink(arg, rawUrl, vaultPost)
		}
	})

	return ret
}
