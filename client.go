package zczy

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultGateway 联调环境网关地址
	DefaultGateway = "https://seal.zczy100.com/zczy-erp/api"
	ProdGateway = "https://connect.zczy56.com/zczy-erp/api"
	// Format 数据格式，暂时只支持json
	Format = "json"
	// SignMethod 签名方法，暂时只支持md5
	SignMethod = "md5"
	// Version API版本，固定3.0
	Version = "3.0"
)

// Client 中储智运SDK客户端
type Client struct {
	appKey      string
	appSecret   string
	publicKey   string
	gateway     string
	consignorId string
	httpClient  *http.Client
}

// Config 客户端配置
type Config struct {
	AppKey      string // 接入时申请的app_key
	AppSecret   string // 接入时申请的app_secret
	PublicKey   string // RSA公钥，用于加密appSecret
	Gateway     string // API网关地址，默认为联调环境
	ConsignorId string // 货主ID（可选）
	Timeout     int    // HTTP请求超时时间（秒），默认30秒
}

// Response API响应结构
type Response struct {
	Code    string `json:"code"`    // 返回码
	Message string `json:"message"` // 返回消息
	Result  any    `json:"result"`  // 返回数据
}

// NewClient 创建SDK客户端
func NewClient(config *Config) (*Client, error) {
	if config.AppKey == "" {
		return nil, errors.New("appKey is required")
	}
	if config.AppSecret == "" {
		return nil, errors.New("appSecret is required")
	}
	if config.PublicKey == "" {
		return nil, errors.New("publicKey is required")
	}

	gateway := config.Gateway
	if gateway == "" {
		gateway = DefaultGateway
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	return &Client{
		appKey:      config.AppKey,
		appSecret:   config.AppSecret,
		publicKey:   config.PublicKey,
		gateway:     gateway,
		consignorId: config.ConsignorId,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

// SetGateway 设置网关地址（用于切换联调环境和正式环境）
func (c *Client) SetGateway(gateway string) {
	c.gateway = gateway
}

// SetConsignorId 设置货主ID
func (c *Client) SetConsignorId(consignorId string) {
	c.consignorId = consignorId
}

// Execute 执行API调用（POST请求）
func (c *Client) Execute(method string, params any) (*Response, error) {
	// 构建请求参数
	reqParams, err := c.buildRequestParams(method, params)
	if err != nil {
		return nil, fmt.Errorf("build request params error: %w", err)
	}

	// 发送HTTP POST请求
	resp, err := c.doRequest(reqParams)
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}

	return resp, nil
}

// ExecuteGet 执行API调用（GET请求）
func (c *Client) ExecuteGet(method string, params any) (*Response, error) {
	// 构建请求参数
	reqParams, err := c.buildRequestParams(method, params)
	if err != nil {
		return nil, fmt.Errorf("build request params error: %w", err)
	}

	// 发送HTTP GET请求
	resp, err := c.doGetRequest(reqParams)
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}

	return resp, nil
}

// buildRequestParams 构建请求参数
func (c *Client) buildRequestParams(method string, params any) (map[string]string, error) {
	// 使用Unix毫秒时间戳（API要求毫秒级）
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// 转换params为字符串
	var paramsStr string
	if params != nil {
		// 检查是否是map[string]string类型（用于轨迹URL等特殊情况）
		if paramMap, ok := params.(map[string]string); ok {
			// 直接使用params字段的值
			paramsStr = paramMap["params"]
		} else {
			// 正常情况：将params序列化为JSON
			paramsBytes, err := json.Marshal(params)
			if err != nil {
				return nil, fmt.Errorf("marshal params error: %w", err)
			}
			paramsStr = string(paramsBytes)
		}
	} else {
		paramsStr = ""
	}

	// 构建参数map（用于签名）
	signParams := map[string]string{
		"appKey":      c.appKey,
		"method":      method,
		"format":      Format,
		"timestamp":   timestamp,
		"sign_method": SignMethod,
		"version":     Version,
		"params":      paramsStr,
	}

	// 生成签名
	sign := c.generateSign(signParams)

	// 加密appSecret
	encryptedSecret, err := c.encryptAppSecret()
	if err != nil {
		return nil, fmt.Errorf("encrypt appSecret error: %w", err)
	}

	// 构建最终请求参数
	requestParams := map[string]string{
		"appKey":      c.appKey,
		"appSecret":   encryptedSecret,
		"method":      method,
		"format":      Format,
		"timestamp":   timestamp,
		"sign":        sign,
		"sign_method": SignMethod,
		"version":     Version,
		"params":      paramsStr,
	}

	// 如果配置了货主ID，添加到请求参数中
	if c.consignorId != "" {
		requestParams["consignorId"] = c.consignorId
	}

	return requestParams, nil
}

// generateSign 生成签名
// 签名规则：按参数名ASCII顺序排序，拼接成key1value1key2value2...格式，
// 前后加上appSecret，进行MD5加密，转大写
func (c *Client) generateSign(params map[string]string) string {
	// 获取所有key并排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接字符串
	var builder strings.Builder
	builder.WriteString(c.appSecret)
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(params[key])
	}
	builder.WriteString(c.appSecret)

	// MD5加密并转大写
	hash := md5.Sum([]byte(builder.String()))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// encryptAppSecret 使用RSA公钥加密appSecret
func (c *Client) encryptAppSecret() (string, error) {
	var rsaPub *rsa.PublicKey

	// 尝试解析PEM格式的公钥
	block, _ := pem.Decode([]byte(c.publicKey))
	if block != nil {
		// PEM格式公钥
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			// 尝试PKCS1格式
			pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return "", fmt.Errorf("failed to parse PEM public key: %w", err)
			}
			var ok bool
			rsaPub, ok = pub.(*rsa.PublicKey)
			if !ok {
				return "", errors.New("not RSA public key")
			}
		} else {
			var ok bool
			rsaPub, ok = pub.(*rsa.PublicKey)
			if !ok {
				return "", errors.New("not RSA public key")
			}
		}
	} else {
		// 尝试Base64编码的原始公钥数据
		keyBytes, err := base64.StdEncoding.DecodeString(c.publicKey)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 public key: %w", err)
		}

		// 尝试PKIX格式
		pub, err := x509.ParsePKIXPublicKey(keyBytes)
		if err != nil {
			// 尝试PKCS1格式
			pub, err = x509.ParsePKCS1PublicKey(keyBytes)
			if err != nil {
				return "", fmt.Errorf("failed to parse public key: %w", err)
			}
			var ok bool
			rsaPub, ok = pub.(*rsa.PublicKey)
			if !ok {
				return "", errors.New("not RSA public key")
			}
		} else {
			var ok bool
			rsaPub, ok = pub.(*rsa.PublicKey)
			if !ok {
				return "", errors.New("not RSA public key")
			}
		}
	}

	// RSA加密
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(c.appSecret))
	if err != nil {
		return "", fmt.Errorf("RSA encrypt error: %w", err)
	}

	// Base64编码
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// buildTrackURL 构建车辆轨迹网址
func (c *Client) buildTrackURL(params map[string]string) (string, error) {
	// 构建GET请求URL（轨迹接口使用 /zczy-erp/html 路径）
	baseURL := strings.Replace(c.gateway, "/zczy-erp/api", "/zczy-erp/html", 1)

	// 构建查询参数
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Set(key, value)
	}

	// 拼接完整URL
	fullURL := baseURL + "?" + queryParams.Encode()
	return fullURL, nil
}

// doRequest 发送HTTP POST请求
func (c *Client) doRequest(params map[string]string) (*Response, error) {
	// 构建form数据
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, value)
	}

	// 创建请求
	req, err := http.NewRequest("POST", c.gateway, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request error: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	// 解析响应
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response error: %w, body: %s", err, string(body))
	}

	return &result, nil
}

// doGetRequest 发送HTTP GET请求
func (c *Client) doGetRequest(params map[string]string) (*Response, error) {
	// 构建GET请求URL（轨迹接口使用 /zczy-erp/html 路径）
	baseURL := strings.Replace(c.gateway, "/zczy-erp/api", "/zczy-erp/html", 1)

	// 构建查询参数
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Set(key, value)
	}

	// 拼接完整URL
	fullURL := baseURL + "?" + queryParams.Encode()

	// 创建GET请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request error: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	// 解析响应
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response error: %w, body: %s", err, string(body))
	}

	return &result, nil
}

// IsSuccess 判断响应是否成功
func (r *Response) IsSuccess() bool {
	return r.Code == "0000"
}

// GetData 获取响应数据并反序列化到指定类型
func (r *Response) GetData(v any) error {
	if !r.IsSuccess() {
		return fmt.Errorf("api error: code=%s, message=%s", r.Code, r.Message)
	}

	dataBytes, err := json.Marshal(r.Result)
	if err != nil {
		return fmt.Errorf("marshal data error: %w", err)
	}

	if err := json.Unmarshal(dataBytes, v); err != nil {
		return fmt.Errorf("unmarshal data error: %w", err)
	}

	return nil
}
