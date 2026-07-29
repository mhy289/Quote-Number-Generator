//go:build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"quote_generator/tools/sources"
)

func main() {
	// 命令行参数
	sourceName := flag.String("source", "yuwenmi", "数据源名称，可选: "+availableSources())
	listSources := flag.Bool("list", false, "列出所有可用数据源")
	outputPath := flag.String("output", "data/quotes.json", "输出文件路径")
	flag.Parse()

	// 列出数据源
	if *listSources {
		fmt.Println("可用数据源:")
		for _, s := range sources.DefaultSources {
			fmt.Printf("  %-12s %s\n", s.Name, s.Description)
		}
		return
	}

	// 查找数据源
	src := sources.GetSourceByName(*sourceName)
	if src == nil {
		fmt.Printf("❌ 未知数据源: %s\n", *sourceName)
		fmt.Println("可用数据源:")
		for _, s := range sources.DefaultSources {
			fmt.Printf("  %s - %s\n", s.Name, s.Description)
		}
		os.Exit(1)
	}

	fmt.Printf("正在从 [%s] 抓取名言列表...\n", src.Description)
	fmt.Println("URL:", src.URL)

	quotes, err := src.Fetch()
	if err != nil {
		fmt.Printf("❌ 抓取失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 共抓取 %d 条名言\n", len(quotes))

	// 打印前 5 条预览
	fmt.Println("\n预览（前5条）:")
	for i := 0; i < 5 && i < len(quotes); i++ {
		fmt.Printf("  %d. %s\n", i+1, quotes[i].FullQuote())
	}

	// 提取纯文本列表（向后兼容）
	quoteTexts := make([]string, len(quotes))
	for i, q := range quotes {
		quoteTexts[i] = q.FullQuote()
	}

	// 保存到 JSON
	data, err := json.MarshalIndent(quoteTexts, "", "  ")
	if err != nil {
		fmt.Printf("❌ JSON 序列化失败: %v\n", err)
		os.Exit(1)
	}

	// 确保目录存在
	outputDir := *outputPath
	if idx := strings.LastIndex(outputDir, string(os.PathSeparator)); idx > 0 {
		dir := outputDir[:idx]
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ 创建目录失败: %v\n", err)
			os.Exit(1)
		}
	} else if idx := strings.LastIndex(outputDir, "/"); idx > 0 {
		dir := outputDir[:idx]
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ 创建目录失败: %v\n", err)
			os.Exit(1)
		}
	}

	if err := os.WriteFile(*outputPath, data, 0644); err != nil {
		fmt.Printf("❌ 写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ 名言已保存到: %s\n", *outputPath)
	fmt.Printf("   共 %d 条名言，%d 字节\n", len(quotes), len(data))
}

func availableSources() string {
	names := make([]string, len(sources.DefaultSources))
	for i, s := range sources.DefaultSources {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}
