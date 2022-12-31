package postprocess

const errStr = "[OBSIDIAN_HUGO_SYNC_ERROR]"

func Process(post, filename, baseUrl, vaultPost, vaultAsset, assetUrl string) string {
	bypass := false

	post = modifyFrontMatter(
		post,
		func(s string) string {
			return titleModifier(s, filename)
		},
		func(s string) string {
			return assetRefModifier(s, vaultAsset, assetUrl)
		},
		func(s string) string {
			bypass = DetectBypassFlag(s)
			return s
		},
	)

	if bypass {
		return post
	}

	post = convertImageMark(post, baseUrl, vaultAsset)
	post = convertLinkMark(post, baseUrl, vaultPost)
	post = convertWikilink(post, vaultPost, vaultAsset, baseUrl)

	return post
}
