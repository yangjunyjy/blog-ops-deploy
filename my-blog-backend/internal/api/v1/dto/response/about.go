package response

// AboutInfo 关于信息
type AboutInfo struct {
	Name        string   `json:"name"`
	Avatar      string   `json:"avatar"`
	Bio         string   `json:"bio"`
	Email       string   `json:"email"`
	Github      string   `json:"github"`
	Website     string   `json:"website"`
	Skills      []string `json:"skills"`
	Description string   `json:"description"`
}
