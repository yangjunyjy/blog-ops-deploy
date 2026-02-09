package api

import (
	"fmt"
	"io"
	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/services"
	"my-blog-backend/internal/ssh"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	// CHUNK_SIZE 每个分片的大小（5MB）
	CHUNK_SIZE = 5 * 1024 * 1024
)

var typeMap = map[string]string{
	".jpg":  "图片",
	".jpeg": "图片",
	".png":  "图片",
	".gif":  "图片",
	".bmp":  "图片",
	".svg":  "图片",
	".webp": "图片",
	".ico":  "图片",

	".pdf":  "PDF文档",
	".doc":  "Word文档",
	".docx": "Word文档",
	".xls":  "Excel表格",
	".xlsx": "Excel表格",
	".ppt":  "PowerPoint演示",
	".pptx": "PowerPoint演示",

	".txt": "文本文件",
	".md":  "Markdown文档",
	".rtf": "富文本",
	".csv": "CSV文件",

	".zip": "压缩文件",
	".rar": "压缩文件",
	".7z":  "压缩文件",
	".tar": "压缩文件",
	".gz":  "压缩文件",
	".bz2": "压缩文件",

	".mp4":  "视频文件",
	".avi":  "视频文件",
	".mov":  "视频文件",
	".wmv":  "视频文件",
	".flv":  "视频文件",
	".mkv":  "视频文件",
	".webm": "视频文件",

	".mp3":  "音频文件",
	".wav":  "音频文件",
	".flac": "音频文件",
	".aac":  "音频文件",
	".m4a":  "音频文件",
	".ogg":  "音频文件",

	".exe": "可执行文件",
	".msi": "安装程序",
	".apk": "Android应用",
	".dmg": "Mac安装包",
	".deb": "Debian包",
	".rpm": "RPM包",

	".json": "JSON数据",
	".xml":  "XML数据",
	".yaml": "YAML配置",
	".yml":  "YAML配置",

	".html": "HTML文件",
	".css":  "CSS样式",
	".js":   "JavaScript脚本",
	".ts":   "TypeScript脚本",
	".vue":  "Vue组件",

	".java": "Java源码",
	".py":   "Python脚本",
	".go":   "Go源码",
	".c":    "C源码",
	".cpp":  "C++源码",
	".h":    "C头文件",
	".hpp":  "C++头文件",
	".rs":   "Rust源码",
	".php":  "PHP脚本",
	".rb":   "Ruby脚本",
	".sh":   "Shell脚本",
	".bat":  "批处理文件",
	".ps1":  "PowerShell脚本",

	".sql":    "SQL脚本",
	".db":     "数据库文件",
	".sqlite": "数据库文件",

	".log":    "日志文件",
	".conf":   "配置文件",
	".config": "配置文件",
	".env":    "环境变量",

	".ttf":   "字体文件",
	".otf":   "字体文件",
	".woff":  "字体文件",
	".woff2": "字体文件",
	".eot":   "字体文件",
}

type ChunkManager struct {
	mu       sync.Mutex
	chunkMap map[string]map[int]bool // filepath -> 对应切片
}

var chunkManager = ChunkManager{
	chunkMap: make(map[string]map[int]bool),
}

type SshFileHandler struct {
	hostService *services.HostService
	pool        *ssh.Pool
	sessions    map[string]*ssh.Session // 共享 SshHandler 的 sessions
	mu          sync.RWMutex
}

func NewSshFileHandler(hostService *services.HostService, pool *ssh.Pool, sessions map[string]*ssh.Session) *SshFileHandler {
	return &SshFileHandler{
		hostService: hostService,
		pool:        pool,
		sessions:    sessions,
	}
}

// UploadFile 处理文件分片上传
func (h *SshFileHandler) UploadFile(c *gin.Context) {
	// 获取参数
	sessionID := c.Query("session_id")
	if sessionID == "" {
		response.Error(c, http.StatusBadRequest, "session_id不能为空", fmt.Errorf("session_id不能为空"))
		return
	}

	path2 := c.PostForm("path")
	fileName := c.PostForm("file_name")
	chunkIndexStr := c.PostForm("chunk_index")
	totalChunksStr := c.PostForm("total_chunks")

	// 参数验证
	if path2 == "" || fileName == "" || chunkIndexStr == "" || totalChunksStr == "" {
		response.Error(c, http.StatusBadRequest, "FormData参数缺失", fmt.Errorf("FormData参数缺失"))
		return
	}

	chunkIndex, err1 := strconv.Atoi(chunkIndexStr)
	totalChunks, err2 := strconv.Atoi(totalChunksStr)
	if err1 != nil || err2 != nil {
		response.Error(c, http.StatusBadRequest, "数字参数格式错误", fmt.Errorf("数字参数格式错误"))
		return
	}

	req := request.UploadFileRequest{
		SessionID:   sessionID,
		Path:        path2,
		FileName:    fileName,
		ChunkIndex:  chunkIndex,
		TotalChunks: totalChunks,
	}

	// 获取 SSH 会话和 SFTP 客户端
	h.mu.RLock()
	session, ok := h.sessions[req.SessionID]
	h.mu.RUnlock()
	if !ok {
		response.Error(c, http.StatusBadRequest, "无效的session_id", fmt.Errorf("无效的session_id: %s", req.SessionID))
		return
	}

	client, err := session.Client.GetSFTP()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取sftp连接失败", fmt.Errorf("获取sftp时发生错误: %v", err))
		return
	}

	// 构造目标文件路径
	targetFile := path.Join(req.Path, req.FileName)

	// 初始化切片管理器
	chunkManager.mu.Lock()
	if _, exists := chunkManager.chunkMap[targetFile]; !exists {
		chunkManager.chunkMap[targetFile] = make(map[int]bool)
	}
	chunkManager.mu.Unlock()

	// 创建临时文件路径（根目录特殊处理）
	var tempChunkPath string
	if req.Path == "/" {
		tempChunkPath = "/" + fmt.Sprintf(".tmp_%s_chunk_%d", req.FileName, req.ChunkIndex)
	} else {
		tempChunkPath = path.Join(req.Path, fmt.Sprintf(".tmp_%s_chunk_%d", req.FileName, req.ChunkIndex))
	}

	// 从 FormData 获取文件分片数据
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取文件数据失败", fmt.Errorf("获取文件数据发生错误: %v", err))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "打开文件失败", fmt.Errorf("打开文件发生错误: %v", err))
		return
	}
	defer file.Close()

	// 读取文件内容到内存
	chunkData := make([]byte, 0)
	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunkData = append(chunkData, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("读取文件数据失败", logger.Err("error", err))
			break
		}
	}

	logger.Info("接收到分片数据",
		logger.Int("data_len", len(chunkData)),
		logger.Int("chunk", req.ChunkIndex),
		logger.Int("expected_size", int(req.ChunkIndex+1)*CHUNK_SIZE))

	if len(chunkData) == 0 {
		logger.Error("接收到空数据", logger.Int("chunk", req.ChunkIndex))
		response.Error(c, http.StatusBadRequest, "接收到空数据", fmt.Errorf("分片 %d 数据为空", req.ChunkIndex))
		return
	}

	// 创建临时文件并写入分片数据
	tempFile, err := client.Create(tempChunkPath)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建临时文件失败", fmt.Errorf("创建临时文件失败: %v", err))
		return
	}
	defer tempFile.Close()

	written, err := tempFile.Write(chunkData)
	if err != nil {
		logger.Error("写流式数据到临时文件出错", logger.Err("error", err))
		response.Error(c, http.StatusInternalServerError, "写数据出错", fmt.Errorf("写数据到临时文件出错 %v", err))
		return
	}

	logger.Info("分片数据写入临时文件",
		logger.Int64("bytes", int64(written)),
		logger.Int("chunk", req.ChunkIndex),
		logger.Int("data_len", len(chunkData)),
		logger.String("file", tempChunkPath))

	// 更新切片管理器
	chunkManager.mu.Lock()
	chunkManager.chunkMap[targetFile][req.ChunkIndex] = true

	// 统计已完成分片数
	fileCount := 0
	for _, completed := range chunkManager.chunkMap[targetFile] {
		if completed {
			fileCount++
		}
	}

	process := float64(fileCount) / float64(req.TotalChunks) * 100
	chunkManager.mu.Unlock()

	// 检查是否所有分片完成
	if fileCount >= req.TotalChunks {
		go h.mergeChunksAsync(req.SessionID, req.Path, req.FileName, req.TotalChunks)
		response.Success(c, map[string]interface{}{
			"message":       "所有切片上传完成，正在合并...",
			"completed":     true,
			"total_chunks":  req.TotalChunks,
			"chunk_index":   req.ChunkIndex,
			"progress":      100.0,
			"bytes_written": written,
		}, "上传完成")
		return
	}

	response.Success(c, map[string]interface{}{
		"message":       "切片上传成功",
		"completed":     false,
		"total_chunks":  req.TotalChunks,
		"chunk_index":   req.ChunkIndex,
		"progress":      process,
		"bytes_written": written,
	}, "上传成功")
}

// mergeChunksAsync 异步合并所有分片为最终文件
func (h *SshFileHandler) mergeChunksAsync(sessionID, filePath, name string, count int) {
	// 获取 SSH 会话和 SFTP 客户端
	h.mu.Lock()
	session, ok := h.sessions[sessionID]
	h.mu.Unlock()
	if !ok {
		logger.Error("异步合并时，获取session出错")
		return
	}

	sftpClient, err := session.Client.GetSFTP()
	if err != nil {
		logger.Error("获取sftp客户端出错: %v", logger.Err("error", err))
		return
	}

	// 构造目标文件路径（根目录特殊处理）
	var targetPath string
	if filePath == "/" {
		targetPath = "/" + name
	} else {
		targetPath = path.Join(filePath, name)
	}

	// 创建目标文件
	targetFile, err := sftpClient.Create(targetPath)
	if err != nil {
		logger.Error("创建目标文件失败", logger.Err("error", err))
		return
	}
	defer targetFile.Close()

	// 按顺序合并所有临时分片文件
	for i := 0; i < count; i++ {
		var tempChunkPath string
		if filePath == "/" {
			tempChunkPath = "/" + fmt.Sprintf(".tmp_%s_chunk_%d", name, i)
		} else {
			tempChunkPath = path.Join(filePath, fmt.Sprintf(".tmp_%s_chunk_%d", name, i))
		}

		// 打开临时分片文件
		tempFile, err := sftpClient.Open(tempChunkPath)
		if err != nil {
			logger.Info("临时文件不存在，跳过", logger.String("filename", tempChunkPath))
			continue
		}

		_, err = tempFile.Stat()
		if err != nil {
			logger.Error("获取临时文件信息失败", logger.Err("error", err), logger.Int("chunk", i))
			tempFile.Close()
			continue
		}

		// 复制分片数据到目标文件
		written, err := io.Copy(targetFile, tempFile)
		if err != nil {
			logger.Error("复制分片数据失败", logger.Err("error", err), logger.Int("chunk", i))
		} else {
			logger.Info("复制分片成功", logger.Int64("bytes", written), logger.Int("chunk", i))
		}

		tempFile.Close()

		// 删除临时分片文件
		if err := sftpClient.Remove(tempChunkPath); err != nil {
			logger.Warn("删除临时分片文件失败", logger.Err("error", err), logger.String("file", tempChunkPath))
		}
	}

	// 清理切片管理器
	h.mu.Lock()
	delete(chunkManager.chunkMap, targetPath)
	h.mu.Unlock()

	logger.Info("文件合并完成", logger.String("file", targetPath))
}

func (h *SshFileHandler) List(c *gin.Context) {
	// 获取参数
	sessionID := c.Query("session_id")
	if sessionID == "" {
		response.Error(c, http.StatusBadRequest, "session_id不能为空", fmt.Errorf("session_id不能为空"))
		return
	}
	path := c.Query("path")
	if path == "" {
		response.Error(c, http.StatusBadRequest, "请求路径不能为空", nil)
		return
	}

	// 获取 SSH 会话
	h.mu.RLock()
	session, ok := h.sessions[sessionID]
	h.mu.RUnlock()
	if !ok {
		response.Error(c, http.StatusBadRequest, "无效的session_id", fmt.Errorf("无效的session_id: %s", sessionID))
		return
	}

	// 获取 SFTP 客户端
	sftpClient, err := session.Client.GetSFTP()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取sftp连接出错",
			fmt.Errorf("获取sftp连接出错: %v", err))
		return
	}

	// 读取目录内容
	entries, err := sftpClient.ReadDir(path)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取路径信息失败",
			fmt.Errorf("获取路径信息失败: %v", err))
		return
	}

	// 构造文件列表
	fileInfos := make([]response.FileInfo, 0, len(entries))

	// 如果不是根目录，添加 ".." 目录用于返回上一级
	if path != "/" && path != "" {
		parentPath := getParentPath(path)
		fileInfos = append(fileInfos, response.FileInfo{
			Name:    "..",
			Path:    parentPath,
			IsDir:   true,
			Size:    0,
			ModTime: "",
			Kind:    "上级目录",
		})
	}

	// 添加目录中的文件和文件夹
	for _, entry := range entries {
		fileInfo := response.FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    entry.Size(),
			ModTime: entry.ModTime().Format("2006-01-02 15:04:05"),
			Kind:    getFileKind(entry),
		}

		// 构造完整路径
		if path == "/" {
			fileInfo.Path = "/" + entry.Name()
		} else {
			fileInfo.Path = path + "/" + entry.Name()
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	// 返回结果
	res := response.PathListInfoResponse{
		Path:  path,
		Files: fileInfos,
	}

	response.Success(c, res, "获取文件列表成功")
}

// getParentPath 获取父目录路径
func getParentPath(path string) string {
	if path == "" || path == "/" {
		return "/"
	}

	// 去除末尾的斜杠
	path = strings.TrimSuffix(path, "/")

	// 找到最后一个斜杠的位置
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash <= 0 {
		return "/"
	}

	return path[:lastSlash]
}

func (h *SshFileHandler) DownLoadFile(c *gin.Context) {

}

// getFileKind 根据文件信息获取文件类型描述
func getFileKind(entry os.FileInfo) string {
	if entry.IsDir() {
		return "目录"
	}

	name := entry.Name()
	ext := ""
	if idx := strings.LastIndex(name, "."); idx > 0 {
		ext = strings.ToLower(name[idx:])
	}

	// 文件类型映射

	if kind, ok := typeMap[ext]; ok {
		return kind
	}
	return "文件"
}
