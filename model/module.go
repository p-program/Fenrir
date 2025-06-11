package model

type Wallet struct {
	Money  float32
	UserID string
}

type Food struct {
	Price  float32 //总价格
	Name   string  //食物名称
	Weight float32 //重量
}

type User struct {
	UserID string
}

// Plate 餐盘
type Plate struct {
	PlateID        string
	Owner          *User    // 餐盘的主人
	WeightSensor   struct{} //食物重量传感器
	FoodClip       struct{} //餐夹
	QrCode         struct{} //二维码
	RFIDManagement struct{} //射频识别,主要是用来关联餐盘
}

type PlateManager struct {
	Plates     []Plate // 管理的餐盘列表
	tablewares []tableware
}

// Tableware 餐具
// 餐具包括: 碗, 盘, 筷子, 勺子等
type tableware struct{}

type Worker struct{}

// Cauldron 装汤的大锅
type Cauldron struct {
	Soup []Soup
}

type Soup struct{}

// FoodGC 食物残渣处理
func (w *Worker) FoodGC() {}

// Guide 引导用户使用这套系统，已经异常处理
func (w *Worker) Guide() {}

// Binding 将用户和餐盘绑定在一起
func (u *User) Binding() {}

// Unbinding 将用户和餐盘解绑,分手动和自动2种方式
// 手动解绑: 用户主动解绑
// 自动解绑: 用户未使用餐盘超过15~20分钟（用餐中）, 系统自动解绑
func (u *User) Unbinding() {}

// Order 用户点餐
func (u *User) Order() Food {
	return Food{}
}

func (u *User) Charge(money float32) Food {
	// Get my wallet
	wallet := &Wallet{UserID: u.UserID}
	wallet.Charge(money)
	return Food{}
}

// Charge 充值
func (w *Wallet) Charge(money float32) {}
