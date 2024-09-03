package controller

import (
	"encoding/json"
	"fish_gateway/dao"
	"fish_gateway/dto"
	"fish_gateway/middleware"
	"fish_gateway/public"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func RegisterAdmin(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admin_info", admin.AdminInfo)
	group.POST("/change_pwd", admin.ChangePwd)
}

// AdminInfo godoc
// @Summary 获取管理员信息
// @Description 获取管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfOutput} "success"
// @Router /admin/admin_info [get]
func (a *AdminController) AdminInfo(c *gin.Context) {
	// 1、读取sessionkey 对应的json 转换为结构体
	// 2、取出数据然后封装输出结构体
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	out := &dto.AdminInfOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "",
		Roles:        []string{},
	}
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (a *AdminController) ChangePwd(c *gin.Context) {
	param := &dto.ChangePwdInput{}
	if err := param.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 1. session 读取用户信息
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 2. 使用seesion Name 查询数据库信息
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	// 3. 给密码加盐
	newPassword := public.GenSaltPassword(admin.Salt, param.Password)
	admin.Password = newPassword

	// 4. 保存修改信息
	if err = admin.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}
