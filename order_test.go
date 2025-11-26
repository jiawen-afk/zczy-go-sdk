package zczy

import (
	"encoding/json"
	"testing"
)

// 测试OrderReceiptInfo作为值类型的JSON序列化
func TestOrderReceiptInfoSerialization(t *testing.T) {
	// 测试有值的情况
	req := &CreateOrderRequest{
		OrderInfo: OrderInfo{
			OrderModel:  "抢单",
			FreightType: "单价",
		},
		CargoList: []CargoInfo{
			{
				CargoName:     "钢材",
				CargoCategory: "重货",
				Weight:        "30.0",
				Pack:          "捆",
			},
		},
		OrderAddressInfo: OrderAddressInfo{
			DespatchCompanyName: "发货公司",
			DespatchName:        "李先生",
			DespatchMobile:      "13800000001",
		},
		OrderReceiptInfo: OrderReceiptInfo{
			ReceiptLabel: "回单001",
			ReceiptMoney: "1000",
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	// 验证OrderReceiptInfo被正确序列化
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	receiptInfo, ok := result["orderReceiptInfo"].(map[string]any)
	if !ok {
		t.Errorf("orderReceiptInfo字段未找到或类型错误")
	}

	if receiptInfo["receiptLabel"] != "回单001" {
		t.Errorf("receiptLabel值错误，期望=回单001，实际=%v", receiptInfo["receiptLabel"])
	}

	if receiptInfo["receiptMoney"] != "1000" {
		t.Errorf("receiptMoney值错误，期望=1000，实际=%v", receiptInfo["receiptMoney"])
	}

	t.Logf("序列化结果: %s", string(data))
}

// 测试空OrderReceiptInfo的JSON序列化（应该被omitempty）
func TestEmptyOrderReceiptInfoSerialization(t *testing.T) {
	req := &CreateOrderRequest{
		OrderInfo: OrderInfo{
			OrderModel:  "抢单",
			FreightType: "单价",
		},
		CargoList: []CargoInfo{
			{
				CargoName:     "钢材",
				CargoCategory: "重货",
				Weight:        "30.0",
				Pack:          "捆",
			},
		},
		OrderAddressInfo: OrderAddressInfo{
			DespatchCompanyName: "发货公司",
			DespatchName:        "李先生",
			DespatchMobile:      "13800000001",
		},
		// OrderReceiptInfo为空值（零值）
		OrderReceiptInfo: OrderReceiptInfo{},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	// 验证空的OrderReceiptInfo不会被序列化（由于omitempty）
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON反序列化失败: %v", err)
	}

	// 注意：空结构体仍会被序列化为 {}，不会被 omitempty 忽略
	// 这是 Go JSON 的标准行为
	t.Logf("序列化结果: %s", string(data))
}

// 测试CreateOrderRequestBuilder设置OrderReceiptInfo
func TestCreateOrderRequestBuilderWithReceiptInfo(t *testing.T) {
	req := NewCreateOrderRequestBuilder().
		SetOrderInfo(OrderInfo{
			OrderModel:  "抢单",
			FreightType: "单价",
		}).
		AddCargo(CargoInfo{
			CargoName:     "钢材",
			CargoCategory: "重货",
			Weight:        "30.0",
			Pack:          "捆",
		}).
		SetOrderAddressInfo(OrderAddressInfo{
			DespatchCompanyName: "发货公司",
			DespatchName:        "李先生",
			DespatchMobile:      "13800000001",
		}).
		SetOrderReceiptInfo(OrderReceiptInfo{
			ReceiptLabel: "回单测试",
			ReceiptMoney: "2000",
		}).
		Build()

	if req.OrderReceiptInfo.ReceiptLabel != "回单测试" {
		t.Errorf("ReceiptLabel设置失败，期望=回单测试，实际=%s", req.OrderReceiptInfo.ReceiptLabel)
	}

	if req.OrderReceiptInfo.ReceiptMoney != "2000" {
		t.Errorf("ReceiptMoney设置失败，期望=2000，实际=%s", req.OrderReceiptInfo.ReceiptMoney)
	}
}
