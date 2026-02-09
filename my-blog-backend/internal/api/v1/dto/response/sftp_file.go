package response

// FileInfo 文件信息
type FileInfo struct {
	Name    string    `json:"name"`    // 文件名
	Path    string    `json:"path"`    // 完整路径
	IsDir   bool      `json:"isDir"`   // 是否为目录
	Size    int64     `json:"size"`    // 文件大小（字节）
	ModTime string    `json:"modTime"` // 修改时间（格式化字符串）
	Kind    string    `json:"kind"`    // 文件类型描述
}

// PathListInfoResponse 文件列表响应
type PathListInfoResponse struct {
	Path  string      `json:"path"`  // 当前路径
	Files []FileInfo  `json:"files"` // 文件列表
}

