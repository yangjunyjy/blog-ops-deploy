package request

type CreateSshConnectRequest struct {
	HostID    uint   `json:"host_id" binding:"required"`
	SessionID string `json:"session_id"`
}

type UploadFileRequest struct {
	SessionID   string `form:"session_id" binding:"required"`
	Path        string `form:"path" binding:"required"`
	FileName    string `form:"file_name" binding:"required"`
	ChunkIndex  int    `form:"chunk_index" binding:"required"`
	TotalChunks int    `form:"total_chunks" binding:"required"`
}

type UploadFileFormData struct {
	Path        string `form:"path" binding:"required"`
	FileName    string `form:"file_name" binding:"required"`
	ChunkIndex  int    `form:"chunk_index" binding:"required"`
	TotalChunks int    `form:"total_chunks" binding:"required"`
}
