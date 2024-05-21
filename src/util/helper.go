package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// email verify
func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Hide(str string) (result string) {
	if len(str) == 0 {
		return
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		star := ""
		if len(res[0]) < 3 {
			star = "***"
		} else {
			star = Substr(str, 0, 3) + "***"
		}
		result = star + "@" + res[1]
		return
	}
	reg := `^1[0-9]\d{9}$`
	rgx := regexp.MustCompile(reg)
	mobileMatch := rgx.MatchString(str)
	if mobileMatch {
		result = Substr(str, 0, 3) + "****" + Substr(str, 7, 11)
	} else {
		nameRune := []rune(str)
		lens := len(nameRune)
		if lens <= 1 {
			result = "***"
		} else if lens == 2 {
			result = string(nameRune[:1]) + "*"
		} else if lens == 3 {
			result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
		} else if lens == 4 {
			result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
		} else if lens > 4 {
			result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
		}
	}
	return
}

// Substr 截取字符
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

func AmountToString(value int64) string {
	if value < 0 {
		return "0.00"
	}
	amountString := strconv.FormatInt(value, 10)
	switch len(amountString) {
	case 1:
		amountString = "00" + amountString
	case 2:
		amountString = "0" + amountString
	}
	return amountString[:len(amountString)-2] + "." + amountString[len(amountString)-2:]
}

var (
	headerNums    = [...]string{"139", "138", "137", "136", "135", "134", "159", "158", "157", "150", "151", "152", "188", "187", "182", "183", "184", "178", "130", "131", "132", "156", "155", "186", "185", "176", "133", "153", "189", "180", "181", "177"}
	headerNumsLen = len(headerNums)
)

func RandomPhone() string {
	rand.Seed(time.Now().UTC().UnixNano())
	header := headerNums[rand.Intn(headerNumsLen)]
	body := fmt.Sprintf("%08d", rand.Intn(99999999))
	phone := header + body
	return phone
}

func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

// 去重 空间换时间
func RemoveDuplicateElement(args []string) []string {
	result := make([]string, 0, len(args))
	temp := map[string]struct{}{}
	for _, item := range args {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Int32Ptr(i int32) *int32 {
	return &i
}

// StringToArray 字符串转数组
func StringToArray(str string) []string {
	var arr []string
	if len(str) > 0 {
		arr = strings.Split(str, ",")
	}
	return arr
}

// StringInArray 字符串是否在数组中
func StringInArray(arr []string, str string) bool {
	if len(arr) > 0 {
		for _, v := range arr {
			if v == str {
				return true
			}
		}
	}
	return false
}

// StringContainsArray 字符串是否在数组中
func StringContainsArray(arr []string, str string) bool {
	if len(arr) > 0 {
		for _, v := range arr {
			if strings.Contains(v, str) {
				return true
			}
		}
	}
	return false
}

func ContainsChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func ConvertToSliceOfInts(data interface{}) ([]int, bool) {
	var result []int
	slice, ok := data.([]interface{})
	if !ok {
		return nil, false
	}
	for _, item := range slice {
		if num, ok := item.(float64); ok {
			result = append(result, int(num))
		} else {
			return nil, false
		}
	}
	return result, true
}

// IsUrl 判断是否为 URL
func IsUrl(path string) bool {
	u, err := url.ParseRequestURI(path)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// RoundToTwoDecimalPlaces 保留两位小数
func RoundToTwoDecimalPlaces(f float64) float64 {
	return math.Round(f*100) / 100
}

// RoundToFourDecimalPlaces 保留四位小数
func RoundToFourDecimalPlaces(f float64) float64 {
	return math.Round(f*10000) / 10000
}

func ReplacerServiceName(name string) string {
	replacer := strings.NewReplacer(
		"_", "-",
		".", "-",
		"::", "-", // 这个可能不需要，因为前一个已经将单个冒号替换了
		":", "-",
	)
	return strings.ToLower(replacer.Replace(name))
}

// LastNChars 截取最后多少字符
func LastNChars(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[len(runes)-n:])
}

// CleanString 过滤掉字符串中的所有非 UTF-8 字符和控制字符，只保留可打印的 UTF-8 字符
func CleanString(s string) string {
	var result []rune
	for _, r := range s {
		if r == utf8.RuneError {
			continue // 排除无效的 UTF-8 字符
		}
		if unicode.IsControl(r) {
			continue // 排除控制字符
		}
		if unicode.IsPrint(r) {
			result = append(result, r) // 保留可打印字符
		}
	}
	return string(result)
}

func ExtractJSONFromMarkdown(markdownText string) ([]map[string]interface{}, error) {
	// 正则表达式匹配 Markdown 代码块
	codeBlockRegex := regexp.MustCompile("```json\n(.*?)\n```")
	matches := codeBlockRegex.FindAllStringSubmatch(markdownText, -1)

	var jsonObjects []map[string]interface{}
	for _, match := range matches {
		if len(match) > 1 {
			jsonStr := match[1]
			var jsonObj map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
				// 如果解析失败，跳过这个代码块
				continue
			}
			jsonObjects = append(jsonObjects, jsonObj)
		}
	}

	return jsonObjects, nil
}

// UnescapeString 反转义字符串
func UnescapeString(s string) string {
	// 处理常见的转义序列
	replacements := map[string]string{
		"\\n":  "\n",
		"\\r":  "\r",
		"\\t":  "\t",
		"\\\"": "\"",
		"\\'":  "'",
		"\\\\": "\\",
		// 可以在这里添加更多的转义序列
	}

	for old, _new := range replacements {
		s = strings.ReplaceAll(s, old, _new)
	}

	return s
}
