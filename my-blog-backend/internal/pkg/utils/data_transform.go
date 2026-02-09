package utils

import "my-blog-backend/internal/models"

func ExtractRolePerms(user *models.SysUser) []string {
	if user == nil || len(user.Roles) == 0 {
		return []string{}
	}

	// 使用map直接去重，无需中间结构
	permSet := make(map[string]struct{})

	// 单次遍历，直接收集权限码
	for _, role := range user.Roles {
		if role == nil {
			continue
		}
		for _, menu := range role.Menus {
			if menu != nil && menu.MenuType == 3 {
				permSet[menu.MenuCode] = struct{}{}
			}
		}
	}

	// 将map转换为切片
	result := make([]string, 0, len(permSet))
	for permCode := range permSet {
		result = append(result, permCode)
	}

	return result
}
