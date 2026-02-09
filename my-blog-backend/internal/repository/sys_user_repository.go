package repository

import (
	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
)

// SysUserRepository 系统用户仓储接口
type SysUserRepository interface {
	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*models.SysUser, error)
	// FindByID 根据ID查找用户
	FindByID(id uint64) (*models.SysUser, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*models.SysUser, error)
	// List 查询用户列表
	List(query *request.SearchSysUserQueryRequest) ([]*response.UserListResponse, int64, error)
	// Create 创建用户
	Create(user *models.SysUser) error
	// Update 更新用户
	Update(user *models.SysUser) error
	// Delete 删除用户（软删除）
	Delete(id uint64) error
	// GetUserRoles 获取用户的角色列表
	GetUserRoles(userID uint64) ([]*models.SysRole, error)
	// GetUserPosts 获取用户的岗位列表
	GetUserPosts(userID uint64) ([]*models.SysPost, error)
	// AssignRoles 分配角色给用户
	AssignRoles(userID uint64, roleIDs []uint64) error
	// AssignPosts 分配岗位给用户
	AssignPosts(userID uint64, postIDs []uint64) error
	// UpdateStatus 更新用户状态
	UpdateStatus(id uint64, status models.Status) error
	// ResetPassword 重置密码
	ResetPassword(id uint64, newPassword string) error
}
