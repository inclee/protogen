package utool

import (
	"strings"
	"unicode"
)

// 工具函数：首字母小写
func LowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func Lower(s string) string {
	return strings.ToLower(s)
}

// camelToSnakeCase 将驼峰命名转换为下划线命名
func CamelToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func CamelToSplitCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '/')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
