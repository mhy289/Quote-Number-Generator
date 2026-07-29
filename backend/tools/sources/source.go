// Package sources 定义名言抓取的数据源结构，支持多网站扩展。
// 每个 Source 封装了：网站 URL、编码方式、正则匹配规则、解析逻辑。
package sources

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Quote 表示一条解析后的名言
type Quote struct {
	Text   string // 名言正文
	Author string // 作者（可能为空）
}

// FullQuote 返回带作者的名言完整字符串
func (q Quote) FullQuote() string {
	if q.Author != "" {
		return fmt.Sprintf("%s —— %s", q.Text, q.Author)
	}
	return q.Text
}

// Source 定义一个名言数据源
type Source struct {
	Name        string            // 数据源名称（用于显示和命令行选择）
	URL         string            // 网页 URL
	Encoding    encoding.Encoding // 网页编码（nil 表示 UTF-8）
	Pattern     *regexp.Regexp    // 匹配名言的正则表达式
	Description string            // 数据源描述
}

// Fetch 抓取并解析该数据源的名言列表
func (s *Source) Fetch() ([]Quote, error) {
	resp, err := http.Get(s.URL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if s.Encoding != nil {
		reader = transform.NewReader(resp.Body, s.Encoding.NewDecoder())
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取失败: %w", err)
	}

	html := string(body)
	matches := s.Pattern.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("未匹配到任何名言，请检查页面结构或正则")
	}

	var quotes []Quote
	seen := make(map[string]bool)
	for _, m := range matches {
		text := strings.TrimSpace(m[1])
		if text == "" || seen[text] {
			continue
		}
		seen[text] = true

		q := Quote{Text: text}
		if len(m) > 2 {
			q.Author = strings.TrimSpace(m[2])
		}
		quotes = append(quotes, q)
	}

	return quotes, nil
}

// ==================== 预定义数据源 ====================

// YuwenmiSource 语文迷网名言数据源
// 页面结构: <p>　　[数字]、[名言]——[作者]</p> 或 <p>　　[数字]、[名言]</p>
var YuwenmiSource = &Source{
	Name:        "yuwenmi",
	URL:         "https://www.yuwenmi.com/yulu/mingren/519858.html",
	Encoding:    simplifiedchinese.GBK,
	Description: "语文迷网 - 名人名言（约 300 条）",
	Pattern: regexp.MustCompile(
		`<p>[\s\x{3000}]*\d+[、.][\s\x{3000}]*([^<]+?)(?:——([^<]+))?</p>`,
	),
}

// DefaultSources 所有可用的数据源列表
// 新增数据源时，在这里 append 即可
var DefaultSources = []*Source{
	YuwenmiSource,
}

// GetSourceByName 根据名称查找数据源
func GetSourceByName(name string) *Source {
	for _, s := range DefaultSources {
		if s.Name == name {
			return s
		}
	}
	return nil
}
