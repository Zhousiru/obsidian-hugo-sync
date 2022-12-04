package postprocess

import (
	"regexp"
	"strings"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

// ModifyTitle replaces `title: {auto}` to `title: <Post Filename Without Ext>`
func ModifyTitle(filename string, post string) string {
	title := util.TrimExt(filename)

	reFrontMatter := regexp.MustCompile(`^---\n[\S\s]*?\n---`)
	ret := reFrontMatter.ReplaceAllStringFunc(post, func(frontMatter string) string {
		reTitle := regexp.MustCompile(`title: .*?\n`)
		return reTitle.ReplaceAllStringFunc(frontMatter, func(titleField string) string {
			return strings.ReplaceAll(titleField, "{{auto}}", title)
		})
	})

	return ret
}
