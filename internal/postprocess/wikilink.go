package postprocess

import (
	"regexp"
	"strings"
)

func convertWikilinkMark(post, vaultPost, vaultAsset, baseUrl string) string {
	reWikilinkMark := regexp.MustCompile(`(!?)\[\[(.*?)\]\]`)

	ret := reWikilinkMark.ReplaceAllStringFunc(post, func(wikilinkMark string) string {
		submatch := reWikilinkMark.FindStringSubmatch(wikilinkMark)

		flag := submatch[1]
		if flag == "!" {
			// `![[...]]`
			// regard as a image
			splited := strings.Split(submatch[2], "|")

			switch len(splited) {
			case 1:
				// `![[<URL>]]`
				return genImage("", splited[0], vaultAsset, baseUrl)
			default:
				// `![[<URL>|<Arg 1>|...]]`
				return genImage(strings.Join(splited[1:], "|"), splited[0], vaultAsset, baseUrl)
			}
		} else {
			// `[[...]]`
			// regard as a post link
			splited := strings.Split(submatch[2], "|")

			switch len(splited) {
			case 1:
				// `[[<URL>]]`
				return genPostLink(splited[0], splited[0], vaultPost)
			case 2:
				// `[[<URL>|<Display Text>]]`
				return genPostLink(splited[1], splited[0], vaultPost)
			default:
				// unexpected
				return errStr
			}
		}
	})

	return ret
}
