//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	url := "https://www.yuwenmi.com/yulu/mingren/519858.html"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 用 GBK 解码（兼容 GB2312）
	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}

	html := string(body)

	// 找 "302" 附近
	idx := strings.Index(html, "302")
	if idx > 0 {
		start := idx - 50
		if start < 0 {
			start = 0
		}
		end := idx + 200
		if end > len(html) {
			end = len(html)
		}
		fmt.Println("302 附近内容:")
		fmt.Println(html[start:end])
	}

	// 找 "1、" 格式
	idx2 := strings.Index(html, "1、")
	if idx2 > 0 {
		start := idx2
		end := idx2 + 500
		if end > len(html) {
			end = len(html)
		}
		fmt.Println("\n\n1、附近内容:")
		fmt.Println(html[start:end])
	}

	// 统计名言数量 - 用 rune 方式
	count := 0
	runes := []rune(html)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '、' {
			j := i - 1
			for j >= 0 && runes[j] >= '0' && runes[j] <= '9' {
				j--
			}
			if i-j > 1 {
				count++
			}
		}
	}
	fmt.Printf("\n\n共找到约 %d 条名言\n", count)

	// 打印页面总长度
	fmt.Printf("页面总长度: %d 字节, %d 字符\n", len(html), utf8.RuneCountInString(html))
}
