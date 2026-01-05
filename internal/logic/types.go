package logic

// WalletChargeRequest 钱包充值请求
type WalletChargeRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

// BindPlateRequest 绑定餐盘请求
type BindPlateRequest struct {
	UserID  string `json:"user_id"`
	PlateID string `json:"plate_id"`
}

// UnbindPlateRequest 解绑餐盘请求
type UnbindPlateRequest struct {
	UserID  string `json:"user_id"`
	PlateID string `json:"plate_id"`
}

// OrderFoodRequest 订单食物请求
type OrderFoodRequest struct {
	FoodID string  `json:"food_id"`
	Weight float64 `json:"weight,optional"`
}

// OrderRequest 点餐请求
type OrderRequest struct {
	UserID  string             `json:"user_id"`
	PlateID string             `json:"plate_id"`
	Foods   []OrderFoodRequest `json:"foods"`
}

// UserOrderListRequest 获取用户订单列表请求
type UserOrderListRequest struct {
	UserID   string `json:"user_id"`
	Page     int    `json:"page,optional,default=1"`
	PageSize int    `json:"page_size,optional,default=10"`
}

// WorkerExceptionRequest 工作人员异常处理请求
type WorkerExceptionRequest struct {
	WorkerID  string `json:"worker_id"`
	PlateID   string `json:"plate_id,optional"`
	Exception string `json:"exception"`
	Action    string `json:"action"`
}

// GCProcessRequest GC处理请求
type GCProcessRequest struct {
	PlateID string `json:"plate_id"`
	Type    string `json:"type"` // "plate" or "food_waste"
}
