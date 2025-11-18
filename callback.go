package zczy

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DelistNotification 摘单通知回调数据
type DelistNotification struct {
	OrderModel        string `json:"orderModel"`        // 0-抢单，1-竞价
	OrderID           string `json:"orderId"`           // 订单号，如果是批量货，就是子单号
	YardID            string `json:"yardId"`            // 如果为批量货，就为母单号，普通货为空
	SelfComment       string `json:"selfComment"`       // 自定义单号
	ConsignorUserName string `json:"consignorUserName"` // 货主名称
	ConsignorMobile   string `json:"consignorMobile"`   // 货主名称手机号
	ConsignorState    string `json:"consignorState"`    // 承运状态：5-摘单，6-确认发货，7-确认收货，8-已终止
	CargoName         string `json:"cargoName"`         // 货物名称
	CarrierName       string `json:"carrierName"`       // 承运方姓名
	CarrierMobile     string `json:"carrierMobile"`     // 承运方手机号
	DelistTime        string `json:"delistTime"`        // 摘牌时间，格式：yyyy-mm-dd hh:mm:ss
	Weight            string `json:"weight"`            // 摘单吨位
	PlateNumber       string `json:"plateNumber"`       // 车牌号
	DriverUserName    string `json:"driverUserName"`    // 司机姓名
	DriverMobile      string `json:"driverMobile"`      // 司机手机号
}

// BreachResultNotification 违约结果通知回调数据
type BreachResultNotification struct {
	OrderID         string `json:"orderId"`         // 订单号
	ConsignorState  string `json:"consignorState"`  // 运单状态
	Operation       string `json:"operation"`       // 操作：1-同意，2-驳回（拒绝）
	ConsignorAmount string `json:"consignorAmount"` // 违约金额
	IsStop          string `json:"isStop"`          // 运单是否终止：1-是，0-否
	PlatformResults string `json:"platformResults"` // 是否最终处理结果（固定值1）
}

// CallbackRequest 回调请求（包含验签参数）
type CallbackRequest struct {
	AppKey    string `json:"app_key"`    // 应用标识
	Timestamp string `json:"timestamp"`  // 时间戳（秒）
	Sign      string `json:"sign"`       // 签名
	Data      string `json:"data"`       // 业务数据（JSON字符串）
}

// VerifyCallbackSign 验证回调签名
// 签名规则：MD5(appSecret + key1value1key2value2... + appSecret)，转大写
// 参数按照ASCII码升序排序，不包括sign字段
func (c *Client) VerifyCallbackSign(req *CallbackRequest) error {
	if req.AppKey != c.appKey {
		return errors.New("app_key不匹配")
	}

	// 验证时间戳（允许5分钟误差）
	timestamp, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("时间戳格式错误: %v", err)
	}

	now := time.Now().Unix()
	if abs(now-timestamp) > 1800 { // 30分钟
		return errors.New("时间戳过期")
	}

	// 构建待签名参数（不包括sign字段）
	params := map[string]string{
		"app_key":   req.AppKey,
		"timestamp": req.Timestamp,
		"data":      req.Data,
	}

	// 计算签名
	expectedSign := c.generateCallbackSign(params)

	if req.Sign != expectedSign {
		return fmt.Errorf("签名验证失败: 期望=%s, 实际=%s", expectedSign, req.Sign)
	}

	return nil
}

// generateCallbackSign 生成回调签名
func (c *Client) generateCallbackSign(params map[string]string) string {
	// 按照key的ASCII顺序排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串：key1value1key2value2...
	var builder strings.Builder
	for _, k := range keys {
		builder.WriteString(k)
		builder.WriteString(params[k])
	}

	// 前后加上appSecret
	signStr := c.appSecret + builder.String() + c.appSecret

	// MD5加密并转大写
	hash := md5.Sum([]byte(signStr))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// ParseCallback 解析回调数据（通用方法）
// 参数 result 必须是指向结构体的指针
// 示例：
//   var delist DelistNotification
//   err := client.ParseCallback(req, &delist)
//
//   var breach BreachResultNotification
//   err := client.ParseCallback(req, &breach)
func (c *Client) ParseCallback(req *CallbackRequest, result any) error {
	// 先验证签名
	if err := c.VerifyCallbackSign(req); err != nil {
		return fmt.Errorf("签名验证失败: %v", err)
	}

	// 解析业务数据到指定的结构体
	if err := json.Unmarshal([]byte(req.Data), result); err != nil {
		return fmt.Errorf("解析业务数据失败: %v", err)
	}

	return nil
}

// abs 返回绝对值
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
