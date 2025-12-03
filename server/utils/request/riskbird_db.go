package request

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// RiskBirdDBConfig RiskBird数据库配置
type RiskBirdDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewRiskBirdDB 创建RiskBird数据库连接
func NewRiskBirdDB(config RiskBirdDBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	return sql.Open("mysql", dsn)
}

// UpdateProductCfg 修改产品配置价格
func UpdateProductCfg(db *sql.DB, id int, value float64) error {
	sql := "UPDATE p_product_cfg SET cfg_value = ? WHERE id = ?"
	_, err := db.Exec(sql, value, id)
	return err
}

// UpdateRechargeProduct 修改充值套餐
func UpdateRechargeProduct(db *sql.DB, id int, amount, giftAmount float64) error {
	sql := "UPDATE p_recharge_product SET amount = ?, gift_amount = ? WHERE id = ?"
	_, err := db.Exec(sql, amount, giftAmount, id)
	return err
}

// UpdatePointExpireTime 修改积分失效时间
func UpdatePointExpireTime(db *sql.DB, userID int64, expireTime interface{}) error {
	sql := "UPDATE point_acquisition SET expire_time = ? WHERE user_id = ? AND left_points > 0"
	_, err := db.Exec(sql, expireTime, userID)
	return err
}

// GetLatestPointAcquisitionID 获取最新的积分获取记录ID
func GetLatestPointAcquisitionID(db *sql.DB, userID int64) (int64, error) {
	var id int64
	sql := "SELECT id FROM point_acquisition WHERE user_id = ? ORDER BY create_time DESC LIMIT 1"
	err := db.QueryRow(sql, userID).Scan(&id)
	return id, err
}

// UpdatePointAcquisitionTime 修改积分获取时间
func UpdatePointAcquisitionTime(db *sql.DB, pointAcquisitionID int64, pointTime interface{}) error {
	sql := "UPDATE point_acquisition SET point_time = ? WHERE id = ?"
	_, err := db.Exec(sql, pointTime, pointAcquisitionID)
	return err
}
