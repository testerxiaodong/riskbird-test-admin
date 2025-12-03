package request

// ModifyUserBalance 修改外部系统用户余额请求结构
type ModifyUserBalance struct {
	Phone          string  `json:"phone" binding:"required"`    // 用户手机号
	Password       string  `json:"password" binding:"required"` // 用户密码
	RechargeAmount float64 `json:"rechargeAmount"`              // 充值金额（最多小数点后2位）
	GiftAmount     float64 `json:"giftAmount"`                  // 赠送金额（最多小数点后2位）
}
