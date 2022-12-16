package postprocess

func Process(post, filename, baseUrl, vaultAsset string) string {
	post = ModifyTitle(filename, post)
	post = ConvertImageMark(post, baseUrl, vaultAsset)

	return post
}
