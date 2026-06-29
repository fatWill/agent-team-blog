package wechat

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
)

// TiptapNode Tiptap JSON 节点结构
type TiptapNode struct {
	Type    string            `json:"type"`
	Content []TiptapNode      `json:"content,omitempty"`
	Text    string            `json:"text,omitempty"`
	Marks   []TiptapMark      `json:"marks,omitempty"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
}

// TiptapMark Tiptap 行内标记
type TiptapMark struct {
	Type  string                 `json:"type"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

// ConvertResult 转换结果
type ConvertResult struct {
	HTML      string   // 最终 HTML
	ImageURLs []string // 文中所有图片 URL（需要上传到微信）
}

// codeBlockStyle 代码块样式
const codeBlockStyle = `background-color:#f6f8fa;padding:16px;border-radius:4px;font-family:Consolas,Monaco,'Courier New',monospace;font-size:13px;line-height:1.6;overflow-x:auto;white-space:pre-wrap;word-wrap:break-word;`

// blockquoteStyle 引用块样式
const blockquoteStyle = `border-left:4px solid #dfe2e5;padding:8px 16px;margin:16px 0;color:#6a737d;`

// TiptapToWechatHTML 将 Tiptap JSON 转换为公众号兼容 HTML
func TiptapToWechatHTML(contentJSON []byte, articleID string) (*ConvertResult, error) {
	var doc TiptapNode
	if err := json.Unmarshal(contentJSON, &doc); err != nil {
		return nil, fmt.Errorf("解析 Tiptap JSON 失败: %w", err)
	}

	result := &ConvertResult{}
	var sb strings.Builder

	for _, node := range doc.Content {
		renderNode(&sb, &node, result)
	}

	// 末尾追加博客原文链接
	sb.WriteString(fmt.Sprintf(
		`<p style="margin-top:32px;padding-top:16px;border-top:1px solid #eee;font-size:14px;color:#999;">📖 本文同步发布于博客：<a href="https://fatwill.cloud/articles/%s">https://fatwill.cloud/articles/%s</a></p>`,
		articleID, articleID,
	))
	sb.WriteString(`<p style="font-size:14px;color:#999;">欢迎关注公众号，获取更多技术分享 🎉</p>`)

	result.HTML = sb.String()
	return result, nil
}

// renderNode 渲染单个节点
func renderNode(sb *strings.Builder, node *TiptapNode, result *ConvertResult) {
	switch node.Type {
	case "paragraph":
		sb.WriteString(`<p style="margin:16px 0;line-height:1.8;">`)
		renderInline(sb, node.Content, result)
		sb.WriteString("</p>")

	case "heading":
		level := 2
		if l, ok := node.Attrs["level"]; ok {
			if lf, ok := l.(float64); ok {
				level = int(lf)
			}
		}
		// h4+ 降级为加粗段落
		if level >= 4 {
			sb.WriteString(`<p style="margin:16px 0;line-height:1.8;"><strong>`)
			renderInline(sb, node.Content, result)
			sb.WriteString("</strong></p>")
		} else {
			sb.WriteString(fmt.Sprintf(`<h%d style="margin:24px 0 12px 0;">`, level))
			renderInline(sb, node.Content, result)
			sb.WriteString(fmt.Sprintf("</h%d>", level))
		}

	case "bulletList":
		sb.WriteString(`<ul style="padding-left:24px;margin:12px 0;">`)
		for _, item := range node.Content {
			renderNode(sb, &item, result)
		}
		sb.WriteString("</ul>")

	case "orderedList":
		sb.WriteString(`<ol style="padding-left:24px;margin:12px 0;">`)
		for _, item := range node.Content {
			renderNode(sb, &item, result)
		}
		sb.WriteString("</ol>")

	case "listItem":
		sb.WriteString(`<li style="margin:4px 0;line-height:1.8;">`)
		for _, child := range node.Content {
			// listItem 内部通常包裹 paragraph，直接渲染内联内容
			if child.Type == "paragraph" {
				renderInline(sb, child.Content, result)
			} else {
				renderNode(sb, &child, result)
			}
		}
		sb.WriteString("</li>")

	case "blockquote":
		sb.WriteString(fmt.Sprintf(`<blockquote style="%s">`, blockquoteStyle))
		for _, child := range node.Content {
			renderNode(sb, &child, result)
		}
		sb.WriteString("</blockquote>")

	case "codeBlock":
		sb.WriteString(fmt.Sprintf(`<pre style="%s"><code>`, codeBlockStyle))
		for _, child := range node.Content {
			if child.Text != "" {
				sb.WriteString(html.EscapeString(child.Text))
			}
		}
		sb.WriteString("</code></pre>")

	case "image":
		src := ""
		if s, ok := node.Attrs["src"]; ok {
			if str, ok := s.(string); ok {
				src = str
			}
		}
		if src != "" {
			result.ImageURLs = append(result.ImageURLs, src)
			sb.WriteString(fmt.Sprintf(`<img src="%s" style="max-width:100%%;margin:12px 0;" />`, html.EscapeString(src)))
		}

	case "horizontalRule":
		sb.WriteString(`<hr style="border:none;border-top:1px solid #eee;margin:24px 0;" />`)

	case "table":
		sb.WriteString(`<table style="border-collapse:collapse;width:100%;margin:16px 0;">`)
		for _, row := range node.Content {
			renderNode(sb, &row, result)
		}
		sb.WriteString("</table>")

	case "tableRow":
		sb.WriteString("<tr>")
		for _, cell := range node.Content {
			renderNode(sb, &cell, result)
		}
		sb.WriteString("</tr>")

	case "tableHeader":
		sb.WriteString(`<th style="border:1px solid #ddd;padding:8px;background-color:#f6f8fa;font-weight:bold;text-align:left;">`)
		for _, child := range node.Content {
			if child.Type == "paragraph" {
				renderInline(sb, child.Content, result)
			} else {
				renderNode(sb, &child, result)
			}
		}
		sb.WriteString("</th>")

	case "tableCell":
		sb.WriteString(`<td style="border:1px solid #ddd;padding:8px;text-align:left;">`)
		for _, child := range node.Content {
			if child.Type == "paragraph" {
				renderInline(sb, child.Content, result)
			} else {
				renderNode(sb, &child, result)
			}
		}
		sb.WriteString("</td>")

	case "hardBreak":
		sb.WriteString("<br/>")

	default:
		// 未知节点类型，尝试递归渲染子节点
		for _, child := range node.Content {
			renderNode(sb, &child, result)
		}
	}
}

// renderInline 渲染行内内容
func renderInline(sb *strings.Builder, nodes []TiptapNode, result *ConvertResult) {
	for _, node := range nodes {
		if node.Type == "image" {
			renderNode(sb, &node, result)
			continue
		}
		if node.Type == "hardBreak" {
			sb.WriteString("<br/>")
			continue
		}
		if node.Text == "" && len(node.Content) > 0 {
			renderNode(sb, &node, result)
			continue
		}

		text := html.EscapeString(node.Text)

		// 检测 mermaid 代码块（行内不太可能出现，但以防万一）
		// 应用行内标记
		for _, mark := range node.Marks {
			switch mark.Type {
			case "bold":
				text = "<strong>" + text + "</strong>"
			case "italic":
				text = "<em>" + text + "</em>"
			case "underline":
				text = "<u>" + text + "</u>"
			case "strike":
				text = "<s>" + text + "</s>"
			case "code":
				text = fmt.Sprintf(`<code style="background-color:#f0f0f0;padding:2px 4px;border-radius:3px;font-family:Consolas,monospace;font-size:0.9em;">%s</code>`, text)
			case "link":
				href := ""
				if h, ok := mark.Attrs["href"]; ok {
					if str, ok := h.(string); ok {
						href = str
					}
				}
				if href != "" {
					text = fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), text)
				}
			}
		}

		sb.WriteString(text)
	}
}

// CleanHTMLForWechat 清理 HTML 中不兼容公众号的内容
func CleanHTMLForWechat(htmlContent string) string {
	// 移除 <style> 标签
	htmlContent = removeTagContent(htmlContent, "style")
	// 移除 <script> 标签
	htmlContent = removeTagContent(htmlContent, "script")
	return htmlContent
}

// removeTagContent 移除指定标签及其内容（简单实现）
func removeTagContent(s, tag string) string {
	for {
		start := strings.Index(s, "<"+tag)
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "</"+tag+">")
		if end == -1 {
			// 自闭合或无结束标签，移除到下一个 >
			endTag := strings.Index(s[start:], ">")
			if endTag == -1 {
				break
			}
			s = s[:start] + s[start+endTag+1:]
		} else {
			s = s[:start] + s[start+end+len("</"+tag+">"):]
		}
	}
	return s
}
