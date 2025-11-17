# 中储智运开放平台 Golang SDK

中储智运开放平台的官方 Golang SDK，提供便捷的 API 调用接口。

## 功能特性

- 自动签名生成（MD5）
- appSecret RSA 加密
- 支持联调环境和正式环境切换
- 完整的错误处理
- 类型安全的响应解析

## 安装

```bash
go get github.com/Jiawen-AFK/zczy-go-sdk
```

## 快速开始

### 1. 初始化客户端

```go
package main

import (
    "fmt"
    "log"
    "github.com/Jiawen-AFK/zczy-go-sdk"
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

### 3. 切换到正式环境

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

## API 参数说明

### Config 配置参数

| 参数      | 类型   | 必填 | 说明                              |
| --------- | ------ | ---- | --------------------------------- |
| AppKey    | string | 是   | 接入时申请的 app_key              |
| AppSecret | string | 是   | 接入时申请的 app_secret           |
| PublicKey | string | 是   | RSA 公钥，用于加密 appSecret      |
| Gateway   | string | 否   | API 网关地址，默认为联调环境      |
| Timeout   | int    | 否   | HTTP 请求超时时间（秒），默认 30 秒 |

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
