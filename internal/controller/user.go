package controller

import (
	"github.com/gin-gonic/gin"

	"cinexus/internal/service"
	"cinexus/pkg/response"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: service.UserService{},
	}
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	resp, err := c.userService.Login(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "登录成功", resp)
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	err := c.userService.Register(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "注册成功", nil)
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	user, err := c.userService.GetUserByID(userID.(uint))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, user)
}

// UpdateUserInfo 更新用户信息
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req service.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	err := c.userService.UpdateUser(userID.(uint), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// UpdatePassword 更新用户密码
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req service.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	err := c.userService.UpdatePassword(userID.(uint), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "密码更新成功", nil)
}
