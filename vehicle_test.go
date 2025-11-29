package zczy

import (
	"strings"
	"testing"
)

// 测试VehicleTrackRequest结构体
func TestVehicleTrackRequest(t *testing.T) {
	req := &VehicleTrackRequest{
		OrderID:          "102019010101018811",
		CreatedStartTime: "2021-08-02 12:20",
		CreatedEndTime:   "2021-08-02 13:20",
	}

	if req.OrderID != "102019010101018811" {
		t.Errorf("OrderID设置失败")
	}

	if req.CreatedStartTime != "2021-08-02 12:20" {
		t.Errorf("CreatedStartTime设置失败")
	}

	if req.CreatedEndTime != "2021-08-02 13:20" {
		t.Errorf("CreatedEndTime设置失败")
	}
}

// 测试VehicleTrackResponse结构体
func TestVehicleTrackResponse(t *testing.T) {
	resp := &VehicleTrackResponse{
		URL: "https://track.example.com/order/123",
	}

	if resp.URL != "https://track.example.com/order/123" {
		t.Errorf("URL设置失败")
	}
}

// 测试GetVehicleTrack方法生成URL
func TestGetVehicleTrackURL(t *testing.T) {
	config := &Config{
		AppKey:    "test_key",
		AppSecret: "test_secret",
		PublicKey: "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALT8QammE81aGfzzmFj0LjHKAOWiyRLESX4fwomlvWr3nVvx4rSzKGz176M/c9UsLQFqJkA0KIk0YxDgS1QG5K8CAwEAAQ==",
		Gateway:   "https://seal.zczy100.com/zczy-erp/api",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() 失败: %v", err)
	}

	req := &VehicleTrackRequest{
		OrderID: "102019010101018811",
	}

	resp, err := client.GetVehicleTrack(req)
	if err != nil {
		t.Fatalf("GetVehicleTrack() 失败: %v", err)
	}

	if resp.URL == "" {
		t.Errorf("GetVehicleTrack() 返回空URL")
	}

	// 验证URL包含必要的参数
	if !strings.Contains(resp.URL, "https://seal.zczy100.com/zczy-erp/html") {
		t.Errorf("URL应该包含正确的基础路径，实际URL: %s", resp.URL)
	}

	if !strings.Contains(resp.URL, "method="+MethodVehicleTrack) {
		t.Errorf("URL应该包含method参数，实际URL: %s", resp.URL)
	}

	if !strings.Contains(resp.URL, "appKey=test_key") {
		t.Errorf("URL应该包含appKey参数，实际URL: %s", resp.URL)
	}

	if !strings.Contains(resp.URL, "sign=") {
		t.Errorf("URL应该包含sign参数，实际URL: %s", resp.URL)
	}

	t.Logf("生成的轨迹URL: %s", resp.URL)
}
