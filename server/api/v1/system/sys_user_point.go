package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserPointApi struct{}

// ModifyUserPoint
// @Tags     UserPoint
// @Summary  修改用户积分
// @Produce   application/json
// @Param    data  body      systemReq.ModifyUserPoint                      true  "手机号, 密码, 修改积分"
// @Success  200   {object}  response.Response{msg=string}                  "修改用户积分成功"
// @Router   /riskbird/user/modifyUserPoint [post]
func (u *UserPointApi) ModifyUserPoint(c *gin.Context) {
	var req systemReq.ModifyUserPoint
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证必填字段
	if req.Phone == "" {
		response.FailWithMessage("用户手机号不能为空", c)
		return
	}
	if req.Password == "" {
		response.FailWithMessage("用户密码不能为空", c)
		return
	}

	// 验证积分
	if req.PointAmount < 0 {
		response.FailWithMessage("积分不能为负数", c)
		return
	}
	if req.PointAmount%5 != 0 {
		response.FailWithMessage("积分必须是5的倍数", c)
		return
	}

	userPointService := service.ServiceGroupApp.SystemServiceGroup.UserPointService
	err = userPointService.ModifyUserPoint(req)
	if err != nil {
		global.GVA_LOG.Error("修改用户积分失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("修改用户积分成功", c)
}
