package example

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jiawen-afk/zczy-go-sdk"
)

// CallbackExample 回调接口示例
func CallbackExample() {
	// 配置客户端
	config := &zczy.Config{
		AppKey:    "your_app_key_here",
		AppSecret: "your_app_secret_here",
		PublicKey: "your_public_key_here",
		Gateway:   zczy.DefaultGateway,
		Timeout:   30,
	}

	client, err := zczy.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 启动HTTP服务器接收回调
	http.HandleFunc("/callback/delist", func(w http.ResponseWriter, r *http.Request) {
		handleDelistCallback(client, w, r)
	})

	http.HandleFunc("/callback/breach", func(w http.ResponseWriter, r *http.Request) {
		handleBreachResultCallback(client, w, r)
	})

	fmt.Println("回调服务器启动在 :8080")
	fmt.Println("摘单通知: http://localhost:8080/callback/delist")
	fmt.Println("违约结果通知: http://localhost:8080/callback/breach")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleDelistCallback 处理摘单通知回调
func handleDelistCallback(client *zczy.Client, w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("读取请求体失败: %v", err)
		http.Error(w, "读取请求失败", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析回调请求
	var callbackReq zczy.CallbackRequest
	if err := json.Unmarshal(body, &callbackReq); err != nil {
		log.Printf("解析JSON失败: %v", err)
		http.Error(w, "解析请求失败", http.StatusBadRequest)
		return
	}

	// 使用通用方法解析摘单通知
	var notification zczy.DelistNotification
	if err := client.ParseCallback(&callbackReq, &notification); err != nil {
		log.Printf("验证签名或解析数据失败: %v", err)
		http.Error(w, "验证失败", http.StatusUnauthorized)
		return
	}

	// 处理业务逻辑
	fmt.Printf("=== 收到摘单通知 ===\n")
	fmt.Printf("订单号: %s\n", notification.OrderID)
	fmt.Printf("订单类型: %s\n", getOrderModelDesc(notification.OrderModel))
	fmt.Printf("货主名称: %s (%s)\n", notification.ConsignorUserName, notification.ConsignorMobile)
	fmt.Printf("承运方: %s (%s)\n", notification.CarrierName, notification.CarrierMobile)
	fmt.Printf("司机: %s (%s)\n", notification.DriverUserName, notification.DriverMobile)
	fmt.Printf("车牌号: %s\n", notification.PlateNumber)
	fmt.Printf("货物名称: %s\n", notification.CargoName)
	fmt.Printf("摘单吨位: %s\n", notification.Weight)
	fmt.Printf("摘牌时间: %s\n", notification.DelistTime)
	fmt.Printf("承运状态: %s\n", getConsignorStateDesc(notification.ConsignorState))

	// 返回成功响应
	sendSuccessResponse(w)
}

// handleBreachResultCallback 处理违约结果通知回调
func handleBreachResultCallback(client *zczy.Client, w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("读取请求体失败: %v", err)
		http.Error(w, "读取请求失败", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析回调请求
	var callbackReq zczy.CallbackRequest
	if err := json.Unmarshal(body, &callbackReq); err != nil {
		log.Printf("解析JSON失败: %v", err)
		http.Error(w, "解析请求失败", http.StatusBadRequest)
		return
	}

	// 使用通用方法解析违约结果通知
	var notification zczy.BreachResultNotification
	if err := client.ParseCallback(&callbackReq, &notification); err != nil {
		log.Printf("验证签名或解析数据失败: %v", err)
		http.Error(w, "验证失败", http.StatusUnauthorized)
		return
	}

	// 处理业务逻辑
	fmt.Printf("=== 收到违约结果通知 ===\n")
	fmt.Printf("订单号: %s\n", notification.OrderID)
	fmt.Printf("运单状态: %s\n", notification.ConsignorState)
	fmt.Printf("操作结果: %s\n", getOperationDesc(notification.Operation))
	fmt.Printf("违约金额: %s\n", notification.ConsignorAmount)
	fmt.Printf("运单是否终止: %s\n", getIsStopDesc(notification.IsStop))
	fmt.Printf("是否最终处理结果: %s\n", notification.PlatformResults)

	// 返回成功响应
	sendSuccessResponse(w)
}

// sendSuccessResponse 发送成功响应
func sendSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"code":    "0000",
		"message": "success",
	})
}

// getOrderModelDesc 获取订单类型描述
func getOrderModelDesc(model string) string {
	switch model {
	case "0":
		return "抢单"
	case "1":
		return "竞价"
	default:
		return "未知"
	}
}

// getConsignorStateDesc 获取承运状态描述
func getConsignorStateDesc(state string) string {
	switch state {
	case "5":
		return "摘单"
	case "6":
		return "确认发货"
	case "7":
		return "确认收货"
	case "8":
		return "已终止"
	default:
		return "未知"
	}
}

// getOperationDesc 获取操作描述
func getOperationDesc(operation string) string {
	switch operation {
	case "1":
		return "同意"
	case "2":
		return "驳回（拒绝）"
	default:
		return "未知"
	}
}

// getIsStopDesc 获取运单是否终止描述
func getIsStopDesc(isStop string) string {
	switch isStop {
	case "1":
		return "是"
	case "0":
		return "否"
	default:
		return "未知"
	}
}
