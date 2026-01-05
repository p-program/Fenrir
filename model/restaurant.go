package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Username  string         `gorm:"type:varchar(100);uniqueIndex" json:"username"`
	Phone     string         `gorm:"type:varchar(20);index" json:"phone,omitempty"`
	Email     string         `gorm:"type:varchar(100)" json:"email,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Wallet *Wallet `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
	Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// Wallet 钱包表
type Wallet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"user_id"`
	Balance   float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
}

// Transaction 交易记录表
type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WalletID  uint      `gorm:"index;not null" json:"wallet_id"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"` // "charge", "consume", "refund"
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Balance   float64   `gorm:"type:decimal(10,2);not null" json:"balance"` // 交易后余额
	OrderID   string    `gorm:"type:varchar(64);index" json:"order_id,omitempty"`
	Remark    string    `gorm:"type:varchar(255)" json:"remark,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Wallet *Wallet `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
}

// Plate 餐盘表
type Plate struct {
	ID          string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	QRCode      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"qr_code"`
	RFIDTag     string         `gorm:"type:varchar(255);uniqueIndex" json:"rfid_tag,omitempty"`
	Weight      float64        `gorm:"type:decimal(8,2);default:0" json:"weight"` // 当前重量（克）
	IsBound     bool           `gorm:"default:false;index" json:"is_bound"`
	BoundUserID string         `gorm:"type:varchar(64);index" json:"bound_user_id,omitempty"`
	BoundAt     *time.Time     `json:"bound_at,omitempty"`
	Status      string         `gorm:"type:varchar(20);default:'available'" json:"status"` // available, in_use, cleaning, maintenance
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	BoundUser *User   `gorm:"foreignKey:BoundUserID" json:"bound_user,omitempty"`
	Orders    []Order `gorm:"foreignKey:PlateID" json:"orders,omitempty"`
}

// Food 食物表
type Food struct {
	ID          string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Price       float64   `gorm:"type:decimal(8,2);not null" json:"price"`    // 单价（每100克）
	Category    string    `gorm:"type:varchar(50)" json:"category,omitempty"` // 菜品分类
	Description string    `gorm:"type:text" json:"description,omitempty"`
	IsAvailable bool      `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	OrderItems []OrderItem `gorm:"foreignKey:FoodID" json:"order_items,omitempty"`
}

// Order 订单表
type Order struct {
	ID         string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	UserID     string         `gorm:"type:varchar(64);index;not null" json:"user_id"`
	PlateID    string         `gorm:"type:varchar(64);index;not null" json:"plate_id"`
	TotalPrice float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, completed, cancelled
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plate      *Plate      `gorm:"foreignKey:PlateID" json:"plate,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

// OrderItem 订单明细表
type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   string    `gorm:"type:varchar(64);index;not null" json:"order_id"`
	FoodID    string    `gorm:"type:varchar(64);index;not null" json:"food_id"`
	FoodName  string    `gorm:"type:varchar(100);not null" json:"food_name"`
	Weight    float64   `gorm:"type:decimal(8,2);not null" json:"weight"`     // 重量（克）
	UnitPrice float64   `gorm:"type:decimal(8,2);not null" json:"unit_price"` // 单价
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`     // 总价
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Food  *Food  `gorm:"foreignKey:FoodID" json:"food,omitempty"`
}

// PlateDepot 餐盘托管处表
type PlateDepot struct {
	ID        string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Location  string         `gorm:"type:varchar(255)" json:"location,omitempty"`
	Capacity  int            `gorm:"default:100" json:"capacity"`
	Available int            `gorm:"default:100" json:"available"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Worker 工作人员表
type Worker struct {
	ID        string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Role      string         `gorm:"type:varchar(50);not null" json:"role"` // "staff", "manager", "gc"
	Phone     string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// ExceptionLog 异常处理记录表
type ExceptionLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	WorkerID  string         `gorm:"type:varchar(64);index;not null" json:"worker_id"`
	PlateID   string         `gorm:"type:varchar(64);index" json:"plate_id,omitempty"`
	Exception string         `gorm:"type:text;not null" json:"exception"`
	Action    string         `gorm:"type:varchar(255);not null" json:"action"`
	Status    string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, resolved
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Worker *Worker `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
	Plate  *Plate  `gorm:"foreignKey:PlateID" json:"plate,omitempty"`
}

// GCProcessLog GC处理记录表
type GCProcessLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PlateID   string         `gorm:"type:varchar(64);index;not null" json:"plate_id"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"`            // "plate", "food_waste"
	Status    string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, processing, completed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Plate *Plate `gorm:"foreignKey:PlateID" json:"plate,omitempty"`
}
