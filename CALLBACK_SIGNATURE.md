# 回调接口签名规范

本文档说明中储智运平台在回调果子熟了回调接口时需要遵循的签名加密规范。

## 回调请求格式

### HTTP 请求规范

- **请求方法**: POST
- **Content-Type**: application/json
- **编码格式**: UTF-8

### 请求参数结构

```json
{
  "app_key": "your_app_key",
  "timestamp": "1737187200",
  "sign": "A1B2C3D4E5F6...",
  "data": "{\"orderId\":\"102019010101018811\",\"orderModel\":\"0\",...}"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| app_key | String | 是 | 应用标识（与接入时申请的app_key一致） |
| timestamp | String | 是 | 当前Unix时间戳（秒），服务端允许±30分钟误差 |
| sign | String | 是 | 签名字符串（详见签名算法） |
| data | String | 是 | 业务数据JSON字符串 |

## 签名算法

### 签名生成步骤

1. **构建待签名参数**
   - 获取除 `sign` 外的所有参数：`app_key`、`timestamp`、`data`

2. **参数排序**
   - 将参数按照 key 的 ASCII 码升序排序

3. **拼接字符串**
   - 格式：`key1value1key2value2key3value3...`
   - 示例：`app_key{app_key值}data{data值}timestamp{timestamp值}`

4. **加盐处理**
   - 在拼接字符串前后添加 `appSecret`
   - 格式：`appSecret + 拼接字符串 + appSecret`

5. **MD5 加密**
   - 对最终字符串进行 MD5 加密
   - 将结果转为大写

### 签名算法伪代码

```
待签名参数 = {
  "app_key": "your_app_key",
  "timestamp": "1737187200",
  "data": "{...}"
}

// 按ASCII排序后拼接
拼接字符串 = "app_key" + "your_app_key" + "data" + "{...}" + "timestamp" + "1737187200"

// 加盐
签名字符串 = appSecret + 拼接字符串 + appSecret

// MD5加密并转大写
sign = UPPER(MD5(签名字符串))
```

### 各语言签名示例

#### Java 示例

```java
import java.security.MessageDigest;
import java.util.Map;
import java.util.TreeMap;

public class SignatureUtil {

    public static String generateSign(String appKey, String appSecret, String timestamp, String data) {
        // 使用TreeMap自动按key排序
        Map<String, String> params = new TreeMap<>();
        params.put("app_key", appKey);
        params.put("timestamp", timestamp);
        params.put("data", data);

        // 拼接字符串
        StringBuilder sb = new StringBuilder();
        for (Map.Entry<String, String> entry : params.entrySet()) {
            sb.append(entry.getKey()).append(entry.getValue());
        }

        // 加盐
        String signStr = appSecret + sb.toString() + appSecret;

        // MD5加密并转大写
        return md5(signStr).toUpperCase();
    }

    private static String md5(String str) {
        try {
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] bytes = md.digest(str.getBytes("UTF-8"));
            StringBuilder sb = new StringBuilder();
            for (byte b : bytes) {
                sb.append(String.format("%02x", b));
            }
            return sb.toString();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}
```


## 签名验证流程

果子熟了回调接口收到回调后，会执行以下验证流程：

1. **验证 app_key**：确认是否与预期一致
2. **验证时间戳**：检查是否在有效时间范围内（±30分钟）
3. **验证签名**：使用相同算法重新计算签名，与传入的 sign 对比

**只有所有验证通过后，才会处理业务数据。**

## 完整示例

### 摘单通知回调

假设参数如下：
- app_key: `test_app_key`
- app_secret: `test_app_secret`
- timestamp: `1737187200`
- data（业务数据）:
```json
{
  "orderId": "102019010101018811",
  "orderModel": "0",
  "carrierName": "张三",
  "driverUserName": "李四",
  "plateNumber": "苏A12345",
  "weight": "12.0"
}
```

### 签名计算过程

1. **data JSON字符串**（需要压缩，无格式化空格）:
```
{"orderId":"102019010101018811","orderModel":"0","carrierName":"张三","driverUserName":"李四","plateNumber":"苏A12345","weight":"12.0"}
```

2. **参数排序拼接**:
```
app_keytest_app_keydata{"orderId":"102019010101018811","orderModel":"0","carrierName":"张三","driverUserName":"李四","plateNumber":"苏A12345","weight":"12.0"}timestamp1737187200
```

3. **加盐**:
```
test_app_secretapp_keytest_app_keydata{"orderId":"102019010101018811","orderModel":"0","carrierName":"张三","driverUserName":"李四","plateNumber":"苏A12345","weight":"12.0"}timestamp1737187200test_app_secret
```

4. **MD5并转大写**:
```
签名结果（示例）: A1B2C3D4E5F6789...
```

5. **最终请求体**:
```json
{
  "app_key": "test_app_key",
  "timestamp": "1737187200",
  "sign": "A1B2C3D4E5F6789...",
  "data": "{\"orderId\":\"102019010101018811\",\"orderModel\":\"0\",\"carrierName\":\"张三\",\"driverUserName\":\"李四\",\"plateNumber\":\"苏A12345\",\"weight\":\"12.0\"}"
}
```

## 注意事项

1. **时间戳格式**：必须是 Unix 时间戳（毫秒）
2. **data 字段**：必须是 JSON 字符串（压缩格式，无多余空格）
3. **字符编码**：统一使用 UTF-8
4. **签名大小写**：最终签名必须是大写
5. **参数排序**：严格按照 ASCII 码升序排序
6. **不包含 sign**：计算签名时，不包含 sign 字段本身

## 响应格式

果子熟了回调接口处理完回调后，会返回以下格式的JSON响应：

### 成功响应
```json
{
  "code": "0000",
  "message": "success"
}
```

### 失败响应
```json
{
  "code": "1001",
  "message": "签名验证失败"
}
```

## 联系方式

如有签名验证相关问题，请联系技术支持进行排查。
