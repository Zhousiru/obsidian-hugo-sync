package postprocess

const errStr = "[OBSIDIAN_HUGO_SYNC_ERROR]"

func Process(post, filename, baseUrl, vaultPost, vaultAsset string) string {
	post = modifyTitle(filename, post)
	post = convertImageMark(post, baseUrl, vaultAsset)
	post = convertLinkMark(post, baseUrl, vaultPost)
	post = convertWikilink(post, vaultPost, vaultAsset, baseUrl)

	return post
}
