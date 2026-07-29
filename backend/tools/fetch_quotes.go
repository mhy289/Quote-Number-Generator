//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	url := "https://www.yuwenmi.com/yulu/mingren/519858.html"
	fmt.Println("正在抓取名言列表...")
	fmt.Println("来源:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("❌ 请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 页面是 GBK 编码，需要解码为 UTF-8
	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := io.ReadAll(utf8Reader)
	if err != nil {
		fmt.Println("❌ 读取失败:", err)
		return
	}

	html := string(body)

	// 匹配格式: <p>　　[数字]、[名言]——[作者]</p>
	// 或: <p>　　[数字]、[名言]</p>（无作者）
	re := regexp.MustCompile(`<p>[\s\x{3000}]*(\d+)[、.][\s\x{3000}]*([^<]+?)(?:——([^<]+))?</p>`)
	matches := re.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		fmt.Println("❌ 未匹配到任何名言，请检查页面结构")
		return
	}

	var quotes []string
	for _, m := range matches {
		quote := strings.TrimSpace(m[2])
		author := strings.TrimSpace(m[3])

		var fullQuote string
		if author != "" {
			fullQuote = fmt.Sprintf("%s —— %s", quote, author)
		} else {
			fullQuote = quote
		}
		quotes = append(quotes, fullQuote)
	}

	fmt.Printf("✅ 共抓取 %d 条名言\n", len(quotes))

	// 打印前 5 条预览
	fmt.Println("\n预览（前5条）:")
	for i := 0; i < 5 && i < len(quotes); i++ {
		fmt.Printf("  %d. %s\n", i+1, quotes[i])
	}

	// 保存到 data/quotes.json
	outputPath := "data/quotes.json"
	data, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		fmt.Println("❌ JSON 序列化失败:", err)
		return
	}

	// 确保 data 目录存在
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Println("❌ 创建目录失败:", err)
		return
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Println("❌ 写入文件失败:", err)
		return
	}

	fmt.Printf("\n✅ 名言已保存到: %s\n", outputPath)
	fmt.Printf("   共 %d 条名言，%d 字节\n", len(quotes), len(data))
}
