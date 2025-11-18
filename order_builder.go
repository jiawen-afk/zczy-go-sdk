package zczy

// NewOrderInfoBuilder 创建订单信息构建器
func NewOrderInfoBuilder() *OrderInfoBuilder {
	return &OrderInfoBuilder{info: &OrderInfo{}}
}

// OrderInfoBuilder 订单信息构建器
type OrderInfoBuilder struct {
	info *OrderInfo
}

func (b *OrderInfoBuilder) SetOrderModel(model string) *OrderInfoBuilder {
	b.info.OrderModel = model
	return b
}

func (b *OrderInfoBuilder) SetFreightType(freightType string) *OrderInfoBuilder {
	b.info.FreightType = freightType
	return b
}

func (b *OrderInfoBuilder) SetSelfComment(comment string) *OrderInfoBuilder {
	b.info.SelfComment = comment
	return b
}

func (b *OrderInfoBuilder) SetContact(name, phone string) *OrderInfoBuilder {
	b.info.ContactName = name
	b.info.ContactPhone = phone
	return b
}

func (b *OrderInfoBuilder) SetVehicle(vehicleType, vehicleLength string) *OrderInfoBuilder {
	b.info.VehicleType = vehicleType
	b.info.VehicleLength = vehicleLength
	return b
}

func (b *OrderInfoBuilder) SetTimeSchedule(despatchStart, despatchEnd, receiveDate string) *OrderInfoBuilder {
	b.info.DespatchStart = despatchStart
	b.info.DespatchEnd = despatchEnd
	b.info.ReceiveDate = receiveDate
	return b
}

func (b *OrderInfoBuilder) SetExpectTime(expectTime string) *OrderInfoBuilder {
	b.info.ExpectTime = expectTime
	return b
}

func (b *OrderInfoBuilder) SetTotalAmount(amount string) *OrderInfoBuilder {
	b.info.TotalAmount = amount
	return b
}

func (b *OrderInfoBuilder) SetConsignorNoTaxMoney(money string) *OrderInfoBuilder {
	b.info.ConsignorNoTaxMoney = money
	return b
}

func (b *OrderInfoBuilder) SetCargoMoney(money string) *OrderInfoBuilder {
	b.info.CargoMoney = money
	return b
}

func (b *OrderInfoBuilder) SetPrompt(prompt string) *OrderInfoBuilder {
	b.info.Prompt = prompt
	return b
}

func (b *OrderInfoBuilder) SetSettleBasis(basis string) *OrderInfoBuilder {
	b.info.SettleBasis = basis
	return b
}

func (b *OrderInfoBuilder) SetInterceptPrice(price string) *OrderInfoBuilder {
	b.info.InterceptPrice = price
	return b
}

func (b *OrderInfoBuilder) SetUrgent(urgent bool) *OrderInfoBuilder {
	if urgent {
		b.info.UrgentFlag = "是"
	} else {
		b.info.UrgentFlag = "否"
	}
	return b
}

func (b *OrderInfoBuilder) SetAdvance(enabled bool, ratio string) *OrderInfoBuilder {
	if enabled {
		b.info.AdvanceFlag = "是"
		b.info.AdvanceRatio = ratio
	} else {
		b.info.AdvanceFlag = "否"
	}
	return b
}

func (b *OrderInfoBuilder) SetReceipt(enabled bool) *OrderInfoBuilder {
	if enabled {
		b.info.ReceiptFlag = "是"
	} else {
		b.info.ReceiptFlag = "否"
	}
	return b
}

func (b *OrderInfoBuilder) SetPolicy(policyFlag string) *OrderInfoBuilder {
	b.info.PolicyFlag = policyFlag
	return b
}

func (b *OrderInfoBuilder) SetOilCard(enabled bool, oilRatio, gasPercent string) *OrderInfoBuilder {
	if enabled {
		b.info.SupportSdOilCardFlag = "是"
		b.info.OilCardRatio = oilRatio
		b.info.GasPercent = gasPercent
	} else {
		b.info.SupportSdOilCardFlag = "否"
	}
	return b
}

func (b *OrderInfoBuilder) SetOilCardFixed(enabled bool, oilFixed, gasFixed string) *OrderInfoBuilder {
	if enabled {
		b.info.SupportSdOilCardFlag = "是"
		b.info.OilFixedCredit = oilFixed
		b.info.GasFixedCredit = gasFixed
	} else {
		b.info.SupportSdOilCardFlag = "否"
	}
	return b
}

func (b *OrderInfoBuilder) SetAdvisoryPhones(pickPhone, settlementPhone string) *OrderInfoBuilder {
	b.info.PickOrderAdvisoryPhone = pickPhone
	b.info.SettlementAdvisoryPhone = settlementPhone
	return b
}

func (b *OrderInfoBuilder) SetOrderMarking(marking string) *OrderInfoBuilder {
	b.info.OrderMarking = marking
	return b
}

func (b *OrderInfoBuilder) Build() *OrderInfo {
	return b.info
}

// NewCargoInfoBuilder 创建货物信息构建器
func NewCargoInfoBuilder() *CargoInfoBuilder {
	return &CargoInfoBuilder{info: &CargoInfo{}}
}

// CargoInfoBuilder 货物信息构建器
type CargoInfoBuilder struct {
	info *CargoInfo
}

func (b *CargoInfoBuilder) SetCargoName(name string) *CargoInfoBuilder {
	b.info.CargoName = name
	return b
}

func (b *CargoInfoBuilder) SetCargoVersion(version string) *CargoInfoBuilder {
	b.info.CargoVersion = version
	return b
}

func (b *CargoInfoBuilder) SetCargoCategory(category string) *CargoInfoBuilder {
	b.info.CargoCategory = category
	return b
}

func (b *CargoInfoBuilder) SetWeight(weight string) *CargoInfoBuilder {
	b.info.Weight = weight
	return b
}

func (b *CargoInfoBuilder) SetDimensions(length, width, height string) *CargoInfoBuilder {
	b.info.CargoLength = length
	b.info.CargoWidth = width
	b.info.CargoHeight = height
	return b
}

func (b *CargoInfoBuilder) SetPack(pack string) *CargoInfoBuilder {
	b.info.Pack = pack
	return b
}

func (b *CargoInfoBuilder) SetWarehouse(name, location string) *CargoInfoBuilder {
	b.info.WarehouseName = name
	b.info.WarehouseLocation = location
	return b
}

func (b *CargoInfoBuilder) Build() *CargoInfo {
	return b.info
}

// NewOrderAddressInfoBuilder 创建收发货信息构建器
func NewOrderAddressInfoBuilder() *OrderAddressInfoBuilder {
	return &OrderAddressInfoBuilder{info: &OrderAddressInfo{}}
}

// OrderAddressInfoBuilder 收发货信息构建器
type OrderAddressInfoBuilder struct {
	info *OrderAddressInfo
}

func (b *OrderAddressInfoBuilder) SetDespatchCompany(name string) *OrderAddressInfoBuilder {
	b.info.DespatchCompanyName = name
	return b
}

func (b *OrderAddressInfoBuilder) SetDespatchContact(company, name, mobile, backupMobile string) *OrderAddressInfoBuilder {
	b.info.DespatchCompanyName = company
	b.info.DespatchName = name
	b.info.DespatchMobile = mobile
	b.info.DespatchBackupMobile = backupMobile
	return b
}

func (b *OrderAddressInfoBuilder) SetDespatchAddress(pro, city, dis, place string) *OrderAddressInfoBuilder {
	b.info.DespatchPro = pro
	b.info.DespatchCity = city
	b.info.DespatchDis = dis
	b.info.DespatchPlace = place
	return b
}

func (b *OrderAddressInfoBuilder) SetDeliverCompany(name string) *OrderAddressInfoBuilder {
	b.info.DeliverCompanyName = name
	return b
}

func (b *OrderAddressInfoBuilder) SetDeliverContact(name, mobile, backupMobile string) *OrderAddressInfoBuilder {
	b.info.DeliverName = name
	b.info.DeliverMobile = mobile
	b.info.DeliverBackupMobile = backupMobile
	return b
}

func (b *OrderAddressInfoBuilder) SetDeliverAddress(pro, city, dis, place string) *OrderAddressInfoBuilder {
	b.info.DeliverPro = pro
	b.info.DeliverCity = city
	b.info.DeliverDis = dis
	b.info.DeliverPlace = place
	return b
}

func (b *OrderAddressInfoBuilder) Build() *OrderAddressInfo {
	return b.info
}

// NewCreateOrderRequestBuilder 创建订单请求构建器
func NewCreateOrderRequestBuilder() *CreateOrderRequestBuilder {
	return &CreateOrderRequestBuilder{
		req: &CreateOrderRequest{CargoList: []CargoInfo{}},
	}
}

// CreateOrderRequestBuilder 创建订单请求构建器
type CreateOrderRequestBuilder struct {
	req *CreateOrderRequest
}

func (b *CreateOrderRequestBuilder) SetOrderInfo(orderInfo OrderInfo) *CreateOrderRequestBuilder {
	b.req.OrderInfo = orderInfo
	return b
}

func (b *CreateOrderRequestBuilder) AddCargo(cargo CargoInfo) *CreateOrderRequestBuilder {
	b.req.CargoList = append(b.req.CargoList, cargo)
	return b
}

func (b *CreateOrderRequestBuilder) SetCargoList(cargoList []CargoInfo) *CreateOrderRequestBuilder {
	b.req.CargoList = cargoList
	return b
}

func (b *CreateOrderRequestBuilder) SetOrderAddressInfo(addressInfo OrderAddressInfo) *CreateOrderRequestBuilder {
	b.req.OrderAddressInfo = addressInfo
	return b
}

func (b *CreateOrderRequestBuilder) SetOrderReceiptInfo(receiptInfo *OrderReceiptInfo) *CreateOrderRequestBuilder {
	b.req.OrderReceiptInfo = receiptInfo
	return b
}

func (b *CreateOrderRequestBuilder) Build() *CreateOrderRequest {
	return b.req
}
