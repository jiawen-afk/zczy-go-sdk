package main

import (
	"fmt"
	"log"

	"github.com/Jiawen-AFK/zczy-go-sdk"
)

func main() {
	// 配置客户端参数
	config := &zczy.Config{
		AppKey:    "your_app_key_here",
		AppSecret: "your_app_secret_here",
		PublicKey: "your_public_key_here",
		Gateway:   zczy.DefaultGateway, // 使用联调环境
		Timeout:   30,
	}

	// 创建客户端实例
	client, err := zczy.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 示例1：调用API（带参数）
	fmt.Println("=== 示例1：调用带参数的API ===")
	params := map[string]any{
		"orderId":   "TEST123456",
		"orderNo":   "ORD202301010001",
		"status":    1,
		"pageNum":   1,
		"pageSize":  10,
		"startTime": "2023-01-01 00:00:00",
		"endTime":   "2023-12-31 23:59:59",
	}

	resp, err := client.Execute("order.query", params)
	if err != nil {
		log.Printf("API调用失败: %v", err)
	} else {
		if resp.IsSuccess() {
			fmt.Printf("调用成功！\n")
			fmt.Printf("返回码: %s\n", resp.Code)
			fmt.Printf("返回消息: %s\n", resp.Message)
			fmt.Printf("返回数据: %+v\n", resp.Data)
		} else {
			fmt.Printf("API返回错误: [%s] %s\n", resp.Code, resp.Message)
		}
	}

	// 示例2：调用API（无参数）
	fmt.Println("\n=== 示例2：调用无参数的API ===")
	resp2, err := client.Execute("system.heartbeat", nil)
	if err != nil {
		log.Printf("API调用失败: %v", err)
	} else {
		if resp2.IsSuccess() {
			fmt.Printf("心跳检测成功！\n")
		} else {
			fmt.Printf("心跳检测失败: [%s] %s\n", resp2.Code, resp2.Message)
		}
	}

	// 示例3：解析响应数据到结构体
	fmt.Println("\n=== 示例3：解析响应数据到结构体 ===")

	// 定义订单数据结构
	type OrderInfo struct {
		OrderID   string `json:"orderId"`
		OrderNo   string `json:"orderNo"`
		Status    int    `json:"status"`
		Amount    string `json:"amount"`
		CreatedAt string `json:"createdAt"`
	}

	// 定义订单列表响应结构
	type OrderListResponse struct {
		Total int          `json:"total"`
		List  []*OrderInfo `json:"list"`
	}

	resp3, err := client.Execute("order.list", map[string]any{
		"pageNum":  1,
		"pageSize": 20,
	})

	if err != nil {
		log.Printf("API调用失败: %v", err)
	} else {
		var orderList OrderListResponse
		if err := resp3.GetData(&orderList); err != nil {
			log.Printf("解析数据失败: %v", err)
		} else {
			fmt.Printf("订单总数: %d\n", orderList.Total)
			fmt.Printf("订单列表: %+v\n", orderList.List)
		}
	}

	// 示例4：切换到正式环境
	fmt.Println("\n=== 示例4：切换到正式环境 ===")
	// 注意：正式环境地址需要从中储智运获取
	productionGateway := "https://production.zczy56.com/api"
	client.SetGateway(productionGateway)
	fmt.Printf("已切换到正式环境: %s\n", productionGateway)

	// 示例5：错误处理
	fmt.Println("\n=== 示例5：完整的错误处理示例 ===")
	resp5, err := client.Execute("test.method", map[string]any{
		"testParam": "value",
	})

	if err != nil {
		// 网络错误、超时等
		log.Printf("请求失败: %v", err)
		return
	}

	// 检查业务错误
	if !resp5.IsSuccess() {
		fmt.Printf("业务错误: 错误码=%s, 错误信息=%s\n", resp5.Code, resp5.Message)

		// 根据错误码进行不同处理
		switch resp5.Code {
		case "1001":
			fmt.Println("处理逻辑: 参数错误，请检查参数")
		case "1002":
			fmt.Println("处理逻辑: 签名验证失败，请检查appKey和appSecret")
		case "1003":
			fmt.Println("处理逻辑: 时间戳过期，请检查系统时间")
		default:
			fmt.Printf("处理逻辑: 未知错误码，请查阅API文档\n")
		}
		return
	}

	fmt.Printf("请求成功: %+v\n", resp5.Data)
}
