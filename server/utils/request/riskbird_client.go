package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// RiskBirdAPIClient 外部 RiskBird 系统 API 客户端
type RiskBirdAPIClient struct {
	BaseURL string
	Client  *http.Client
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// BalanceResponse 余额响应结构
type BalanceResponse struct {
	Code int `json:"code"`
	Data struct {
		TotalBalance float64 `json:"totalBalance"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// PreOrderResponse 预订单响应结构
type PreOrderResponse struct {
	Code int `json:"code"`
	Data struct {
		OrderNo string `json:"orderNo"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// OrderResponse 订单响应结构
type OrderResponse struct {
	Code int `json:"code"`
	Data struct {
		OrderNo string `json:"orderNo"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// NewRiskBirdAPIClient 创建 RiskBird API 客户端
func NewRiskBirdAPIClient(baseURL string) *RiskBirdAPIClient {
	return &RiskBirdAPIClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login 用户登录
func (c *RiskBirdAPIClient) Login(mobile, password string) (map[string]interface{}, error) {
	// 构建 URL 并正确编码参数
	params := url.Values{}
	params.Add("mobile", mobile)
	params.Add("password", password)
	urlStr := fmt.Sprintf("%s/loginByPass?%s", c.BaseURL, params.Encode())

	global.GVA_LOG.Info("调用RiskBird登录接口",
		zap.String("url", urlStr),
		zap.String("mobile", mobile))

	resp, err := c.Client.Post(urlStr, "application/json", nil)
	if err != nil {
		global.GVA_LOG.Error("RiskBird登录请求失败", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		global.GVA_LOG.Error("RiskBird登录HTTP状态错误", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		global.GVA_LOG.Error("RiskBird登录响应解析失败", zap.Error(err))
		return nil, err
	}

	if code, ok := result["code"].(float64); ok && int(code) != 20000 {
		global.GVA_LOG.Error("RiskBird登录失败",
			zap.Float64("code", code),
			zap.String("msg", result["msg"].(string)))
		return nil, fmt.Errorf("登录失败: %s", result["msg"])
	}

	global.GVA_LOG.Info("RiskBird登录成功",
		zap.String("mobile", mobile))

	return result["data"].(map[string]interface{}), nil
}

// GetBalance 获取用户余额
func (c *RiskBirdAPIClient) GetBalance(token string) (float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/recharge/account/balance", c.BaseURL), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Data.TotalBalance, nil
}

// CreatePreOrder 创建预订单
func (c *RiskBirdAPIClient) CreatePreOrder(token string, payload map[string]interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/payment/createPreOrder", c.BaseURL), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result PreOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.OrderNo, nil
}

// CreateOrder 创建订单
func (c *RiskBirdAPIClient) CreateOrder(token string, payload map[string]interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/payment/createOrder", c.BaseURL), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.OrderNo, nil
}

// UpdateOrder 更新订单状态
func (c *RiskBirdAPIClient) UpdateOrder(token string, orderNo string, status string) error {
	payload := map[string]string{
		"orderNo": orderNo,
		"result":  status,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/payment/updateOrder", c.BaseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// GetPointOverview 获取用户积分信息
func (c *RiskBirdAPIClient) GetPointOverview(token string) (int64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/point/overview", c.BaseURL), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", token)

	resp, err := c.Client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("获取积分信息请求失败", zap.Error(err))
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		global.GVA_LOG.Error("获取积分信息HTTP状态错误", zap.Int("status", resp.StatusCode))
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		global.GVA_LOG.Error("积分信息响应解析失败", zap.Error(err))
		return 0, err
	}

	if code, ok := result["code"].(float64); ok && int(code) != 20000 {
		global.GVA_LOG.Error("获取积分信息失败", zap.String("msg", result["msg"].(string)))
		return 0, fmt.Errorf("code %d", int(code))
	}

	data := result["data"].(map[string]interface{})
	availablePoints := int64(data["availablePoints"].(float64))

	return availablePoints, nil
}

// ExpirePoint 使积分失效
func (c *RiskBirdAPIClient) ExpirePoint(token string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/guest/job/expirePoint", c.BaseURL), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := c.Client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("积分失效请求失败", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		global.GVA_LOG.Error("积分失效HTTP状态错误", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// PointAuditDay 积分日审核定时任务
func (c *RiskBirdAPIClient) PointAuditDay(token string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/guest/job/pointAuditDay", c.BaseURL), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := c.Client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("积分日审核定时任务请求失败", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		global.GVA_LOG.Error("积分日审核定时任务HTTP状态错误", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// AdminLogin 管理员登录
func (c *RiskBirdAPIClient) AdminLogin(username, password string) (string, error) {
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	urlStr := fmt.Sprintf("http://mgrtest.riskbird.com/prod-api/account/login?%s", params.Encode())

	global.GVA_LOG.Info("调用RiskBird管理员登录接口",
		zap.String("url", urlStr),
		zap.String("username", username))

	resp, err := c.Client.Post(urlStr, "application/json", nil)
	if err != nil {
		global.GVA_LOG.Error("管理员登录请求失败", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		global.GVA_LOG.Error("管理员登录HTTP状态错误", zap.Int("status", resp.StatusCode))
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		global.GVA_LOG.Error("管理员登录响应解析失败", zap.Error(err))
		return "", err
	}

	if code, ok := result["code"].(float64); ok && int(code) != 20000 {
		global.GVA_LOG.Error("管理员登录失败", zap.String("msg", result["msg"].(string)))
		return "", fmt.Errorf("code %d", int(code))
	}

	data := result["data"].(map[string]interface{})
	token := data["token"].(string)

	return token, nil
}

// AuditPointAcquisition 对积分获取记录进行审核
func (c *RiskBirdAPIClient) AuditPointAcquisition(adminToken string, pointAcquisitionID int64) error {
	payload := map[string]interface{}{
		"auditResult": 1,
		"ids":         []int64{pointAcquisitionID},
		"type":        2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://mgrtest.riskbird.com/prod-api/admin/point/acquisition/audit/operate", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("积分审核请求失败", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		global.GVA_LOG.Error("积分审核HTTP状态错误", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}
