package zczy

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestVerifyCallbackSign(t *testing.T) {
	// 创建测试客户端
	client := &Client{
		appKey:    "test_app_key",
		appSecret: "test_app_secret",
	}

	// 构建测试数据
	notification := DelistNotification{
		OrderModel:     "0",
		OrderID:        "102019010101018811",
		CarrierName:    "张三",
		DriverUserName: "李四",
		PlateNumber:    "苏A12345",
		Weight:         "12.0",
	}

	dataJSON, _ := json.Marshal(notification)
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	// 生成正确的签名
	params := map[string]string{
		"app_key":   "test_app_key",
		"timestamp": timestampStr,
		"data":      string(dataJSON),
	}
	validSign := client.generateCallbackSign(params)

	tests := []struct {
		name    string
		req     *CallbackRequest
		wantErr bool
	}{
		{
			name: "有效签名",
			req: &CallbackRequest{
				AppKey:    "test_app_key",
				Timestamp: timestampStr,
				Sign:      validSign,
				Data:      string(dataJSON),
			},
			wantErr: false,
		},
		{
			name: "app_key不匹配",
			req: &CallbackRequest{
				AppKey:    "wrong_app_key",
				Timestamp: timestampStr,
				Sign:      validSign,
				Data:      string(dataJSON),
			},
			wantErr: true,
		},
		{
			name: "签名错误",
			req: &CallbackRequest{
				AppKey:    "test_app_key",
				Timestamp: timestampStr,
				Sign:      "INVALID_SIGN",
				Data:      string(dataJSON),
			},
			wantErr: true,
		},
		{
			name: "时间戳过期",
			req: &CallbackRequest{
				AppKey:    "test_app_key",
				Timestamp: "1000000000", // 很久以前的时间戳
				Sign:      validSign,
				Data:      string(dataJSON),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.VerifyCallbackSign(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyCallbackSign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseCallback(t *testing.T) {
	client := &Client{
		appKey:    "test_app_key",
		appSecret: "test_app_secret",
	}

	// 测试摘单通知
	t.Run("解析摘单通知", func(t *testing.T) {
		expectedNotification := DelistNotification{
			OrderModel:        "0",
			OrderID:           "102019010101018811",
			SelfComment:       "TEST001",
			ConsignorUserName: "小王",
			ConsignorMobile:   "13712345678",
			ConsignorState:    "5",
			CargoName:         "钢铁",
			CarrierName:       "张三",
			CarrierMobile:     "13687654321",
			DelistTime:        "2025-01-18 10:00:00",
			Weight:            "12.0",
			PlateNumber:       "苏A12345",
			DriverUserName:    "李四",
			DriverMobile:      "13598765432",
		}

		dataJSON, _ := json.Marshal(expectedNotification)
		timestamp := time.Now().Unix()
		timestampStr := strconv.FormatInt(timestamp, 10)

		params := map[string]string{
			"app_key":   "test_app_key",
			"timestamp": timestampStr,
			"data":      string(dataJSON),
		}
		sign := client.generateCallbackSign(params)

		req := &CallbackRequest{
			AppKey:    "test_app_key",
			Timestamp: timestampStr,
			Sign:      sign,
			Data:      string(dataJSON),
		}

		var notification DelistNotification
		err := client.ParseCallback(req, &notification)
		if err != nil {
			t.Errorf("ParseCallback() error = %v", err)
			return
		}

		if notification.OrderID != expectedNotification.OrderID {
			t.Errorf("OrderID = %v, want %v", notification.OrderID, expectedNotification.OrderID)
		}

		if notification.CarrierName != expectedNotification.CarrierName {
			t.Errorf("CarrierName = %v, want %v", notification.CarrierName, expectedNotification.CarrierName)
		}
	})

	// 测试违约结果通知
	t.Run("解析违约结果通知", func(t *testing.T) {
		expectedNotification := BreachResultNotification{
			OrderID:         "102019010101018811",
			ConsignorState:  "8",
			Operation:       "1",
			ConsignorAmount: "1000.00",
			IsStop:          "1",
			PlatformResults: "1",
		}

		dataJSON, _ := json.Marshal(expectedNotification)
		timestamp := time.Now().Unix()
		timestampStr := strconv.FormatInt(timestamp, 10)

		params := map[string]string{
			"app_key":   "test_app_key",
			"timestamp": timestampStr,
			"data":      string(dataJSON),
		}
		sign := client.generateCallbackSign(params)

		req := &CallbackRequest{
			AppKey:    "test_app_key",
			Timestamp: timestampStr,
			Sign:      sign,
			Data:      string(dataJSON),
		}

		var notification BreachResultNotification
		err := client.ParseCallback(req, &notification)
		if err != nil {
			t.Errorf("ParseCallback() error = %v", err)
			return
		}

		if notification.OrderID != expectedNotification.OrderID {
			t.Errorf("OrderID = %v, want %v", notification.OrderID, expectedNotification.OrderID)
		}

		if notification.Operation != expectedNotification.Operation {
			t.Errorf("Operation = %v, want %v", notification.Operation, expectedNotification.Operation)
		}
	})
}

func TestGenerateCallbackSign(t *testing.T) {
	client := &Client{
		appKey:    "test_key",
		appSecret: "test_secret",
	}

	params := map[string]string{
		"app_key":   "test_key",
		"timestamp": "1737187200",
		"data":      `{"orderId":"123"}`,
	}

	sign := client.generateCallbackSign(params)

	// 签名应该是32位的大写MD5字符串
	if len(sign) != 32 {
		t.Errorf("签名长度错误，期望32，实际%d", len(sign))
	}

	// 验证是否是大写
	for _, c := range sign {
		if c >= 'a' && c <= 'z' {
			t.Errorf("签名包含小写字母: %s", sign)
			break
		}
	}

	// 同样的参数应该产生同样的签名
	sign2 := client.generateCallbackSign(params)
	if sign != sign2 {
		t.Errorf("相同参数产生不同签名: %s != %s", sign, sign2)
	}
}
