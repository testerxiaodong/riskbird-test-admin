import service from '@/utils/request'

// 修改用户积分
export const modifyUserPoint = (data) => {
  return service({
    url: '/riskbird/user/modifyUserPoint',
    method: 'post',
    data
  })
}
