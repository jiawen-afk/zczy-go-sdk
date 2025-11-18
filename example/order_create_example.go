package example

import (
	"fmt"
	"log"

	"github.com/jiawen-afk/zczy-go-sdk"
)

func OrderCreateExample() {
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

	// 示例1：使用完整构建器模式创建订单（推荐）
	fmt.Println("=== 示例1：使用完整构建器模式创建订单 ===")
	createOrderWithFullBuilder(client)

	// 示例2：使用部分构建器创建订单
	fmt.Println("\n=== 示例2：使用部分构建器创建订单 ===")
	createOrderWithPartialBuilder(client)

	// 示例3：创建多货物订单
	fmt.Println("\n=== 示例3：创建多货物订单（完整构建器）===")
	createMultiCargoOrder(client)
}

// 示例1：使用完整构建器模式创建订单（推荐）
func createOrderWithFullBuilder(client *zczy.Client) {
	// 使用构建器创建订单信息
	orderInfo := zczy.NewOrderInfoBuilder().
		SetOrderModel("抢单").
		SetFreightType("单价").
		SetSelfComment("TEST20250117001").
		SetContact("赵先生", "13800000000").
		SetVehicle("高栏车", "12").
		SetTimeSchedule("2025-01-20 08:00", "2025-01-20 12:00", "2025-01-25 12:00").
		SetTotalAmount("5000.00").
		SetCargoMoney("50000.00").
		SetPrompt("需要叉车装卸").
		SetSettleBasis("按发货磅单结算").
		SetInterceptPrice("5500").
		SetUrgent(false).
		SetAdvance(false, "").
		SetReceipt(false).
		SetPolicy("1").
		SetOilCard(false, "", "").
		SetAdvisoryPhones("13800000001", "13800000002").
		Build()

	// 使用构建器创建货物信息
	cargo := zczy.NewCargoInfoBuilder().
		SetCargoName("钢材").
		SetCargoVersion("Q235").
		SetCargoCategory("重货").
		SetWeight("30.0").
		SetDimensions("6", "2", "1.5").
		SetPack("捆").
		Build()

	// 使用构建器创建收发货地址信息
	addressInfo := zczy.NewOrderAddressInfoBuilder().
		SetDespatchContact("江苏中储智运物流有限公司", "李先生", "13800000003", "13800000004").
		SetDespatchAddress("江苏省", "南京市", "鼓楼区", "燕江路201号").
		SetDeliverCompany("上海物流有限公司").
		SetDeliverContact("王先生", "13800000005", "13800000006").
		SetDeliverAddress("上海市", "上海市", "浦东新区", "张江高科技园区100号").
		Build()

	// 使用顶层构建器组合所有信息
	req := zczy.NewCreateOrderRequestBuilder().
		SetOrderInfo(*orderInfo).
		AddCargo(*cargo).
		SetOrderAddressInfo(*addressInfo).
		Build()

	resp, err := client.CreateOrder(req)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}

	fmt.Printf("订单创建成功！订单号: %s\n", resp.OrderID)
}

// 示例2：使用部分构建器创建订单
func createOrderWithPartialBuilder(client *zczy.Client) {
	// 订单信息使用构建器
	orderInfo := zczy.NewOrderInfoBuilder().
		SetOrderModel("抢单").
		SetFreightType("单价").
		SetSelfComment("TEST20250117002").
		SetContact("张先生", "13800000010").
		SetVehicle("厢式车", "17").
		SetTimeSchedule("2025-01-22 09:00", "2025-01-22 17:00", "2025-01-28 18:00").
		SetCargoMoney("80000.00").
		SetSettleBasis("按收货磅单结算").
		SetInterceptPrice("6000").
		SetUrgent(true).        // 加急订单
		SetAdvance(true, "30"). // 预付30%
		SetReceipt(false).
		SetPolicy("1").
		Build()

	// 货物信息直接使用结构体
	cargoList := []zczy.CargoInfo{
		{
			CargoName:     "煤炭",
			CargoCategory: "重货",
			Weight:        "25.0",
			Pack:          "散装",
		},
	}

	// 地址信息使用构建器
	addressInfo := zczy.NewOrderAddressInfoBuilder().
		SetDespatchContact("山西煤业有限公司", "赵先生", "13800000011", "").
		SetDespatchAddress("山西省", "太原市", "小店区", "煤炭基地1号").
		SetDeliverCompany("河北钢铁有限公司").
		SetDeliverContact("刘先生", "13800000012", "").
		SetDeliverAddress("河北省", "石家庄市", "长安区", "工业园区88号").
		Build()

	// 组合请求
	req := &zczy.CreateOrderRequest{
		OrderInfo:        *orderInfo,
		CargoList:        cargoList,
		OrderAddressInfo: *addressInfo,
	}

	resp, err := client.CreateOrder(req)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}

	fmt.Printf("订单创建成功！订单号: %s\n", resp.OrderID)
}

// 示例3：创建多货物订单
func createMultiCargoOrder(client *zczy.Client) {
	orderInfo := zczy.NewOrderInfoBuilder().
		SetOrderModel("抢单").
		SetFreightType("包车价").
		SetSelfComment("MULTI20250117001").
		SetContact("陈先生", "13800000020").
		SetVehicle("厢式车", "17").
		SetTimeSchedule("2025-01-22 09:00", "2025-01-22 17:00", "2025-01-28 18:00").
		SetCargoMoney("100000.00").
		SetSettleBasis("按收货磅单结算").
		SetInterceptPrice("8000").
		SetUrgent(true).
		SetAdvance(true, "50").
		SetReceipt(true).
		SetPolicy("2").
		SetOilCard(true, "20", "5"). // 油品20%，汽品5%
		SetOrderMarking("A仓库").
		Build()

	// 使用构建器创建多个货物
	cargo1 := zczy.NewCargoInfoBuilder().
		SetCargoName("电子产品").
		SetCargoVersion("笔记本电脑").
		SetCargoCategory("泡货").
		SetWeight("2.5").
		SetDimensions("1.2", "0.8", "1.0").
		SetPack("纸箱").
		Build()

	cargo2 := zczy.NewCargoInfoBuilder().
		SetCargoName("五金配件").
		SetCargoVersion("标准件").
		SetCargoCategory("重货").
		SetWeight("5.0").
		SetDimensions("0.8", "0.6", "0.5").
		SetPack("木箱").
		Build()

	cargo3 := zczy.NewCargoInfoBuilder().
		SetCargoName("办公用品").
		SetCargoCategory("泡货").
		SetWeight("1.2").
		SetPack("纸箱").
		Build()

	addressInfo := zczy.NewOrderAddressInfoBuilder().
		SetDespatchContact("深圳电子科技有限公司", "周先生", "13800000030", "").
		SetDespatchAddress("广东省", "深圳市", "南山区", "科技园南区A栋").
		SetDeliverCompany("北京贸易有限公司").
		SetDeliverContact("吴先生", "13800000031", "").
		SetDeliverAddress("北京市", "北京市", "朝阳区", "CBD商务中心5号楼").
		Build()

	// 使用顶层构建器，链式添加多个货物
	req := zczy.NewCreateOrderRequestBuilder().
		SetOrderInfo(*orderInfo).
		AddCargo(*cargo1).
		AddCargo(*cargo2).
		AddCargo(*cargo3).
		SetOrderAddressInfo(*addressInfo).
		SetOrderReceiptInfo(&zczy.OrderReceiptInfo{
			ReceiptLabel: "回单标签001",
			ReceiptMoney: "1000",
		}).
		Build()

	resp, err := client.CreateOrder(req)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}

	fmt.Printf("多货物订单创建成功！订单号: %s\n", resp.OrderID)
}

// 示例4：取消订单
func cancelOrderExample(client *zczy.Client, orderID string) {
	fmt.Printf("\n=== 示例4：取消订单 ===\n")
	fmt.Printf("正在取消订单: %s\n", orderID)

	err := client.CancelOrder(orderID)
	if err != nil {
		log.Printf("取消订单失败: %v", err)
		return
	}

	fmt.Printf("订单取消成功！\n")
}

// 示例5：回单确认
func confirmReceiptExample(client *zczy.Client, orderID string) {
	fmt.Printf("\n=== 示例5：回单确认 ===\n")
	fmt.Printf("正在确认订单回单: %s\n", orderID)

	// 使用结算金额方式
	req := &zczy.ConfirmReceiptRequest{
		OrderID:         orderID,
		Tonnage:         "20.0",
		SettleMoney:     "580.00",
		SettleApplyFlag: "0",
		Remark:          "货物已签收",
	}

	err := client.ConfirmReceipt(req)
	if err != nil {
		log.Printf("回单确认失败: %v", err)
		return
	}

	fmt.Printf("回单确认成功！\n")
}

// 示例6：回单确认（使用承运方预估到手价）
func confirmReceiptWithConsignorMoneyExample(client *zczy.Client, orderID string) {
	fmt.Printf("\n=== 示例6：回单确认（使用承运方预估到手价）===\n")
	fmt.Printf("正在确认订单回单: %s\n", orderID)

	// 使用承运方预估到手价方式
	req := &zczy.ConfirmReceiptRequest{
		OrderID:             orderID,
		Tonnage:             "25.5",
		ConsignorNoTaxMoney: "620.00",
		SettleApplyFlag:     "1", // 同时提交结算申请
		Remark:              "",
	}

	err := client.ConfirmReceipt(req)
	if err != nil {
		log.Printf("回单确认失败: %v", err)
		return
	}

	fmt.Printf("回单确认成功！\n")
}
