package example

import (
	"fmt"
	"log"

	zczy "github.com/jiawen-afk/zczy-go-sdk"
)

// CoordinateExample 在途轨迹坐标查询示例
func CoordinateExample() {
	// 创建客户端配置
	config := &zczy.Config{
		AppKey:    "your_app_key",
		AppSecret: "your_app_secret",
		PublicKey: "your_public_key",
		Gateway:   "https://seal.zczy100.com/zczy-erp/api", // 联调环境
	}

	// 创建SDK客户端
	client, err := zczy.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 示例1: 获取订单在途轨迹坐标（只传订单号）
	fmt.Println("=== 示例1: 获取订单在途轨迹坐标 ===")
	req := &zczy.OrderCoordinateRequest{
		OrderID: "102019010101018811",
	}

	resp, err := client.GetOrderCoordinate(req)
	if err != nil {
		log.Printf("获取轨迹坐标失败: %v", err)
	} else {
		fmt.Printf("订单号: %s\n", resp.OrderID)
		fmt.Printf("司机: %s\n", resp.DriverName)
		fmt.Printf("车牌号: %s\n", resp.PlateNumber)
		fmt.Printf("司机电话: %s\n", resp.DriverMobile)
		fmt.Printf("轨迹点数量: %d\n\n", len(resp.CoordinateList))

		// 遍历所有轨迹点
		for i, coord := range resp.CoordinateList {
			fmt.Printf("位置 %d:\n", i+1)
			fmt.Printf("  地址: %s\n", coord.Address)
			fmt.Printf("  经度: %s\n", coord.Longitude)
			fmt.Printf("  纬度: %s\n", coord.Latitude)
			fmt.Printf("  定位时间: %s\n", coord.CreatedTime)
			fmt.Printf("  地图类型: %s (1-高德)\n\n", coord.Type)
		}
	}

	// 示例2: 带时间范围查询轨迹坐标
	fmt.Println("=== 示例2: 带时间范围查询轨迹坐标 ===")
	req2 := &zczy.OrderCoordinateRequest{
		OrderID:          "102019010101018811",
		CreatedStartTime: "2021-08-02 12:20",
		CreatedEndTime:   "2021-08-02 13:20",
	}

	resp2, err := client.GetOrderCoordinate(req2)
	if err != nil {
		log.Printf("获取轨迹坐标失败: %v", err)
	} else {
		fmt.Printf("时间范围内的轨迹点数量: %d\n", len(resp2.CoordinateList))

		// 只显示前3个轨迹点
		displayCount := 3
		if len(resp2.CoordinateList) < displayCount {
			displayCount = len(resp2.CoordinateList)
		}

		for i := 0; i < displayCount; i++ {
			coord := resp2.CoordinateList[i]
			fmt.Printf("位置 %d: %s (定位时间: %s)\n",
				i+1, coord.Address, coord.CreatedTime)
		}
	}

	// 示例3: 实时跟踪 - 获取最新位置
	fmt.Println("\n=== 示例3: 获取最新位置 ===")
	resp3, err := client.GetOrderCoordinate(&zczy.OrderCoordinateRequest{
		OrderID: "102019010101018811",
	})
	if err != nil {
		log.Printf("获取轨迹坐标失败: %v", err)
	} else if len(resp3.CoordinateList) > 0 {
		// 假设最后一个点是最新的位置
		latest := resp3.CoordinateList[len(resp3.CoordinateList)-1]
		fmt.Printf("车辆最新位置:\n")
		fmt.Printf("  地址: %s\n", latest.Address)
		fmt.Printf("  经纬度: (%s, %s)\n", latest.Longitude, latest.Latitude)
		fmt.Printf("  更新时间: %s\n", latest.CreatedTime)
	}
}
