package system

import (
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"go.uber.org/zap"
)

type UserPointService struct{}

var UserPointServiceApp = new(UserPointService)

// ModifyUserPoint 修改外部系统用户积分
func (s *UserPointService) ModifyUserPoint(req systemReq.ModifyUserPoint) error {
	// 验证积分
	if req.PointAmount < 0 {
		return errors.New("修改后的积分不能为负数")
	}
	if req.PointAmount%5 != 0 {
		return errors.New("修改后的积分必须是5的倍数")
	}

	// 连接 RiskBird 数据库
	riskBirdDB, err := request.NewRiskBirdDB(request.RiskBirdDBConfig{
		Host:     global.GVA_CONFIG.RiskBird.DB.Host,
		Port:     global.GVA_CONFIG.RiskBird.DB.Port,
		User:     global.GVA_CONFIG.RiskBird.DB.User,
		Password: global.GVA_CONFIG.RiskBird.DB.Password,
		Database: global.GVA_CONFIG.RiskBird.DB.Database,
	})
	if err != nil {
		global.GVA_LOG.Error("连接RiskBird数据库失败", zap.Error(err))
		return errors.New("连接数据库失败")
	}
	defer riskBirdDB.Close()

	// 创建 RiskBird API 客户端
	riskBirdClient := request.NewRiskBirdAPIClient(global.GVA_CONFIG.RiskBird.API.BaseUrl)

	// 1. 用户登录
	loginResp, err := riskBirdClient.Login(req.Phone, req.Password)
	if err != nil {
		global.GVA_LOG.Error("RiskBird用户登录失败",
			zap.String("phone", req.Phone),
			zap.String("error_message", err.Error()))
		return errors.New("请输入正确的手机号和密码")
	}

	token := loginResp["token"].(string)
	userID := int64(loginResp["user"].(map[string]interface{})["id"].(float64))

	global.GVA_LOG.Info("RiskBird用户登录成功", zap.String("phone", req.Phone))

	// 2. 获取用户当前可用积分
	availablePoints, err := riskBirdClient.GetPointOverview(token)
	if err != nil {
		global.GVA_LOG.Error("获取用户积分信息失败", zap.Error(err))
		return errors.New("获取用户积分信息失败")
	}

	// 3. 如果用户有可用积分，先使其失效
	if availablePoints > 0 {
		// 设置积分失效时间为昨天
		expireTime := time.Now().AddDate(0, 0, -1)

		err = request.UpdatePointExpireTime(riskBirdDB, userID, expireTime)
		if err != nil {
			global.GVA_LOG.Error("修改积分失效时间失败", zap.Error(err))
			return errors.New("修改积分失效时间失败")
		}

		// 调用积分失效定时任务接口使积分失效
		err = riskBirdClient.ExpirePoint(token)
		if err != nil {
			global.GVA_LOG.Error("调用积分失效定时任务接口失败", zap.Error(err))
			return errors.New("调用积分失效定时任务接口失败")
		}

		global.GVA_LOG.Info(fmt.Sprintf("已使用户可用积分失效，失效积分数：%d分", availablePoints))
	}

	// 4. 如果修改后的积分为0，直接返回
	if req.PointAmount == 0 {
		global.GVA_LOG.Info(fmt.Sprintf("用户%s的积分已修改为0分", req.Phone))
		return nil
	}

	// 5. 计算支付金额（每5积分对应1元）
	payAmount := float64(req.PointAmount) / 5

	// 设置积分获取时间为昨天
	pointTime := time.Now().AddDate(0, 0, -1)

	// 5.1 修改企业信用报告导出价格为支付金额
	err = request.UpdateProductCfg(riskBirdDB, 12, payAmount)
	if err != nil {
		global.GVA_LOG.Error("修改产品配置失败", zap.Error(err))
		return errors.New("修改产品配置失败")
	}

	// 5.2 创建企业信用报告预订单
	preOrderPayload := map[string]interface{}{
		"productCode":     "paid_report",
		"productNum":      "2",
		"sendEmail":       "",
		"totalAmount":     payAmount,
		"transactionType": "C",
		"tradeType":       "JSAPI",
		"selectConditionData": map[string]interface{}{
			"entName":     "乐视网信息技术（北京）股份有限公司",
			"entid":       "7jShe5V5mqx",
			"fileType":    "pdf,word",
			"groupIdList": "9,2,5,6,7,8,",
		},
	}

	preOrderNo, err := riskBirdClient.CreatePreOrder(token, preOrderPayload)
	if err != nil {
		global.GVA_LOG.Error("创建企业信用报告预订单失败", zap.Error(err))
		return errors.New("创建企业信用报告预订单失败")
	}

	// 5.3 创建企业信用报告订单
	reportOrderPayload := map[string]interface{}{
		"balanceAmount":     0,
		"payAmount":         payAmount,
		"payMethod":         "webpay",
		"productNum":        "2",
		"totalAmount":       payAmount,
		"tradeType":         "JSAPI",
		"unifiedPreOrderNo": preOrderNo,
	}

	reportOrderNo, err := riskBirdClient.CreateOrder(token, reportOrderPayload)
	if err != nil {
		global.GVA_LOG.Error("创建企业信用报告订单失败", zap.Error(err))
		return errors.New("创建企业信用报告订单失败")
	}

	// 5.4 更新报告订单状态为成功
	err = riskBirdClient.UpdateOrder(token, reportOrderNo, "success")
	if err != nil {
		global.GVA_LOG.Error("更新企业信用报告订单失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info(fmt.Sprintf("已完成企业信用报告导出支付，支付金额：%.2f元，待入账积分：%d分", payAmount, req.PointAmount))

	// 5.5 恢复企业信用报告导出价格为默认值5元
	err = request.UpdateProductCfg(riskBirdDB, 12, 5)
	if err != nil {
		global.GVA_LOG.Error("恢复产品配置失败", zap.Error(err))
		return errors.New("恢复产品配置失败")
	}

	// 6. 等待5秒，确保积分获取记录已创建
	time.Sleep(5 * time.Second)

	// 7. 查询最新的积分获取记录ID并修改其发生时间
	pointAcquisitionID, err := request.GetLatestPointAcquisitionID(riskBirdDB, userID)
	if err != nil {
		global.GVA_LOG.Error("查询积分获取记录失败", zap.Error(err))
		return errors.New("查询积分获取记录失败")
	}

	err = request.UpdatePointAcquisitionTime(riskBirdDB, pointAcquisitionID, pointTime)
	if err != nil {
		global.GVA_LOG.Error("修改积分获取时间失败", zap.Error(err))
		return errors.New("修改积分获取时间失败")
	}

	// 8. 调用积分审核日度定时任务
	err = riskBirdClient.PointAuditDay(token)
	if err != nil {
		global.GVA_LOG.Error("调用积分日审核定时任务接口失败", zap.Error(err))
		return errors.New("调用积分日审核定时任务接口失败")
	}

	// 9. 登录后台管理系统进行积分审核
	adminToken, err := riskBirdClient.AdminLogin("zengdong", "zengdong@123")
	if err != nil {
		global.GVA_LOG.Error("管理员登录失败", zap.Error(err))
		return errors.New("管理员登录失败")
	}

	// 10. 对当前用户进行积分审核
	err = riskBirdClient.AuditPointAcquisition(adminToken, pointAcquisitionID)
	if err != nil {
		global.GVA_LOG.Error("积分审核失败", zap.Error(err))
		return errors.New("积分审核失败")
	}

	global.GVA_LOG.Info(fmt.Sprintf("用户%s的积分已修改为%d分，移动端用户请重新登录后查看最新积分", req.Phone, req.PointAmount))

	return nil
}
