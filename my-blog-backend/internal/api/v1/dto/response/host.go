package response

type HostResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type HostListResponse struct {
	Total int64             `json:"total"`
	Items []HostResponse   `json:"items"`
}

type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
