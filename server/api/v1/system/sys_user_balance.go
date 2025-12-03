package system

import (
	"math"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserBalanceApi struct{}

// ModifyUserBalance
// @Tags     UserBalance
// @Summary  修改用户余额
// @Produce   application/json
// @Param    data  body      systemReq.ModifyUserBalance                      true  "手机号, 密码, 充值金额, 赠送金额"
// @Success  200   {object}  response.Response{msg=string}                    "修改用户余额成功"
// @Router   /riskbird/user/modifyUserBalance [post]
func (u *UserBalanceApi) ModifyUserBalance(c *gin.Context) {
	var req systemReq.ModifyUserBalance
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
	if req.RechargeAmount < 0 || req.GiftAmount < 0 {
		response.FailWithMessage("充值金额和赠送金额不能为负数", c)
		return
	}

	// 验证小数点后最多2位
	if !isValidDecimal(req.RechargeAmount, 2) {
		response.FailWithMessage("充值金额最多支持小数点后2位", c)
		return
	}
	if !isValidDecimal(req.GiftAmount, 2) {
		response.FailWithMessage("赠送金额最多支持小数点后2位", c)
		return
	}

	userBalanceService := service.ServiceGroupApp.SystemServiceGroup.UserBalanceService
	err = userBalanceService.ModifyUserBalance(req)
	if err != nil {
		global.GVA_LOG.Error("修改用户余额失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("用户余额修改成功", c)
}

// isValidDecimal 验证数字的小数点位数
func isValidDecimal(value float64, maxDecimalPlaces int) bool {
	// 将数字乘以10^maxDecimalPlaces后，检查是否为整数
	multiplied := value * math.Pow(10, float64(maxDecimalPlaces))
	// 如果乘以后的值与其四舍五入值相等，则说明小数位数不超过maxDecimalPlaces
	return math.Abs(multiplied-math.Round(multiplied)) < 1e-9
}
