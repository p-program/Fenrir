package model

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
