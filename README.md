# 中储智运开放平台 Golang SDK

中储智运开放平台的官方 Golang SDK，提供便捷的 API 调用接口。

## 功能特性

- 自动签名生成（MD5）
- appSecret RSA 加密
- 支持联调环境和正式环境切换
- 完整的错误处理
- 类型安全的响应解析
- 业务API封装（订单创建等）
- 构建器模式支持，简化复杂参数构建

## 安装

```bash
go get github.com/jiawen-afk/zczy-go-sdk
```

## 快速开始

### 1. 初始化客户端

```go
package main

import (
    "fmt"
    "log"
    "github.com/jiawen-afk/zczy-go-sdk"
)

func main() {
    // 配置客户端
    config := &zczy.Config{
        AppKey:    "your_app_key",
        AppSecret: "your_app_secret",
        PublicKey: "your_public_key_content",
        Gateway:   zczy.DefaultGateway, // 可选，默认为联调环境
        Timeout:   30,                   // 可选，默认30秒
    }

    // 创建客户端
    client, err := zczy.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // 调用API
    params := map[string]any{
        "orderId": "123456",
        "status":  1,
    }

    resp, err := client.Execute("your.api.method", params)
    if err != nil {
        log.Fatal(err)
    }

    // 检查响应
    if resp.IsSuccess() {
        fmt.Printf("Success: %+v\n", resp.Data)
    } else {
        fmt.Printf("Error: %s - %s\n", resp.Code, resp.Message)
    }
}
```

### 2. 处理响应数据

```go
// 定义响应数据结构
type OrderInfo struct {
    OrderID   string `json:"orderId"`
    OrderNo   string `json:"orderNo"`
    Status    int    `json:"status"`
    CreatedAt string `json:"createdAt"`
}

// 调用API
resp, err := client.Execute("order.get", params)
if err != nil {
    log.Fatal(err)
}

// 解析数据到结构体
var orderInfo OrderInfo
if err := resp.GetData(&orderInfo); err != nil {
    log.Fatal(err)
}

fmt.Printf("Order: %+v\n", orderInfo)
```

### 3. 创建订单（使用业务API）

```go
// 使用构建器模式创建订单（推荐）
orderInfo := zczy.NewOrderInfoBuilder().
    SetOrderModel("抢单").
    SetFreightType("单价").
    SetSelfComment("TEST001").
    SetContact("赵先生", "13800000000").
    SetVehicle("高栏车", "12").
    SetTimeSchedule("2025-01-20 08:00", "2025-01-20 12:00", "2025-01-25 12:00").
    SetCargoMoney("50000.00").
    SetSettleBasis("按发货磅单结算").
    SetInterceptPrice("5500").
    SetUrgent(false).
    SetAdvance(true, "30").
    SetReceipt(false).
    SetPolicy("1").
    Build()

req := &zczy.CreateOrderRequest{
    OrderInfo: *orderInfo,
    CargoList: []zczy.CargoInfo{
        {
            CargoName:     "钢材",
            CargoCategory: "重货",
            Weight:        "30.0",
            Pack:          "捆",
        },
    },
    OrderAddressInfo: zczy.OrderAddressInfo{
        DespatchCompanyName: "发货公司",
        DespatchName:        "李先生",
        DespatchMobile:      "13800000001",
        DespatchPro:         "江苏省",
        DespatchCity:        "南京市",
        DespatchDis:         "鼓楼区",
        DespatchPlace:       "燕江路201号",
        DeliverCompanyName:  "收货公司",
        DeliverName:         "王先生",
        DeliverMobile:       "13800000002",
        DeliverPro:          "上海市",
        DeliverCity:         "上海市",
        DeliverDis:          "浦东新区",
        DeliverPlace:        "张江高科技园区100号",
    },
}

resp, err := client.CreateOrder(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("订单创建成功！订单号: %s\n", resp.OrderID)
```

### 4. 切换到正式环境

```go
// 方式1：创建客户端时指定
config := &zczy.Config{
    AppKey:    "your_app_key",
    AppSecret: "your_app_secret",
    PublicKey: "your_public_key",
    Gateway:   "https://production.zczy56.com/api", // 正式环境地址
}
client, _ := zczy.NewClient(config)

// 方式2：动态切换
client.SetGateway("https://production.zczy56.com/api")
```

### 5. 多货主场景配置（可选）

如果您的系统需要支持多个货主，可以通过 `ConsignorId` 参数来指定货主ID：

```go
// 方式1：创建客户端时指定货主ID
config := &zczy.Config{
    AppKey:      "your_app_key",
    AppSecret:   "your_app_secret",
    PublicKey:   "your_public_key",
    ConsignorId: "your_consignor_id", // 货主ID
}
client, _ := zczy.NewClient(config)

// 方式2：动态切换货主ID
client.SetConsignorId("another_consignor_id")

// ConsignorId 会自动添加到所有 API 请求参数中
resp, err := client.CreateOrder(orderReq)
```

## 业务API

### 订单管理

#### 创建普通货订单

方法：`CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error)`

支持单货、多货创建，提供构建器模式简化订单信息构建。

**参数说明：**

- **OrderInfo**: 订单基本信息
  - `OrderModel`: 订单类型（抢单/竞价）
  - `FreightType`: 费用类型（包车价/单价）
  - `VehicleType`: 车型要求
  - `VehicleLength`: 车长要求（米）
  - `SettleBasis`: 结算依据
  - 其他字段详见 [order.go](order.go)

- **CargoList**: 货物信息列表（支持多货物）
  - `CargoName`: 货物名称
  - `CargoCategory`: 货物类别（重货/泡货）
  - `Weight`: 货物重量
  - `Pack`: 包装类型

- **OrderAddressInfo**: 收发货地址信息
  - 发货地址：省、市、区、详细地址
  - 收货地址：省、市、区、详细地址
  - 联系人及电话

- **OrderReceiptInfo**: 押回单信息（可选）

**示例代码：**

参见 [example/order_create_example.go](example/order_create_example.go)

#### 取消订单

方法：`CancelOrder(orderID string) error`

取消指定的订单。

**参数说明：**

- `orderID`: 订单号（必填）

**示例：**

```go
err := client.CancelOrder("102019010101018811")
if err != nil {
    log.Printf("取消订单失败: %v", err)
}
fmt.Println("订单取消成功！")
```

#### 回单确认

方法：`ConfirmReceipt(req *ConfirmReceiptRequest) error`

确认回单并传递收货吨位、结算金额等信息，可选择是否同时提交结算申请。

**参数说明：**

- `OrderID`: 订单号（必填）
- `Tonnage`: 收货吨位（必填）
- `SettleMoney`: 结算金额（与 ConsignorNoTaxMoney 二选一）
- `ConsignorNoTaxMoney`: 承运方预估到手价（与 SettleMoney 二选一）
- `SettleApplyFlag`: 是否提交结算申请（必填），0-否，1-是
- `Remark`: 备注（可选）

**示例：**

```go
// 方式1：使用结算金额
req := &zczy.ConfirmReceiptRequest{
    OrderID:         "102019010101018811",
    Tonnage:         "20.0",
    SettleMoney:     "580.00",
    SettleApplyFlag: "0",
    Remark:          "货物已签收",
}
err := client.ConfirmReceipt(req)
if err != nil {
    log.Printf("回单确认失败: %v", err)
}

// 方式2：使用承运方预估到手价
req := &zczy.ConfirmReceiptRequest{
    OrderID:             "102019010101018811",
    Tonnage:             "25.5",
    ConsignorNoTaxMoney: "620.00",
    SettleApplyFlag:     "1", // 同时提交结算申请
}
err := client.ConfirmReceipt(req)
if err != nil {
    log.Printf("回单确认失败: %v", err)
}
```

## 回调接口

SDK 支持接收中储智运平台的回调通知，并提供完整的签名验证功能。

### 通用回调解析方法

SDK 提供了通用的 `ParseCallback` 方法，可以解析任何类型的回调通知：

```go
// 解析摘单通知
var delistNotification zczy.DelistNotification
err := client.ParseCallback(&callbackReq, &delistNotification)

// 解析违约结果通知
var breachNotification zczy.BreachResultNotification
err := client.ParseCallback(&callbackReq, &breachNotification)
```

### 摘单通知回调

当订单被承运方摘单后，平台会主动推送摘单通知到您配置的回调地址。

**回调数据字段：**

- `orderModel`: 订单类型（0-抢单，1-竞价）
- `orderId`: 订单号（批量货为子单号）
- `yardId`: 母单号（批量货时有值，普通货为空）
- `selfComment`: 自定义单号
- `consignorUserName`: 货主名称
- `consignorMobile`: 货主手机号
- `consignorState`: 承运状态（5-摘单，6-确认发货，7-确认收货，8-已终止）
- `cargoName`: 货物名称
- `carrierName`: 承运方姓名
- `carrierMobile`: 承运方手机号
- `delistTime`: 摘牌时间（格式：yyyy-mm-dd hh:mm:ss）
- `weight`: 摘单吨位
- `plateNumber`: 车牌号
- `driverUserName`: 司机姓名
- `driverMobile`: 司机手机号

**使用示例：**

```go
import (
    "encoding/json"
    "net/http"
    "github.com/jiawen-afk/zczy-go-sdk"
)

func handleDelistCallback(w http.ResponseWriter, r *http.Request) {
    // 读取请求体
    var callbackReq zczy.CallbackRequest
    json.NewDecoder(r.Body).Decode(&callbackReq)

    // 使用通用方法解析
    var notification zczy.DelistNotification
    err := client.ParseCallback(&callbackReq, &notification)
    if err != nil {
        http.Error(w, "验证失败", http.StatusUnauthorized)
        return
    }

    // 处理业务逻辑
    fmt.Printf("收到摘单通知，订单号: %s\n", notification.OrderID)
    fmt.Printf("承运方: %s，司机: %s\n", notification.CarrierName, notification.DriverUserName)

    // 返回成功响应
    json.NewEncoder(w).Encode(map[string]any{
        "code": "0000",
        "message": "success",
    })
}
```

### 违约结果通知回调

当订单发生违约处理后，平台会推送违约结果通知。

**回调数据字段：**

- `orderId`: 订单号
- `consignorState`: 运单状态
- `operation`: 操作结果（1-同意，2-驳回/拒绝）
- `consignorAmount`: 违约金额
- `isStop`: 运单是否终止（1-是，0-否）
- `platformResults`: 是否最终处理结果（固定值1）

**使用示例：**

```go
func handleBreachCallback(w http.ResponseWriter, r *http.Request) {
    var callbackReq zczy.CallbackRequest
    json.NewDecoder(r.Body).Decode(&callbackReq)

    // 使用通用方法解析
    var notification zczy.BreachResultNotification
    err := client.ParseCallback(&callbackReq, &notification)
    if err != nil {
        http.Error(w, "验证失败", http.StatusUnauthorized)
        return
    }

    // 处理业务逻辑
    fmt.Printf("订单号: %s\n", notification.OrderID)
    fmt.Printf("操作结果: %s\n", notification.Operation)
    fmt.Printf("违约金额: %s\n", notification.ConsignorAmount)
    fmt.Printf("运单是否终止: %s\n", notification.IsStop)

    // 返回成功响应
    json.NewEncoder(w).Encode(map[string]any{
        "code": "0000",
        "message": "success",
    })
}
```

**完整示例代码：**

参见 [example/callback_example.go](example/callback_example.go)

**对方加密规范：**

中储智运平台在回调时需要遵循的签名加密规范，详见 [CALLBACK_SIGNATURE.md](CALLBACK_SIGNATURE.md)

该文档包含：
- 完整的签名算法说明
- Java、Python、PHP 等多语言示例代码
- 签名计算完整示例
- 常见问题说明

## API 参数说明

### Config 配置参数

| 参数        | 类型   | 必填 | 说明                              |
| ----------- | ------ | ---- | --------------------------------- |
| AppKey      | string | 是   | 接入时申请的 app_key              |
| AppSecret   | string | 是   | 接入时申请的 app_secret           |
| PublicKey   | string | 是   | RSA 公钥，支持 Base64 编码或 PEM 格式 |
| Gateway     | string | 否   | API 网关地址，默认为联调环境      |
| ConsignorId | string | 否   | 货主ID，用于多货主场景            |
| Timeout     | int    | 否   | HTTP 请求超时时间（秒），默认 30 秒 |

**PublicKey 格式说明：**

SDK 支持两种公钥格式：

1. **Base64 编码格式**（推荐，中储智运默认提供）：
   ```
   MFwwDQ******8CAwEAAQ==
   ```

2. **PEM 格式**：
   ```
   -----BEGIN PUBLIC KEY-----
   MFwwDQ******8CAwEAAQ==
   -----END PUBLIC KEY-----
   ```

### Response 响应结构

| 字段    | 类型        | 说明     |
| ------- | ----------- | -------- |
| Code    | string      | 返回码   |
| Message | string      | 返回消息 |
| Data    | any | 返回数据 |

### 错误码说明

返回码共 4 位，其中前 2 位代表系统码，后 2 位代表错误码：

- `0000`：成功
- 系统码 `00`：开放平台
- 系统码 `10`：订单接口
- 其他错误码请参考开放平台文档

## 接入流程

1. 联系中储进行对接申请（现阶段支持货主方、第三方对接）
2. 获取 appKey、appSecret、publicKey 和联调地址
3. 建立技术讨论组
4. 使用测试信息进行联调测试
5. 申请正式接口并提供服务器地址（白名单）
6. 正式环境联调
7. 完成对接

## 开放平台文档

详细的 API 接口文档请访问：[https://oerp.zczy56.com/apiFile](https://oerp.zczy56.com/apiFile)

## 注意事项

1. params 参数可以为空（nil）
2. 签名规则：按参数名 ASCII 顺序排序，拼接成 `key1value1key2value2...` 格式，前后加上 appSecret，进行 MD5 加密并转大写
3. appSecret 需要使用提供的 RSA 公钥进行加密后传输
4. 时间戳使用 Unix 毫秒时间戳，服务端允许 5 分钟误差
5. 请求方法必须使用 POST，编码格式为 `application/x-www-form-urlencoded; charset=UTF-8`

## 许可证

MIT License
