package request

// ModifyUserPoint 修改用户积分请求
type ModifyUserPoint struct {
	Phone       string `json:"phone" binding:"required"`       // 手机号
	Password    string `json:"password" binding:"required"`    // 密码
	PointAmount int64  `json:"pointAmount" binding:"required"` // 积分数量
}
