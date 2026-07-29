//go:build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/transform"

	"quote_generator/tools/sources"
)

func main() {
	sourceName := flag.String("source", "yuwenmi", "数据源名称，可选: "+availableSources())
	listSources := flag.Bool("list", false, "列出所有可用数据源")
	keyword := flag.String("keyword", "", "额外搜索的关键词（可选）")
	flag.Parse()

	if *listSources {
		fmt.Println("可用数据源:")
		for _, s := range sources.DefaultSources {
			fmt.Printf("  %-12s %s\n", s.Name, s.Description)
		}
		return
	}

	src := sources.GetSourceByName(*sourceName)
	if src == nil {
		fmt.Printf("❌ 未知数据源: %s\n", *sourceName)
		os.Exit(1)
	}

	fmt.Printf("正在探测数据源: %s\n", src.Description)
	fmt.Println("URL:", src.URL)

	resp, err := http.Get(src.URL)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if src.Encoding != nil {
		reader = transform.NewReader(resp.Body, src.Encoding.NewDecoder())
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("❌ 读取失败: %v\n", err)
		os.Exit(1)
	}

	html := string(body)
	fmt.Printf("页面总长度: %d 字节\n\n", len(html))

	// 用正则匹配看看效果
	matches := src.Pattern.FindAllStringSubmatch(html, -1)
	fmt.Printf("正则匹配到 %d 条结果\n", len(matches))
	if len(matches) > 0 {
		fmt.Println("\n前 3 条匹配结果:")
		for i := 0; i < 3 && i < len(matches); i++ {
			fmt.Printf("  [%d] 全文: %s\n", i+1, matches[i][0])
			fmt.Printf("      名言: %s\n", strings.TrimSpace(matches[i][1]))
			if len(matches[i]) > 2 && matches[i][2] != "" {
				fmt.Printf("      作者: %s\n", strings.TrimSpace(matches[i][2]))
			}
		}
	}

	// 搜索关键词附近内容
	kw := *keyword
	if kw == "" {
		kw = "1、"
	}
	idx := strings.Index(html, kw)
	if idx > 0 {
		start := idx - 80
		if start < 0 {
			start = 0
		}
		end := idx + 300
		if end > len(html) {
			end = len(html)
		}
		fmt.Printf("\n关键词 \"%s\" 附近内容:\n", kw)
		fmt.Println(html[start:end])
	} else {
		fmt.Printf("\n未找到关键词 \"%s\"\n", kw)
	}
}

func availableSources() string {
	names := make([]string, len(sources.DefaultSources))
	for i, s := range sources.DefaultSources {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}
