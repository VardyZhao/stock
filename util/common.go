package util

import (
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func ParseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func ConvertFloatStrToInt(s string) int {
	// 将字符串转换为浮点数
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	// 乘以100并转换为整数
	intVal := int(floatVal * 100)
	return intVal
}

func ReadBodyWithCharset(resp *http.Response) (string, error) {
	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ExtractOrgCode 提取营业部编码
func ExtractOrgCode(url string) (string, error) {
	re := regexp.MustCompile(`orgcode/([a-zA-Z0-9]+)/?`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("无法匹配到代码")
	}
	return matches[1], nil
}
