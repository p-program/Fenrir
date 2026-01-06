package model

import (
	"net/http"
	"time"
)

// 此文件保留原有的业务逻辑方法定义
// 数据模型定义已迁移到 model/restaurant.go

// PlateManager 餐盘管理器（业务对象）
type PlateManager struct {
	Plates     []Plate // 管理的餐盘列表
	Tablewares []Tableware
}

// Tableware 餐具
// 餐具包括: 碗, 盘, 筷子, 勺子等
type Tableware struct{}

// Cauldron 装汤的大锅
type Cauldron struct {
	Soup []Soup
}

// Soup 汤
type Soup struct{}

// FoodGC 食物残渣处理（业务方法）
func (w *Worker) FoodGC() {}

// Guide 引导用户使用这套系统，以及异常处理（业务方法）
func (w *Worker) Guide() {}

// Binding 将用户和餐盘绑定在一起（业务方法）
func (u *User) Binding() {}

// Unbinding 将用户和餐盘解绑,分手动和自动2种方式
// 手动解绑: 用户主动解绑
// 自动解绑: 用户未使用餐盘超过15~20分钟（用餐中）, 系统自动解绑
func (u *User) Unbinding() {}

// Order 用户点餐（业务方法）
func (u *User) Order() *Food {
	return &Food{}
}

// Charge 用户充值（业务方法）
func (u *User) Charge(money float64) *Food {
	// Get my wallet
	wallet := &Wallet{UserID: u.ID}
	wallet.Charge(money)
	return &Food{}
}

// Charge 钱包充值（业务方法）
func (w *Wallet) Charge(money float64) {
	w.Balance += money
}

type APIResponse struct {
	Cost    time.Duration `json:"cost,omitempty"`  // 处理耗时
	Code    int           `json:"code"`            // 业务状态码（如 0 表示成功）
	Message string        `json:"message"`         // 消息提示
	Data    interface{}   `json:"data,omitempty"`  // 返回数据体，可为任意结构
	Error   string        `json:"error,omitempty"` // 可选错误描述（一般调试用）
}

func NewErrorAPIResponse(cost time.Duration, msg string) APIResponse {
	return APIResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Cost:    cost,
		// Error:   "An unexpected error occurred",
	}
}

func NewSuccessAPIResponse(cost time.Duration, msg string) APIResponse {
	return APIResponse{
		Code:    http.StatusOK,
		Message: msg,
		Cost:    cost,
		// Data:    nil,
	}
}

type Hermes interface {
	Translate(source Language, location Location) (target []Language, err error)
}

type Language struct {
	Word     string
	Location Location
}

// Location 简化为地球上的位置
type Location struct {
	Latitude  float64 // 纬度
	Longitude float64 // 经度
}
