import service from '@/utils/request'

// @Tags UserBalance
// @Summary 修改用户余额
// @Produce application/json
// @Param data body {phone:"string",password:"string",rechargeAmount:"number",giftAmount:"number"}
// @Router /user/modifyUserBalance [post]
export const modifyUserBalance = (data) => {
  return service({
    url: '/riskbird/user/modifyUserBalance',
    method: 'post',
    data: data
  })
}
