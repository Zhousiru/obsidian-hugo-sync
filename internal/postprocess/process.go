package postprocess

const errStr = "[OBSIDIAN_HUGO_SYNC_ERROR]"

func Process(post, filename, baseUrl, vaultPost, vaultAsset, assetUrl string) string {
	post = modifyTitle(post, filename)
	post = modifyAssetRef(post, vaultAsset, assetUrl)
	post = convertImageMark(post, baseUrl, vaultAsset)
	post = convertLinkMark(post, baseUrl, vaultPost)
	post = convertWikilink(post, vaultPost, vaultAsset, baseUrl)

	return post
}
