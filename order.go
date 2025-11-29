package zczy

const (
	// MethodOrderCreateMore 生成普通货(支持单货、多货)
	MethodOrderCreateMore = "zczy.api.order.create.more"
	// MethodOrderCancel 订单取消
	MethodOrderCancel = "zczy.api.order.cancel"
	// MethodReceiptConfirm 回单确认
	MethodReceiptConfirm = "zczy.api.receipt.confirm"
	// MethodVehicleTrack 获取车辆在途轨迹网址
	MethodVehicleTrack = "zczy.html.order.cordinate"
	// MethodOrderCoordinate 在途轨迹
	MethodOrderCoordinate = "zczy.api.order.cordinate"
)

// OrderInfo 订单信息
type OrderInfo struct {
	OrderModel              string `json:"orderModel"`                        // 订单类型：抢单,竞价
	FreightType             string `json:"freightType"`                       // 费用类型：包车价，单价
	SelfComment             string `json:"selfComment"`                       // 自定义编号/标签
	ContactName             string `json:"contactName"`                       // 紧急联系人
	ContactPhone            string `json:"contactPhone"`                      // 紧急联系电话
	VehicleType             string `json:"vehicleType"`                       // 车型要求
	VehicleLength           string `json:"vehicleLength"`                     // 车长要求，单位米
	UrgentFlag              string `json:"urgentFlag"`                        // 是否加急：否，是
	DespatchStart           string `json:"despatchStart"`                     // 装货开始时间
	DespatchEnd             string `json:"despatchEnd"`                       // 装货结束时间
	ReceiveDate             string `json:"receiveDate"`                       // 收货时间
	ExpectTime              string `json:"expectTime,omitempty"`              // 报价结束时间（竞价时必填）
	TotalAmount             string `json:"totalAmount,omitempty"`             // 运费（抢单时填）
	ConsignorNoTaxMoney     string `json:"consignorNoTaxMoney,omitempty"`     // 承运方预估到手价（抢单时填）
	CargoMoney              string `json:"cargoMoney"`                        // 货物价值
	Prompt                  string `json:"prompt,omitempty"`                  // 装卸货要求
	AdvanceFlag             string `json:"advanceFlag"`                       // 是否预付：否，是
	AdvanceRatio            string `json:"advanceRatio,omitempty"`            // 预付比例
	ReceiptFlag             string `json:"receiptFlag"`                       // 是否押回单：否,是
	PolicyFlag              string `json:"policyFlag"`                        // 是否购买保险
	SupportSdOilCardFlag    string `json:"supportSdOilCardFlag"`              // 是否包含油气品：否,是
	OilCardRatio            string `json:"oilCardRatio,omitempty"`            // 油品比例
	GasPercent              string `json:"gasPercent,omitempty"`              // 汽品比例
	OilFixedCredit          string `json:"oilFixedCredit,omitempty"`          // 油品固定额度
	GasFixedCredit          string `json:"gasFixedCredit,omitempty"`          // 气品固定额度
	RuleID                  string `json:"ruleId,omitempty"`                  // 自动成交规则Id
	RuleName                string `json:"ruleName,omitempty"`                // 自动成交规则名称
	TonRuleID               string `json:"tonRuleId,omitempty"`               // 亏涨吨扣款规则id
	SettleBasis             string `json:"settleBasis"`                       // 结算依据
	PickOrderAdvisoryPhone  string `json:"pickOrderAdvisoryPhone,omitempty"`  // 摘单咨询电话
	SettlementAdvisoryPhone string `json:"settlementAdvisoryPhone,omitempty"` // 结算咨询电话
	InterceptPrice          string `json:"interceptPrice"`                    // 拦标价
	OrderMarking            string `json:"orderMarking,omitempty"`            // 运单标识
}

// CargoInfo 货物信息
type CargoInfo struct {
	CargoName         string `json:"cargoName"`                   // 货物名称
	CargoVersion      string `json:"cargoVersion,omitempty"`      // 规格型号
	CargoCategory     string `json:"cargoCategory"`               // 货物类别：重货，泡货
	Weight            string `json:"weight"`                      // 货物重量或体积
	CargoLength       string `json:"cargoLength,omitempty"`       // 规格-长
	CargoWidth        string `json:"cargoWidth,omitempty"`        // 规格-宽
	CargoHeight       string `json:"cargoHeight,omitempty"`       // 规格-高
	Pack              string `json:"pack"`                        // 包装类型
	WarehouseName     string `json:"warehouseName,omitempty"`     // 仓库名称
	WarehouseLocation string `json:"warehouseLocation,omitempty"` // 仓库位置
}

// OrderAddressInfo 收发货信息
type OrderAddressInfo struct {
	DespatchCompanyName  string `json:"despatchCompanyName"`            // 发货单位名称
	DespatchName         string `json:"despatchName"`                   // 发货人姓名
	DespatchMobile       string `json:"despatchMobile"`                 // 发货人联系电话
	DespatchBackupMobile string `json:"despatchBackupMobile,omitempty"` // 发货人备用电话
	DespatchPro          string `json:"despatchPro"`                    // 启运地省
	DespatchCity         string `json:"despatchCity"`                   // 启运地市
	DespatchDis          string `json:"despatchDis"`                    // 启运地区
	DespatchPlace        string `json:"despatchPlace"`                  // 启运地详细地址
	DeliverCompanyName   string `json:"deliverCompanyName"`             // 收货单位名称
	DeliverName          string `json:"deliverName"`                    // 收货人名称
	DeliverMobile        string `json:"deliverMobile"`                  // 收货人联系电话
	DeliverBackupMobile  string `json:"deliverBackupMobile,omitempty"`  // 收货人备用电话
	DeliverPro           string `json:"deliverPro"`                     // 目的地省
	DeliverCity          string `json:"deliverCity"`                    // 目的地市
	DeliverDis           string `json:"deliverDis"`                     // 目的地区
	DeliverPlace         string `json:"deliverPlace"`                   // 目的地详细地址
}

// OrderReceiptInfo 押回单信息
type OrderReceiptInfo struct {
	ReceiptLabel string `json:"receiptLabel,omitempty"` // 押回单标签
	ReceiptMoney string `json:"receiptMoney,omitempty"` // 回单押金
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	OrderInfo        OrderInfo        `json:"orderInfo"`                  // 订单信息
	CargoList        []CargoInfo      `json:"cargoList"`                  // 货物信息
	OrderAddressInfo OrderAddressInfo `json:"orderAddressInfo"`           // 收发货信息
	OrderReceiptInfo OrderReceiptInfo `json:"orderReceiptInfo,omitempty"` // 押回单信息（可选）
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID string `json:"orderId"` // 订单号
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID string `json:"orderId"` // 订单号
}

// ConfirmReceiptRequest 回单确认请求
type ConfirmReceiptRequest struct {
	OrderID             string `json:"orderId"`                       // 订单号
	Tonnage             string `json:"tonnage"`                       // 收货吨位
	SettleMoney         string `json:"settleMoney,omitempty"`         // 结算金额（与ConsignorNoTaxMoney二选一）
	ConsignorNoTaxMoney string `json:"consignorNoTaxMoney,omitempty"` // 承运方预估到手价（与SettleMoney二选一）
	SettleApplyFlag     string `json:"settleApplyFlag"`               // 是否提交结算申请：0-否，1-是
	Remark              string `json:"remark,omitempty"`              // 备注
}

// VehicleTrackRequest 获取车辆在途轨迹网址请求
type VehicleTrackRequest struct {
	OrderID          string `json:"orderId"`                    // 订单号
	CreatedStartTime string `json:"createdStartTime,omitempty"` // 开始时间（格式：2021-08-02 12:20）
	CreatedEndTime   string `json:"createdEndTime,omitempty"`   // 结束时间（格式：2021-08-02 13:20）
}

// VehicleTrackResponse 获取车辆在途轨迹网址响应
type VehicleTrackResponse struct {
	URL string `json:"url"` // 轨迹网址
}

// Coordinate 位置坐标信息
type Coordinate struct {
	Address     string `json:"address"`     // 位置
	Longitude   string `json:"longitude"`   // 经度
	Latitude    string `json:"latitude"`    // 纬度
	CreatedTime string `json:"createdTime"` // 定位时间
	Type        string `json:"type"`        // 地图类型 1-高德
}

// OrderCoordinateRequest 在途轨迹请求
type OrderCoordinateRequest struct {
	OrderID          string `json:"orderId"`                    // 订单号
	CreatedStartTime string `json:"createdStartTime,omitempty"` // 开始时间（格式：2021-08-02 12:20）
	CreatedEndTime   string `json:"createdEndTime,omitempty"`   // 结束时间（格式：2021-08-02 13:20）
}

// OrderCoordinateResponse 在途轨迹响应
type OrderCoordinateResponse struct {
	OrderID        string       `json:"orderId"`       // 订单号
	DriverName     string       `json:"driverName"`    // 司机名称
	PlateNumber    string       `json:"plateNumber"`   // 车牌号
	DriverMobile   string       `json:"driverMobile"`  // 联系方式
	CoordinateList []Coordinate `json:"cordinateList"` // 位置信息列表（注意：API使用cordinate拼写）
}

// CreateOrder 创建普通货订单（支持单货、多货）
func (c *Client) CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
	resp, err := c.Execute(MethodOrderCreateMore, req)
	if err != nil {
		return nil, err
	}

	var result CreateOrderResponse
	if err := resp.GetData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(orderID string) error {
	req := &CancelOrderRequest{
		OrderID: orderID,
	}

	resp, err := c.Execute(MethodOrderCancel, req)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return nil
	}

	return nil
}

// ConfirmReceipt 回单确认
func (c *Client) ConfirmReceipt(req *ConfirmReceiptRequest) error {
	resp, err := c.Execute(MethodReceiptConfirm, req)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return nil
	}

	return nil
}

// GetVehicleTrack 获取车辆在途轨迹网址
func (c *Client) GetVehicleTrack(req *VehicleTrackRequest) (*VehicleTrackResponse, error) {
	// 轨迹URL的params直接使用订单号字符串，不是JSON对象
	// 使用map[string]string类型，buildRequestParams会自动识别并处理
	params := map[string]string{
		"params": req.OrderID,
	}

	// 构建请求参数（包含签名等）
	reqParams, err := c.buildRequestParams(MethodVehicleTrack, params)
	if err != nil {
		return nil, err
	}

	// 构建轨迹URL
	url, err := c.buildTrackURL(reqParams)
	if err != nil {
		return nil, err
	}

	return &VehicleTrackResponse{URL: url}, nil
}

// GetOrderCoordinate 获取订单在途轨迹坐标
func (c *Client) GetOrderCoordinate(req *OrderCoordinateRequest) (*OrderCoordinateResponse, error) {
	resp, err := c.Execute(MethodOrderCoordinate, req)
	if err != nil {
		return nil, err
	}

	var result OrderCoordinateResponse
	if err := resp.GetData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
