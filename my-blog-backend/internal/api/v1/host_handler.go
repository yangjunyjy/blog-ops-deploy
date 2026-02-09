package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"my-blog-backend/internal/api/v1/dto/request"
	dtoResponse "my-blog-backend/internal/api/v1/dto/response"
	"my-blog-backend/internal/services"
)

type HostHandler struct {
	hostService *services.HostService
}

func NewHostHandler(hostService *services.HostService) *HostHandler {
	return &HostHandler{hostService: hostService}
}

// CreateHost 创建主机
// @Summary 创建主机
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param request body request.CreateHostRequest true "主机信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/hosts [post]
func (h *HostHandler) CreateHost(c *gin.Context) {
	var req request.CreateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, 400, "请求参数错误", err)
		return
	}

	if err := h.hostService.CreateHost(&req); err != nil {
		dtoResponse.Error(c, 500, "创建主机失败", err)
		return
	}

	dtoResponse.Success(c, nil, "创建成功")
}

// UpdateHost 更新主机
// @Summary 更新主机
// @Tags 主机管理
// @Accept json
// @Produce json
// @Param request body request.UpdateHostRequest true "主机信息"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/hosts [put]
func (h *HostHandler) UpdateHost(c *gin.Context) {
	var req request.UpdateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dtoResponse.Error(c, 400, "请求参数错误", err)
		return
	}

	if err := h.hostService.UpdateHost(&req); err != nil {
		dtoResponse.Error(c, 500, "更新主机失败", err)
		return
	}

	dtoResponse.Success(c, nil, "更新成功")
}

// DeleteHost 删除主机
// @Summary 删除主机
// @Tags 主机管理
// @Param id path int true "主机ID"
// @Success 200 {object} dtoResponse.Response
// @Router /api/v1/hosts/{id} [delete]
func (h *HostHandler) DeleteHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtoResponse.Error(c, 400, "无效的主机ID", err)
		return
	}

	if err := h.hostService.DeleteHost(uint(id)); err != nil {
		dtoResponse.Error(c, 500, "删除主机失败", err)
		return
	}

	dtoResponse.Success(c, nil, "删除成功")
}

// GetHost 获取主机详情
// @Summary 获取主机详情
// @Tags 主机管理
// @Param id path int true "主机ID"
// @Success 200 {object} dtoResponse.Response{data=dtoResponse.HostResponse}
// @Router /api/v1/hosts/{id} [get]
func (h *HostHandler) GetHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtoResponse.Error(c, 400, "无效的主机ID", err)
		return
	}

	host, err := h.hostService.GetHost(uint(id))
	if err != nil {
		dtoResponse.Error(c, 500, "获取主机失败", err)
		return
	}

	dtoResponse.Success(c, host, "获取成功")
}

// ListHosts 主机列表
// @Summary 主机列表
// @Tags 主机管理
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param name query string false "主机名称"
// @Param address query string false "主机地址"
// @Param type query string false "认证类型"
// @Param status query string false "状态"
// @Success 200 {object} dtoResponse.Response{data=dtoResponse.HostListResponse}
// @Router /api/v1/hosts [get]
func (h *HostHandler) ListHosts(c *gin.Context) {
	var req request.ListHostRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		dtoResponse.Error(c, 400, "请求参数错误", err)
		return
	}

	list, err := h.hostService.ListHosts(&req)
	if err != nil {
		dtoResponse.Error(c, 500, "获取主机列表失败", err)
		return
	}

	dtoResponse.Success(c, list, "获取成功")
}

// GetAllHosts 获取所有主机（用于下拉选择）
// @Summary 获取所有主机
// @Tags 主机管理
// @Success 200 {object} dtoResponse.Response{data=[]dtoResponse.HostResponse}
// @Router /api/v1/hosts/all [get]
func (h *HostHandler) GetAllHosts(c *gin.Context) {
	hosts, err := h.hostService.GetAllHosts()
	if err != nil {
		dtoResponse.Error(c, 500, "获取主机列表失败", err)
		return
	}

	dtoResponse.Success(c, hosts, "获取成功")
}

// TestConnection 测试连接
// @Summary 测试连接
// @Tags 主机管理
// @Param id path int true "主机ID"
// @Success 200 {object} dtoResponse.Response{data=dtoResponse.TestConnectionResponse}
// @Router /api/v1/hosts/{id}/test [post]
func (h *HostHandler) TestConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		dtoResponse.Error(c, 400, "无效的主机ID", err)
		return
	}

	result, err := h.hostService.TestConnection(uint(id))
	if err != nil {
		dtoResponse.Error(c, 500, "测试连接失败", err)
		return
	}

	dtoResponse.Success(c, result, "测试完成")
}
