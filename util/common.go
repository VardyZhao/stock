package util

import (
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ParseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// 读取并自动转换编码的响应内容
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
