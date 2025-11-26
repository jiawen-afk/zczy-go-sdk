package zczy

import (
	"encoding/json"
	"testing"
)

// 测试客户端创建
func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "正常创建客户端",
			config: &Config{
				AppKey:    "test_key",
				AppSecret: "test_secret",
				PublicKey: "test_public_key",
				Gateway:   DefaultGateway,
				Timeout:   30,
			},
			wantErr: false,
		},
		{
			name: "缺少AppKey",
			config: &Config{
				AppSecret: "test_secret",
				PublicKey: "test_key",
			},
			wantErr: true,
		},
		{
			name: "缺少AppSecret",
			config: &Config{
				AppKey:    "test_key",
				PublicKey: "test_key",
			},
			wantErr: true,
		},
		{
			name: "缺少PublicKey",
			config: &Config{
				AppKey:    "test_key",
				AppSecret: "test_secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 测试签名生成
func TestGenerateSign(t *testing.T) {
	client := &Client{
		appSecret: "testSecret",
	}

	params := map[string]string{
		"appKey":      "yourappKey",
		"format":      "json",
		"method":      "test",
		"params":      "test",
		"sign_method": "md5",
		"timestamp":   "1367819523",
		"version":     "3.0",
	}

	sign := client.generateSign(params)

	// 验证签名是否为32位大写字符串
	if len(sign) != 32 {
		t.Errorf("签名长度应为32位，实际为%d位", len(sign))
	}

	// 验证是否全部为大写
	for _, c := range sign {
		if c >= 'a' && c <= 'z' {
			t.Errorf("签名应为大写，但包含小写字符: %c", c)
		}
	}

	t.Logf("生成的签名: %s", sign)
}

// 测试参数排序
func TestParamsSorting(t *testing.T) {
	client := &Client{
		appSecret: "secret",
	}

	params := map[string]string{
		"z_param": "z",
		"a_param": "a",
		"m_param": "m",
	}

	sign1 := client.generateSign(params)
	sign2 := client.generateSign(params)

	// 验证相同参数生成的签名一致
	if sign1 != sign2 {
		t.Errorf("相同参数应生成相同签名，但得到不同结果: %s vs %s", sign1, sign2)
	}
}

// 测试Response结构
func TestResponse(t *testing.T) {
	// 测试成功响应
	resp := &Response{
		Code:    "0000",
		Message: "success",
		Data:    map[string]any{"key": "value"},
	}

	if !resp.IsSuccess() {
		t.Errorf("IsSuccess() 应返回true，实际返回false")
	}

	// 测试失败响应
	resp2 := &Response{
		Code:    "1001",
		Message: "error",
		Data:    nil,
	}

	if resp2.IsSuccess() {
		t.Errorf("IsSuccess() 应返回false，实际返回true")
	}
}

// 测试GetData方法
func TestResponseGetData(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	// 测试成功解析
	resp := &Response{
		Code:    "0000",
		Message: "success",
		Data: map[string]any{
			"name":  "test",
			"value": 123,
		},
	}

	var data TestData
	err := resp.GetData(&data)
	if err != nil {
		t.Errorf("GetData() 解析失败: %v", err)
	}

	if data.Name != "test" || data.Value != 123 {
		t.Errorf("GetData() 解析结果不正确: %+v", data)
	}

	// 测试错误响应
	resp2 := &Response{
		Code:    "1001",
		Message: "error",
		Data:    nil,
	}

	var data2 TestData
	err = resp2.GetData(&data2)
	if err == nil {
		t.Errorf("GetData() 应返回错误，实际返回nil")
	}
}

// 测试JSON序列化
func TestJSONMarshal(t *testing.T) {
	params := map[string]any{
		"orderId": "123456",
		"status":  1,
		"tags":    []string{"tag1", "tag2"},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Errorf("JSON序列化失败: %v", err)
	}

	t.Logf("序列化结果: %s", string(data))
}

// 测试SetGateway
func TestSetGateway(t *testing.T) {
	client := &Client{
		gateway: DefaultGateway,
	}

	newGateway := "https://production.example.com/api"
	client.SetGateway(newGateway)

	if client.gateway != newGateway {
		t.Errorf("SetGateway() 设置失败，期望=%s，实际=%s", newGateway, client.gateway)
	}
}

// 测试Base64格式公钥加密
func TestEncryptAppSecretWithBase64Key(t *testing.T) {
	// 中储智运提供的Base64格式公钥
	publicKey := "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALT8QammE81aGfzzmFj0LjHKAOWiyRLESX4fwomlvWr3nVvx4rSzKGz176M/c9UsLQFqJkA0KIk0YxDgS1QG5K8CAwEAAQ=="

	client := &Client{
		appSecret: "testSecret123",
		publicKey: publicKey,
	}

	encrypted, err := client.encryptAppSecret()
	if err != nil {
		t.Fatalf("encryptAppSecret() 失败: %v", err)
	}

	if encrypted == "" {
		t.Errorf("encryptAppSecret() 返回空字符串")
	}

	t.Logf("加密后的appSecret: %s", encrypted)
}

// 测试PEM格式公钥加密
func TestEncryptAppSecretWithPEMKey(t *testing.T) {
	// PEM格式的测试公钥
	pemKey := `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALT8QammE81aGfzzmFj0LjHKAOWiyRLE
SX4fwomlvWr3nVvx4rSzKGz176M/c9UsLQFqJkA0KIk0YxDgS1QG5K8CAwEAAQ==
-----END PUBLIC KEY-----`

	client := &Client{
		appSecret: "testSecret123",
		publicKey: pemKey,
	}

	encrypted, err := client.encryptAppSecret()
	if err != nil {
		t.Fatalf("encryptAppSecret() 失败: %v", err)
	}

	if encrypted == "" {
		t.Errorf("encryptAppSecret() 返回空字符串")
	}

	t.Logf("加密后的appSecret: %s", encrypted)
}
