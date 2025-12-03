package system

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"go.uber.org/zap"
)

type UserBalanceService struct{}

var UserBalanceServiceApp = new(UserBalanceService)

// ModifyUserBalance 修改外部系统用户余额
func (s *UserBalanceService) ModifyUserBalance(req systemReq.ModifyUserBalance) error {
	// 验证金额
	if req.RechargeAmount < 0 || req.GiftAmount < 0 {
		return errors.New("修改后的金额不能为负数")
	}

	// 连接 RiskBird 数据库
	riskBirdDB, err := request.NewRiskBirdDB(request.RiskBirdDBConfig{
		Host:     "106.75.65.54",
		Port:     7749,
		User:     "zengdong",
		Password: "Zd@2025",
		Database: "riskbirdmdev",
	})
	if err != nil {
		global.GVA_LOG.Error("连接RiskBird数据库失败", zap.Error(err))
		return errors.New("连接数据库失败")
	}
	defer riskBirdDB.Close()

	// 创建 RiskBird API 客户端
	riskBirdClient := request.NewRiskBirdAPIClient("http://test-assets.riskbird.com/test-qbb-api")

	// 1. 用户登录
	loginResp, err := riskBirdClient.Login(req.Phone, req.Password)
	if err != nil {
		global.GVA_LOG.Error("RiskBird用户登录失败",
			zap.String("phone", req.Phone),
			zap.String("error_message", err.Error()))
		return errors.New("请输入正确的手机号和密码")
	}

	token := loginResp["token"].(string)
	global.GVA_LOG.Info("RiskBird用户登录成功", zap.String("phone", req.Phone))
	fmt.Println(token)

	// 2. 获取当前余额
	currentBalance, err := riskBirdClient.GetBalance(token)
	if err != nil {
		global.GVA_LOG.Error("获取用户余额失败", zap.Error(err))
		return errors.New("获取用户余额失败")
	}

	// 3. 如果余额大于0，先花光当前余额
	if currentBalance > 0 {
		// 3.1 修改企业信用报告导出价格为当前余额（使其能被花光）
		err = request.UpdateProductCfg(riskBirdDB, 12, currentBalance)
		if err != nil {
			global.GVA_LOG.Error("修改产品配置失败", zap.Error(err))
			return errors.New("修改产品配置失败")
		}

		// 3.2 创建企业信用报告预订单
		preOrderPayload := map[string]interface{}{
			"productCode":     "paid_report",
			"productNum":      "2",
			"sendEmail":       "",
			"totalAmount":     currentBalance,
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

		// 3.3 创建企业信用报告订单
		reportOrderPayload := map[string]interface{}{
			"balanceAmount":     currentBalance,
			"payAmount":         "0.00",
			"payMethod":         "balance",
			"productNum":        "2",
			"totalAmount":       currentBalance,
			"tradeType":         "JSAPI",
			"unifiedPreOrderNo": preOrderNo,
		}

		reportOrderNo, err := riskBirdClient.CreateOrder(token, reportOrderPayload)
		if err != nil {
			global.GVA_LOG.Error("创建企业信用报告订单失败", zap.Error(err))
			return errors.New("创建企业信用报告订单失败")
		}

		// 3.4 更新报告订单状态为成功
		err = riskBirdClient.UpdateOrder(token, reportOrderNo, "success")
		if err != nil {
			global.GVA_LOG.Error("更新企业信用报告订单失败", zap.Error(err))
			return err
		}

		global.GVA_LOG.Info(fmt.Sprintf("已花光用户当前余额，总金额：%.2f元", currentBalance))

		// 3.5 恢复企业信用报告导出价格为默认值5元
		err = request.UpdateProductCfg(riskBirdDB, 12, 5)
		if err != nil {
			global.GVA_LOG.Error("恢复产品配置失败", zap.Error(err))
			return errors.New("恢复产品配置失败")
		}
	}

	// 4. 充值指定金额
	// 4.1 修改充值套餐金额
	err = request.UpdateRechargeProduct(riskBirdDB, 5, req.RechargeAmount, req.GiftAmount)
	if err != nil {
		global.GVA_LOG.Error("修改充值套餐失败", zap.Error(err))
		return errors.New("修改充值套餐失败")
	}

	// 4.2 创建充值预订单
	rechargePreOrderPayload := map[string]any{
		"productCode":     "",
		"productNum":      1,
		"totalAmount":     req.RechargeAmount,
		"transactionType": "P",
		"selectConditionData": map[string]any{
			"productId": "5",
		},
	}

	rechargePreOrderNo, err := riskBirdClient.CreatePreOrder(token, rechargePreOrderPayload)
	if err != nil {
		global.GVA_LOG.Error("创建充值预订单失败", zap.Error(err))
		return errors.New("创建充值预订单失败")
	}

	// 4.3 创建充值订单
	rechargeOrderPayload := map[string]interface{}{
		"balanceAmount":     0,
		"payAmount":         req.RechargeAmount,
		"payMethod":         "webpay",
		"productNum":        1,
		"productCode":       "",
		"totalAmount":       req.RechargeAmount,
		"tradeType":         "JSAPI",
		"unifiedPreOrderNo": rechargePreOrderNo,
	}

	rechargeOrderNo, err := riskBirdClient.CreateOrder(token, rechargeOrderPayload)
	if err != nil {
		global.GVA_LOG.Error("创建充值订单失败", zap.Error(err))
		return errors.New("创建充值订单失败")
	}

	// 4.4 更新充值订单状态为成功
	err = riskBirdClient.UpdateOrder(token, rechargeOrderNo, "success")
	if err != nil {
		global.GVA_LOG.Error("更新充值订单失败", zap.Error(err))
		return err
	}

	global.GVA_LOG.Info(fmt.Sprintf("已成功为用户充值，充值金额：%.2f元，赠送金额：%.2f元", req.RechargeAmount, req.GiftAmount))

	return nil
}
