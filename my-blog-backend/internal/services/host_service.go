package services

import (
	"context"
	"fmt"
	"time"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
	opsModel "my-blog-backend/internal/models/opsModel"
	"my-blog-backend/internal/repository"
	"my-blog-backend/internal/ssh"
)

type HostService struct {
	hostRepo repository.HostRepository
	sshPool  *ssh.Pool
}

func NewHostService(hostRepo repository.HostRepository, sshPool *ssh.Pool) *HostService {
	return &HostService{
		hostRepo: hostRepo,
		sshPool:  sshPool,
	}
}

// CreateHost 创建主机
func (s *HostService) CreateHost(req *request.CreateHostRequest) error {
	// 转换认证类型
	var sshType opsModel.SshType
	switch req.Type {
	case "password":
		sshType = opsModel.Pwd
	case "key":
		sshType = opsModel.Key
	default:
		return fmt.Errorf("无效的认证类型")
	}

	// 设置默认状态
	var status models.Status = models.StatusEnabled
	if req.Status != "" {
		if req.Status == "inactive" {
			status = models.StatusDisabled
		}
	}

	host := &opsModel.RemoteHost{
		Name:      req.Name,
		Address:   req.Address,
		Port:      int64(req.Port),
		Username:  req.Username,
		Password:  req.Password,
		SecretKey: req.SecretKey,
		Type:      sshType,
		Status:    status,
	}

	return s.hostRepo.Create(host)
}

// UpdateHost 更新主机
func (s *HostService) UpdateHost(req *request.UpdateHostRequest) error {
	// 检查主机是否存在
	host, err := s.hostRepo.GetByID(req.ID)
	if err != nil {
		return fmt.Errorf("主机不存在")
	}

	// 转换认证类型
	var sshType opsModel.SshType
	switch req.Type {
	case "password":
		sshType = opsModel.Pwd
	case "key":
		sshType = opsModel.Key
	default:
		return fmt.Errorf("无效的认证类型")
	}

	// 更新字段
	host.Name = req.Name
	host.Address = req.Address
	host.Port = int64(req.Port)
	host.Username = req.Username
	host.Password = req.Password
	host.SecretKey = req.SecretKey
	host.Type = sshType

	// 如果提供了状态，则更新状态
	if req.Status != "" {
		if req.Status == "inactive" {
			host.Status = models.StatusDisabled
		} else {
			host.Status = models.StatusEnabled
		}
	}

	return s.hostRepo.Update(host)
}

// DeleteHost 删除主机
func (s *HostService) DeleteHost(id uint) error {
	// 检查主机是否存在
	_, err := s.hostRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("主机不存在")
	}

	// TODO: 检查是否有活跃的会话

	return s.hostRepo.Delete(id)
}

// GetHost 获取主机详情
func (s *HostService) GetHost(id uint) (*response.HostResponse, error) {
	host, err := s.hostRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("主机不存在")
	}

	return s.toHostResponse(host), nil
}

// ListHosts 主机列表
func (s *HostService) ListHosts(req *request.ListHostRequest) (*response.HostListResponse, error) {
	hosts, total, err := s.hostRepo.List(req.Page, req.PageSize, req.Name, req.Address, req.Type, req.Status)
	if err != nil {
		return nil, err
	}

	items := make([]response.HostResponse, len(hosts))
	for i, host := range hosts {
		items[i] = *s.toHostResponse(host)
	}

	return &response.HostListResponse{
		Total: total,
		Items: items,
	}, nil
}

// GetAllHosts 获取所有主机（用于下拉选择）
func (s *HostService) GetAllHosts() ([]*response.HostResponse, error) {
	hosts, err := s.hostRepo.GetAll()
	if err != nil {
		return nil, err
	}

	items := make([]*response.HostResponse, len(hosts))
	for i, host := range hosts {
		items[i] = s.toHostResponse(host)
	}

	return items, nil
}

// TestConnection 测试连接
func (s *HostService) TestConnection(id uint) (*response.TestConnectionResponse, error) {
	host, err := s.hostRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("主机不存在")
	}

	// 构建 SSH 配置
	var authType ssh.AuthType
	var key []byte

	switch host.Type {
	case opsModel.Pwd:
		authType = ssh.AuthTypePassword
	case opsModel.Key:
		authType = ssh.AuthTypeKey
		key = []byte(host.SecretKey)
	}

	cfg := &ssh.Config{
		Host:     host.Address,
		Port:     uint(host.Port),
		Username: host.Username,
		Password: host.Password,
		Key:      key,
		AuthType: authType,
		Timeout:  10 * time.Second,
	}

	// 创建 SSH 客户端测试连接
	opts := []ssh.Option{
		ssh.WithHost(cfg.Host),
		ssh.WithPort(cfg.Port),
		ssh.WithUsername(cfg.Username),
		ssh.WithTimeout(cfg.Timeout),
	}

	switch cfg.AuthType {
	case ssh.AuthTypePassword:
		if cfg.Password != "" {
			opts = append(opts, ssh.WithPassword(cfg.Password))
		}
	case ssh.AuthTypeKey:
		if len(cfg.Key) > 0 {
			opts = append(opts, ssh.WithKey(cfg.Key))
		}
	case ssh.AuthTypeBoth:
		if cfg.Password != "" && len(cfg.Key) > 0 {
			opts = append(opts, ssh.WithBoth(cfg.Password, cfg.Key))
		}
	}

	client, err := ssh.NewClient(context.Background(), opts...)
	if err != nil {
		return &response.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}

	client.Close()

	return &response.TestConnectionResponse{
		Success: true,
		Message: "连接成功",
	}, nil
}

// GetSSHConfig 获取 SSH 配置（用于 WebSocket 连接）
func (s *HostService) GetSSHConfig(hostID uint) (*ssh.Config, error) {
	host, err := s.hostRepo.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("主机不存在")
	}

	var authType ssh.AuthType
	var key []byte

	switch host.Type {
	case opsModel.Pwd:
		authType = ssh.AuthTypePassword
	case opsModel.Key:
		authType = ssh.AuthTypeKey
		key = []byte(host.SecretKey)
	}

	return &ssh.Config{
		Host:     host.Address,
		Port:     uint(host.Port),
		Username: host.Username,
		Password: host.Password,
		Key:      key,
		AuthType: authType,
		Timeout:  30 * time.Second,
	}, nil
}

// toHostResponse 转换为响应对象
func (s *HostService) toHostResponse(host *opsModel.RemoteHost) *response.HostResponse {
	var hostType string
	switch host.Type {
	case opsModel.Pwd:
		hostType = "password"
	case opsModel.Key:
		hostType = "key"
	}

	var status string
	if host.Status == models.StatusEnabled {
		status = "active"
	} else {
		status = "inactive"
	}

	return &response.HostResponse{
		ID:        host.ID,
		Name:      host.Name,
		Address:   host.Address,
		Port:      int(host.Port),
		Username:  host.Username,
		Type:      hostType,
		Status:    status,
		CreatedAt: host.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: host.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
