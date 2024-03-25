/**
 * @Time: 2023/8/1 15:07
 * @Author: varluffy
 */

package util

import "strings"

// 优化prompt，将其中换行替换成句号
func OptimizePrompt(prompt string) string {
	prompt = strings.ReplaceAll(prompt, "\r\n", "。")
	prompt = strings.ReplaceAll(prompt, "\n", "。")
	prompt = strings.ReplaceAll(prompt, "\r", "。")
	return prompt
}
