package generic

import "strings"

func GetSubcategorySuffixFromMarkdownDescription(markdownDescription string) string {
	return " ||| " + strings.Split(markdownDescription, " ||| ")[1]
}
