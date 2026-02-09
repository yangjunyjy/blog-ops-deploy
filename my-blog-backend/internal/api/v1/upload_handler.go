package api

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(uploadService services.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadMarkdown 上传并解析 Markdown 文件
// @Summary 上传Markdown文件
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Markdown文件"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/markdown [post]
func (h *UploadHandler) UploadMarkdown(c *gin.Context) {
	// 从表单获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件失败", err)
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".md" && ext != ".markdown" {
		response.Error(c, http.StatusBadRequest, "只支持 .md 或 .markdown 格式的文件", nil)
		return
	}

	// 验证文件大小（限制为 10MB）
	if file.Size > 10*1024*1024 {
		response.Error(c, http.StatusBadRequest, "文件大小不能超过 10MB", nil)
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logger.Error("打开文件失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "打开文件失败", err)
		return
	}
	defer src.Close()

	// 读取文件内容
	content, err := io.ReadAll(src)
	if err != nil {
		logger.Error("读取文件内容失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "读取文件失败", err)
		return
	}

	markdownContent := string(content)

	// 解析 Markdown 文件
	title := parseTitle(markdownContent, file.Filename)
	summary := parseSummary(markdownContent)

	// 存储文件到项目根目录下的 markdown_files 目录
	uploadDir, err := filepath.Abs("../markdown_files")
	if err != nil {
		logger.Error("获取绝对路径失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "获取路径失败", err)
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("创建目录失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "创建目录失败", err)
		return
	}

	uploadPath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		logger.Error("保存文件失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "保存文件失败", err)
		return
	}

	logger.Info("Markdown文件上传成功",
		logger.String("filename", file.Filename),
		logger.Int64("size", file.Size),
		logger.String("path", uploadPath),
	)

	// 返回解析后的内容
	result := map[string]any{
		"title":    title,
		"content":  markdownContent,
		"summary":  summary,
		"filename": file.Filename,
		"path":     uploadPath,
		"size":     file.Size,
	}

	response.Success(c, result, "文件上传并解析成功")
}

// UploadImage 上传图片
// @Summary 上传图片
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 从表单获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件失败", err)
		return
	}

	// 调用上传服务
	imageURL, err := h.uploadService.UploadImage(file)
	if err != nil {
		logger.Error("图片上传失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// 返回图片URL
	result := map[string]any{
		"url":  imageURL,
		"name": file.Filename,
		"size": file.Size,
	}

	response.Success(c, result, "图片上传成功")
}

// parseTitle 从 Markdown 内容中解析标题
func parseTitle(content, defaultTitle string) string {
	// 尝试匹配第一个一级标题
	re := regexp.MustCompile(`^#\s+(.+)$`)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// 如果没有找到一级标题，使用文件名（去掉扩展名）
	if defaultTitle != "" {
		return strings.TrimSuffix(defaultTitle, filepath.Ext(defaultTitle))
	}

	return "未命名文章"
}

// parseSummary 从 Markdown 内容中解析摘要
func parseSummary(content string) string {
	lines := strings.Split(content, "\n")
	summaryLines := []string{}
	maxLines := 3

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过标题行
		if strings.HasPrefix(line, "#") {
			continue
		}
		// 跳过空行
		if line == "" {
			continue
		}
		// 跳过代码块
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			continue
		}
		// 跳过链接和图片
		if strings.HasPrefix(line, "[") || strings.HasPrefix(line, "!") {
			continue
		}

		// 添加到摘要
		summaryLines = append(summaryLines, line)
		if len(summaryLines) >= maxLines {
			break
		}
	}

	summary := strings.Join(summaryLines, " ")
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}

	return summary
}

// UploadWord 上传并解析 Word 文档
// UploadWord 上传并解析 Word 文档
// @Summary 上传Word文档
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Word文档"
// @Success 200 {object} response.Response
// @Router /api/v1/upload/word [post]
func (h *UploadHandler) UploadWord(c *gin.Context) {
	// 从表单获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件失败", err)
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".doc" && ext != ".docx" {
		response.Error(c, http.StatusBadRequest, "只支持 .doc 或 .docx 格式的文件", nil)
		return
	}

	// 验证文件大小（限制为 20MB）
	if file.Size > 20*1024*1024 {
		response.Error(c, http.StatusBadRequest, "文件大小不能超过 20MB", nil)
		return
	}

	// 存储文件到项目根目录下的 word_files 目录
	uploadDir, err := filepath.Abs("../word_files")
	if err != nil {
		logger.Error("获取绝对路径失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "获取路径失败", err)
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("创建目录失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "创建目录失败", err)
		return
	}

	uploadPath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		logger.Error("保存文件失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "保存文件失败", err)
		return
	}

	// 解析 Word 文档（包含图片）
	docContent, err := parseDocxFileWithImages(uploadPath)
	if err != nil {
		logger.Error("解析Word文档失败", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "解析Word文档失败", err)
		return
	}

	if docContent.Text == "" {
		logger.Error("解析Word文档失败: 内容为空")
		response.Error(c, http.StatusInternalServerError, "解析Word文档失败: 内容为空", nil)
		return
	}

	// 保存图片并获取 URL 映射
	imageURLs, err := saveWordImages(docContent.Images)
	if err != nil {
		logger.Error("保存Word图片失败", logger.Err("error", err))
		// 不返回错误，继续处理文本内容
	}

	// 获取服务器基础 URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	// 将相对路径转换为完整 URL
	fullImageURLs := make(map[string]string)
	for rId, url := range imageURLs {
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			fullImageURLs[rId] = url
		} else {
			fullImageURLs[rId] = baseURL + url
		}
	}

	// 提取标题
	title := parseWordTitle(file.Filename, docContent.Text)
	summary := parseWordSummary(docContent.Text)

	// 将文本内容转换为 HTML 格式，并处理图片占位符
	htmlContent := convertTextToHtmlWithImages(docContent.Text, fullImageURLs)

	logger.Info("解析完成",
		logger.String("title", title),
		logger.Int("textContentLength", len(docContent.Text)),
		logger.Int("htmlContentLength", len(htmlContent)),
		logger.String("summary", summary),
		logger.Int("imageCount", len(imageURLs)),
	)

	logger.Info("Word文件上传成功",
		logger.String("filename", file.Filename),
		logger.Int64("size", file.Size),
		logger.String("path", uploadPath),
	)

	// 返回解析后的内容
	result := map[string]any{
		"title":      title,
		"content":    htmlContent,
		"summary":    summary,
		"filename":   file.Filename,
		"path":       uploadPath,
		"size":       file.Size,
		"images":     fullImageURLs,  // 图片完整 URL 映射
		"imageCount": len(imageURLs), // 图片数量
	}

	response.Success(c, result, "文件上传并解析成功")
}

// parseWordTitle 从 Word 文档内容中解析标题
func parseWordTitle(filename, content string) string {
	// 使用文件名作为标题（去掉扩展名）
	title := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	// 如果内容不为空，尝试提取第一行作为标题
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) <= 100 {
			return line
		}
	}

	if title == "" {
		return "未命名文档"
	}

	return title
}

// parseWordSummary 从 Word 文档内容中解析摘要
func parseWordSummary(content string) string {
	lines := strings.Split(content, "\n")
	summaryLines := []string{}
	maxLines := 3

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行
		if line == "" {
			continue
		}

		// 添加到摘要
		summaryLines = append(summaryLines, line)
		if len(summaryLines) >= maxLines {
			break
		}
	}

	summary := strings.Join(summaryLines, " ")
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}

	return summary
}

// convertTextToHtml 将文本转换为 HTML（保留 Word 样式）
func convertTextToHtml(text string) string {
	lines := strings.Split(text, "\n")
	var htmlParts []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检测是否为标题
		if isTitleLine(line) {
			level := getTitleLevel(line)
			htmlParts = append(htmlParts, fmt.Sprintf("<h%d>%s</h%d>", level, strings.TrimSpace(line), level))
		} else {
			htmlParts = append(htmlParts, fmt.Sprintf("<p>%s</p>", line))
		}
	}

	return strings.Join(htmlParts, "\n")
}

// isTitleLine 判断是否为标题行
func isTitleLine(line string) bool {
	// 如果行很短且全大写，可能是标题
	if len(line) > 0 && len(line) < 50 {
		upper := strings.ToUpper(line)
		if line == upper && !strings.Contains(line, " ") {
			return true
		}
	}
	return false
}

// getTitleLevel 获取标题级别
func getTitleLevel(line string) int {
	// 简单的标题级别判断
	if len(line) < 10 {
		return 1
	}
	if len(line) < 20 {
		return 2
	}
	return 3
}

// parseDocxFile 解析 .docx 文件并提取文本内容
// .docx 文件实际上是一个 ZIP 压缩包，包含多个 XML 文件
func parseDocxFile(filePath string) (string, error) {
	// 打开 .docx 文件作为 ZIP 文件
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("打开 .docx 文件失败: %w", err)
	}
	defer r.Close()

	// .docx 文件的主要内容在 word/document.xml 中
	var docContent []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("读取 document.xml 失败: %w", err)
			}
			defer rc.Close()

			docContent, err = io.ReadAll(rc)
			if err != nil {
				return "", fmt.Errorf("读取 document.xml 内容失败: %w", err)
			}
			break
		}
	}

	if docContent == nil {
		return "", fmt.Errorf("未找到 word/document.xml 文件")
	}

	logger.Info("读取到 document.xml，大小", logger.Int("bytes", len(docContent)))

	// 定义 XML 结构体用于解析（Word XML 命名空间）
	type Text struct {
		Content string `xml:",chardata"`
	}

	type Run struct {
		Text Text `xml:"w:t"`
	}

	type Paragraph struct {
		Runs []Run `xml:"w:r"`
	}

	type Document struct {
		Body struct {
			Paragraphs []Paragraph `xml:"w:p"`
		} `xml:"w:body"`
	}

	var doc Document
	if err := xml.Unmarshal(docContent, &doc); err != nil {
		logger.Error("解析 XML 失败", logger.Err("error", err))
		return "", fmt.Errorf("解析 XML 失败: %w", err)
	}

	logger.Info("XML 解析成功，段落数量", logger.Int("count", len(doc.Body.Paragraphs)))

	// 提取所有文本
	var textBuilder strings.Builder
	textCount := 0
	for _, para := range doc.Body.Paragraphs {
		for _, run := range para.Runs {
			if run.Text.Content != "" {
				textBuilder.WriteString(run.Text.Content)
				textCount++
			}
		}
		textBuilder.WriteString("\n")
	}

	result := textBuilder.String()
	logger.Info("提取文本完成，文本长度", logger.Int("length", len(result)), logger.Int("textRuns", textCount))

	if result == "" {
		return "", fmt.Errorf("未能从 Word 文档中提取到任何文本内容")
	}

	return result, nil
}

// WordImageInfo 存储 Word 文档中的图片信息
type WordImageInfo struct {
	RelationshipID string // 关系 ID
	FileName       string // 图片文件名
	Data           []byte // 图片数据
	ContentType    string // 图片类型
}

// WordDocContent 存储 Word 文档解析后的内容
type WordDocContent struct {
	Text   string
	Images map[string]WordImageInfo // relationshipID -> imageInfo
}

// parseDocxFileWithImages 解析 Word 文档并提取文本和图片
func parseDocxFileWithImages(filePath string) (*WordDocContent, error) {
	// 打开 .docx 文件作为 ZIP 文件
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开 .docx 文件失败: %w", err)
	}
	defer r.Close()

	result := &WordDocContent{
		Images: make(map[string]WordImageInfo),
	}

	// 1. 读取 word/_rels/document.xml.rels 获取图片关系
	relationships := make(map[string]string) // rId -> target
	for _, f := range r.File {
		if f.Name == "word/_rels/document.xml.rels" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("读取 _rels 失败: %w", err)
			}
			defer rc.Close()

			relsContent, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("读取 _rels 内容失败: %w", err)
			}

			// 解析关系
			type Relationship struct {
				ID     string `xml:"Id,attr"`
				Type   string `xml:"Type,attr"`
				Target string `xml:"Target,attr"`
			}

			type Relationships struct {
				Rel []Relationship `xml:"Relationship"`
			}

			var rels Relationships
			if err := xml.Unmarshal(relsContent, &rels); err == nil {
				for _, rel := range rels.Rel {
					if strings.Contains(rel.Type, "image") {
						relationships[rel.ID] = rel.Target
					}
				}
			}
			break
		}
	}

	logger.Info("找到", logger.Int("relationships", len(relationships)))

	// 2. 提取图片数据
	for rId, target := range relationships {
		imagePath := "word/" + target
		for _, f := range r.File {
			if f.Name == imagePath {
				rc, err := f.Open()
				if err != nil {
					logger.Error("打开图片失败", logger.Err("error", err), logger.String("path", imagePath))
					continue
				}

				imageData, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					logger.Error("读取图片失败", logger.Err("error", err), logger.String("path", imagePath))
					continue
				}

				// 确定图片类型
				contentType := "image/jpeg"
				if strings.HasSuffix(target, ".png") {
					contentType = "image/png"
				} else if strings.HasSuffix(target, ".gif") {
					contentType = "image/gif"
				} else if strings.HasSuffix(target, ".svg") {
					contentType = "image/svg+xml"
				}

				result.Images[rId] = WordImageInfo{
					RelationshipID: rId,
					FileName:       filepath.Base(target),
					Data:           imageData,
					ContentType:    contentType,
				}
				logger.Info("提取图片成功", logger.String("rId", rId), logger.String("file", filepath.Base(target)), logger.Int("size", len(imageData)))
				break
			}
		}
	}

	// 3. 解析文档文本（并标记图片位置）
	var docContent []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("读取 document.xml 失败: %w", err)
			}
			defer rc.Close()

			docContent, err = io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("读取 document.xml 内容失败: %w", err)
			}
			break
		}
	}

	if docContent == nil {
		return nil, fmt.Errorf("未找到 word/document.xml 文件")
	}

	logger.Info("读取到 document.xml，大小", logger.Int("bytes", len(docContent)))

	// 4. 解析 XML 并提取文本
	xmlString := string(docContent)

	// 提取文本和图片的相对位置
	// 我们需要解析 <w:r> (run) 元素，每个 run 可能包含文本或图片
	// 按顺序提取，保持文档的原始结构
	type ContentElement struct {
		Type string // "text" 或 "image"
		Text string
		RId  string // 图片的 relationship ID
	}

	runRegex := regexp.MustCompile(`<w:r[^>]*>([\s\S]*?)</w:r>`)
	runMatches := runRegex.FindAllStringSubmatch(xmlString, -1)

	var contentElements []ContentElement

	// 正则表达式用于匹配文本和图片引用
	textRegex := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	imageRegex := regexp.MustCompile(`<a:blip[^>]*r:embed="([^"]+)"`)

	for _, runMatch := range runMatches {
		if len(runMatch) > 1 {
			runContent := runMatch[1]

			// 检查这个 run 中是否有图片
			imageMatches := imageRegex.FindStringSubmatch(runContent)
			if len(imageMatches) > 1 {
				contentElements = append(contentElements, ContentElement{
					Type: "image",
					RId:  imageMatches[1],
				})
			}

			// 检查这个 run 中是否有文本
			textMatches := textRegex.FindAllStringSubmatch(runContent, -1)
			for _, textMatch := range textMatches {
				if len(textMatch) > 1 {
					text := decodeXmlEntities(textMatch[1])
					if text != "" {
						contentElements = append(contentElements, ContentElement{
							Type: "text",
							Text: text,
						})
					}
				}
			}
		}
	}

	// 按顺序构建文本内容，插入图片占位符
	var textBuilder strings.Builder
	textCount := 0
	imageCount := 0

	for _, element := range contentElements {
		switch element.Type {
		case "text":
			textBuilder.WriteString(element.Text)
			textCount++
		case "image":
			// 在文本中插入图片占位符
			textBuilder.WriteString(fmt.Sprintf(" [IMG:%s] ", element.RId))
			imageCount++
		}
	}

	result.Text = textBuilder.String()

	// 检查文本中是否包含占位符
	placeholderCount := strings.Count(result.Text, "[IMG:")
	logger.Info("文本中的占位符数量", logger.Int("count", placeholderCount))
	logger.Info("提取文本完成，文本长度", logger.Int("length", len(result.Text)))

	if result.Text == "" {
		return nil, fmt.Errorf("未能从 Word 文档中提取到任何文本内容")
	}

	return result, nil
}

// saveWordImages 保存 Word 文档中的图片并返回图片 URL 映射
func saveWordImages(images map[string]WordImageInfo) (map[string]string, error) {
	imageURLs := make(map[string]string)

	// 创建图片保存目录（直接在 uploads/word_images 下）
	imageDir := "./uploads/word_images"
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建图片目录失败: %w", err)
	}

	for rId, img := range images {
		// 生成唯一文件名
		ext := filepath.Ext(img.FileName)
		if ext == "" {
			ext = ".jpg"
		}
		imageFileName := fmt.Sprintf("word_%s%s", generateRandomString(8), ext)
		imagePath := filepath.Join(imageDir, imageFileName)

		// 保存图片文件
		if err := os.WriteFile(imagePath, img.Data, 0644); err != nil {
			logger.Error("保存图片失败", logger.Err("error", err), logger.String("file", imageFileName))
			continue
		}

		// 生成访问 URL（匹配静态文件路由 /uploads）
		imageURLs[rId] = "/uploads/word_images/" + imageFileName
		logger.Info("保存图片成功",
			logger.String("rId", rId),
			logger.String("file", imageFileName),
			logger.String("url", imageURLs[rId]),
		)
	}

	return imageURLs, nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[generateRandomInt(len(charset))]
	}
	return string(b)
}

// generateRandomInt 生成随机数
func generateRandomInt(max int) int {
	return rand.Intn(max)
}

// decodeXmlEntities 解码 XML 实体
func decodeXmlEntities(text string) string {
	// 移除控制字符（ASCII 0-31，除了换行符和制表符）
	result := strings.Builder{}
	for _, r := range text {
		// 保留可打印字符和换行、制表符
		if r == '\n' || r == '\t' || r >= 32 {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// convertTextToHtmlWithImages 将文本转换为 HTML，并正确处理图片占位符
func convertTextToHtmlWithImages(text string, imageURLs map[string]string) string {
	lines := strings.Split(text, "\n")
	var htmlParts []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否包含图片占位符
		if strings.Contains(line, "[IMG:") {
			// 拆分行，将图片占位符作为独立元素
			parts := strings.Split(line, "[IMG:")

			for i, part := range parts {
				if i == 0 {
					// 第一部分可能是普通文本
					if part != "" {
						htmlParts = append(htmlParts, fmt.Sprintf("<p>%s</p>", part))
					}
				} else {
					// 后续部分格式为 "rIdX] 文本"
					endIndex := strings.Index(part, "]")
					if endIndex > 0 {
						rId := part[:endIndex]
						remainingText := part[endIndex+1:]

						// 插入图片
						if imageUrl, ok := imageURLs[rId]; ok {
							htmlParts = append(htmlParts, fmt.Sprintf(`<img src="%s" alt="文档图片" style="max-width:100%%"/>`, imageUrl))
						}

						// 如果有剩余文本，添加段落
						if remainingText != "" {
							htmlParts = append(htmlParts, fmt.Sprintf("<p>%s</p>", remainingText))
						}
					}
				}
			}
		} else {
			// 检测是否为标题
			if isTitleLine(line) {
				level := getTitleLevel(line)
				htmlParts = append(htmlParts, fmt.Sprintf("<h%d>%s</h%d>", level, strings.TrimSpace(line), level))
			} else {
				htmlParts = append(htmlParts, fmt.Sprintf("<p>%s</p>", line))
			}
		}
	}

	result := strings.Join(htmlParts, "\n")
	logger.Info("转换HTML完成", logger.String("html", result))
	return result
}
