package zczy

import (
	"testing"
)

// 测试Coordinate结构体
func TestCoordinate(t *testing.T) {
	coord := &Coordinate{
		Address:     "江苏省南京市鼓楼区宝塔桥街道江山汇金商务楼",
		Longitude:   "118.765659000",
		Latitude:    "32.116436000",
		CreatedTime: "2018-04-18 17:55:55",
		Type:        "1",
	}

	if coord.Address != "江苏省南京市鼓楼区宝塔桥街道江山汇金商务楼" {
		t.Errorf("Address设置失败")
	}

	if coord.Longitude != "118.765659000" {
		t.Errorf("Longitude设置失败")
	}

	if coord.Latitude != "32.116436000" {
		t.Errorf("Latitude设置失败")
	}

	if coord.Type != "1" {
		t.Errorf("Type设置失败")
	}
}

// 测试OrderCoordinateRequest结构体
func TestOrderCoordinateRequest(t *testing.T) {
	req := &OrderCoordinateRequest{
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

// 测试OrderCoordinateResponse结构体
func TestOrderCoordinateResponse(t *testing.T) {
	resp := &OrderCoordinateResponse{
		OrderID:      "10210901012212",
		DriverName:   "赵四",
		PlateNumber:  "苏AT0001",
		DriverMobile: "13800xxxxxx",
		CoordinateList: []Coordinate{
			{
				Address:     "江苏省南京市鼓楼区宝塔桥街道江山汇金商务楼",
				Longitude:   "118.765659000",
				Latitude:    "32.116436000",
				CreatedTime: "2018-04-18 17:55:55",
				Type:        "1",
			},
		},
	}

	if resp.OrderID != "10210901012212" {
		t.Errorf("OrderID设置失败")
	}

	if resp.DriverName != "赵四" {
		t.Errorf("DriverName设置失败")
	}

	if resp.PlateNumber != "苏AT0001" {
		t.Errorf("PlateNumber设置失败")
	}

	if len(resp.CoordinateList) != 1 {
		t.Errorf("CoordinateList长度错误")
	}

	if resp.CoordinateList[0].Address != "江苏省南京市鼓楼区宝塔桥街道江山汇金商务楼" {
		t.Errorf("CoordinateList[0].Address设置失败")
	}
}
