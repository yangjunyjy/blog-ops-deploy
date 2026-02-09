package mysql

import (
	"gorm.io/gorm"

	"my-blog-backend/internal/api/v1/dto/request"
	"my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/models"
	"my-blog-backend/internal/pkg/logger"
	_ "my-blog-backend/internal/pkg/logger"
	"my-blog-backend/internal/repository"
)

type SysUserRepositoryImpl struct {
	db *gorm.DB
}

func NewSysUserRepositoryImpl(db *gorm.DB) repository.SysUserRepository {
	return &SysUserRepositoryImpl{db: db}
}

func (r *SysUserRepositoryImpl) FindByUsername(username string) (*models.SysUser, error) {
	var user models.SysUser
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SysUserRepositoryImpl) FindByID(id uint64) (*models.SysUser, error) {
	var user models.SysUser
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SysUserRepositoryImpl) FindByEmail(email string) (*models.SysUser, error) {
	var user models.SysUser
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SysUserRepositoryImpl) List(query *request.SearchSysUserQueryRequest) ([]*response.UserListResponse, int64, error) {
	var users []*models.SysUser
	var total int64

	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize

	// 构建基础查询
	db := r.db.Model(&models.SysUser{})

	// 应用查询条件
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.Status > 0 { // 假设0表示查询所有状态
		db = db.Where("status = ?", query.Status)
	}

	// 统计总数（必须在分页之前）
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，通过中间表 sys_user_role 预加载角色信息
	err := db.Preload("Roles").
		Offset(offset).
		Limit(query.PageSize).
		Order("id DESC").
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// 将 models.SysUser 转换为 response.UserListResponse
	userResponses := make([]*response.UserListResponse, len(users))
	for i, user := range users {
		userResponses[i] = &response.UserListResponse{
			ID:         user.ID,
			Username:   user.Username,
			Password:   user.Password, // 注意：通常不返回密码，这里根据你的结构体保留
			Nickname:   user.Nickname,
			Avatar:     user.Avatar,
			Email:      user.Email,
			Gender:     user.Gender,
			Status:     user.Status,
			IsAdmin:    user.IsAdmin,
			CreateTime: *user.CreateTime,
			Roles:      user.Roles, // 通过 Preload 已加载的角色信息
		}
	}

	return userResponses, total, nil
}

func (r *SysUserRepositoryImpl) Create(user *models.SysUser) error {
	return r.db.Create(user).Error
}

func (r *SysUserRepositoryImpl) Update(user *models.SysUser) error {
	return r.db.Save(user).Error
}

func (r *SysUserRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&models.SysUser{}, id).Error
}

func (r *SysUserRepositoryImpl) GetUserRoles(userID uint64) ([]*models.SysRole, error) {
	var roles []*models.SysRole

	// 方法2：通过关联表直接查询用户的所有角色
	err := r.db.
		Joins("JOIN sys_user_role sur ON sur.role_id = sys_role.id").
		Where("sur.user_id = ?", userID).
		Preload("Menus"). // 预加载每个角色的菜单
		Find(&roles).Error

	if err != nil {
		return nil, err
	}
	// logger.Info("repo层查询到的角色信息", logger.Any("roles", roles))
	return roles, nil
}

func (r *SysUserRepositoryImpl) GetUserPosts(userID uint64) ([]*models.SysPost, error) {
	var posts []*models.SysPost
	err := r.db.Table("sys_post").
		Joins("INNER JOIN sys_user_post ON sys_post.id = sys_user_post.post_id").
		Where("sys_user_post.user_id = ?", userID).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *SysUserRepositoryImpl) AssignRoles(userID uint64, roleIDs []uint64) error {
	logger.Info("分配角色", logger.Any("userID", userID), logger.Any("roleIDs", roleIDs))

	// 先删除现有角色
	if err := r.db.Where("user_id = ?", userID).Delete(&models.SysUserRole{}).Error; err != nil {
		logger.Error("删除用户角色失败", logger.Err("error", err))
		return err
	}

	// 批量插入新角色
	if len(roleIDs) == 0 {
		logger.Info("角色ID列表为空,跳过插入")
		return nil
	}

	var userRoles []models.SysUserRole
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, models.SysUserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	logger.Info("准备插入用户角色", logger.Any("count", len(userRoles)), logger.Any("data", userRoles))
	err := r.db.CreateInBatches(userRoles, 100).Error
	if err != nil {
		logger.Error("插入用户角色失败", logger.Err("error", err))
	}
	return err
}

func (r *SysUserRepositoryImpl) AssignPosts(userID uint64, postIDs []uint64) error {
	// 先删除现有岗位
	if err := r.db.Where("user_id = ?", userID).Delete(&models.SysUserPost{}).Error; err != nil {
		return err
	}

	// 批量插入新岗位
	if len(postIDs) == 0 {
		return nil
	}

	var userPosts []models.SysUserPost
	for _, postID := range postIDs {
		userPosts = append(userPosts, models.SysUserPost{
			UserID: userID,
			PostID: postID,
		})
	}

	return r.db.CreateInBatches(userPosts, 100).Error
}

func (r *SysUserRepositoryImpl) UpdateStatus(id uint64, status models.Status) error {
	return r.db.Model(&models.SysUser{}).Where("id = ?", id).Update("status", status).Error
}

func (r *SysUserRepositoryImpl) ResetPassword(id uint64, newPassword string) error {
	return r.db.Model(&models.SysUser{}).Where("id = ?", id).Update("password", newPassword).Error
}
