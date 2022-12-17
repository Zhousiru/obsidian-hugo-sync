package postprocess

const errStr = "[OBSIDIAN_HUGO_SYNC_ERROR]"

func Process(post, filename, baseUrl, vaultPost, vaultAsset string) string {
	post = modifyTitle(filename, post)
	post = convertCommonMarkImage(post, baseUrl, vaultAsset)
	post = convertWikilinkMark(post, vaultPost, vaultAsset, baseUrl)

	return post
}
