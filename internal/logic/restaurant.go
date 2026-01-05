package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"zeusro.com/gotemplate/model"
)

// RestaurantLogic 餐厅业务逻辑
type RestaurantLogic struct {
	db *gorm.DB
}

// NewRestaurantLogic 创建餐厅业务逻辑实例
func NewRestaurantLogic(db *gorm.DB) *RestaurantLogic {
	return &RestaurantLogic{db: db}
}

// ChargeWallet 钱包充值
func (l *RestaurantLogic) ChargeWallet(ctx context.Context, userID string, amount float64) (*model.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("充值金额必须大于0")
	}

	var wallet model.Wallet
	err := l.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新钱包
			wallet = model.Wallet{
				UserID:  userID,
				Balance: amount,
			}
			if err := l.db.WithContext(ctx).Create(&wallet).Error; err != nil {
				return nil, fmt.Errorf("创建钱包失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("查询钱包失败: %w", err)
		}
	} else {
		// 更新余额
		wallet.Balance += amount
		if err := l.db.WithContext(ctx).Save(&wallet).Error; err != nil {
			return nil, fmt.Errorf("更新钱包失败: %w", err)
		}
	}

	// 记录交易
	transaction := model.Transaction{
		WalletID: wallet.ID,
		Type:     "charge",
		Amount:   amount,
		Balance:  wallet.Balance,
		Remark:   "钱包充值",
	}
	if err := l.db.WithContext(ctx).Create(&transaction).Error; err != nil {
		return nil, fmt.Errorf("记录交易失败: %w", err)
	}

	return &wallet, nil
}

// GetUserInfo 获取用户信息
func (l *RestaurantLogic) GetUserInfo(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := l.db.WithContext(ctx).Preload("Wallet").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// BindPlate 绑定餐盘
func (l *RestaurantLogic) BindPlate(ctx context.Context, userID string, plateID string) (*model.Plate, error) {
	// 检查用户是否存在
	var user model.User
	if err := l.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 检查餐盘是否存在
	var plate model.Plate
	if err := l.db.WithContext(ctx).Where("id = ?", plateID).First(&plate).Error; err != nil {
		return nil, fmt.Errorf("餐盘不存在: %w", err)
	}

	// 检查餐盘是否已被绑定
	if plate.IsBound && plate.BoundUserID != userID {
		return nil, errors.New("餐盘已被其他用户绑定")
	}

	// 如果用户已有绑定的餐盘，先解绑
	var existingPlate model.Plate
	if err := l.db.WithContext(ctx).Where("bound_user_id = ? AND is_bound = ?", userID, true).First(&existingPlate).Error; err == nil {
		existingPlate.IsBound = false
		existingPlate.BoundUserID = ""
		existingPlate.BoundAt = nil
		l.db.WithContext(ctx).Save(&existingPlate)
	}

	// 绑定餐盘
	now := time.Now()
	plate.IsBound = true
	plate.BoundUserID = userID
	plate.BoundAt = &now
	plate.Status = "in_use"

	if err := l.db.WithContext(ctx).Save(&plate).Error; err != nil {
		return nil, fmt.Errorf("绑定餐盘失败: %w", err)
	}

	return &plate, nil
}

// UnbindPlate 解绑餐盘
func (l *RestaurantLogic) UnbindPlate(ctx context.Context, userID string, plateID string) error {
	var plate model.Plate
	if err := l.db.WithContext(ctx).Where("id = ? AND bound_user_id = ?", plateID, userID).First(&plate).Error; err != nil {
		return fmt.Errorf("餐盘不存在或未绑定: %w", err)
	}

	plate.IsBound = false
	plate.BoundUserID = ""
	plate.BoundAt = nil
	plate.Status = "available"

	if err := l.db.WithContext(ctx).Save(&plate).Error; err != nil {
		return fmt.Errorf("解绑餐盘失败: %w", err)
	}

	return nil
}

// GetPlateInfo 获取餐盘信息
func (l *RestaurantLogic) GetPlateInfo(ctx context.Context, plateID string) (*model.Plate, error) {
	var plate model.Plate
	err := l.db.WithContext(ctx).Preload("BoundUser").Where("id = ?", plateID).First(&plate).Error
	if err != nil {
		return nil, fmt.Errorf("查询餐盘失败: %w", err)
	}
	return &plate, nil
}

// GetPlateList 获取餐盘列表
func (l *RestaurantLogic) GetPlateList(ctx context.Context, isBound *bool) ([]model.Plate, error) {
	var plates []model.Plate
	query := l.db.WithContext(ctx).Preload("BoundUser")

	if isBound != nil {
		query = query.Where("is_bound = ?", *isBound)
	}

	if err := query.Find(&plates).Error; err != nil {
		return nil, fmt.Errorf("查询餐盘列表失败: %w", err)
	}
	return plates, nil
}

// CreateOrder 创建订单
func (l *RestaurantLogic) CreateOrder(ctx context.Context, userID string, plateID string, foods []OrderFood) (*model.Order, error) {
	// 检查用户和餐盘绑定关系
	var plate model.Plate
	if err := l.db.WithContext(ctx).Where("id = ? AND bound_user_id = ? AND is_bound = ?", plateID, userID, true).First(&plate).Error; err != nil {
		return nil, fmt.Errorf("餐盘未绑定或绑定关系不正确: %w", err)
	}

	// 检查用户钱包
	var wallet model.Wallet
	if err := l.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, fmt.Errorf("用户钱包不存在: %w", err)
	}

	// 计算订单总价
	var totalPrice float64
	var orderItems []model.OrderItem

	for _, foodReq := range foods {
		var food model.Food
		if err := l.db.WithContext(ctx).Where("id = ?", foodReq.FoodID).First(&food).Error; err != nil {
			return nil, fmt.Errorf("食物不存在: %s, %w", foodReq.FoodID, err)
		}

		if !food.IsAvailable {
			return nil, fmt.Errorf("食物不可用: %s", food.Name)
		}

		weight := foodReq.Weight
		if weight <= 0 {
			weight = 100 // 默认100克
		}

		itemPrice := (food.Price / 100.0) * weight
		totalPrice += itemPrice

		orderItems = append(orderItems, model.OrderItem{
			FoodID:    food.ID,
			FoodName:  food.Name,
			Weight:    weight,
			UnitPrice: food.Price,
			Price:     itemPrice,
		})
	}

	// 检查余额
	if wallet.Balance < totalPrice {
		return nil, fmt.Errorf("余额不足，当前余额: %.2f, 需要: %.2f", wallet.Balance, totalPrice)
	}

	// 创建订单
	orderID := uuid.New().String()
	order := model.Order{
		ID:         orderID,
		UserID:     userID,
		PlateID:    plateID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	if err := l.db.WithContext(ctx).Create(&order).Error; err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 创建订单明细
	for i := range orderItems {
		orderItems[i].OrderID = orderID
	}
	if err := l.db.WithContext(ctx).Create(&orderItems).Error; err != nil {
		return nil, fmt.Errorf("创建订单明细失败: %w", err)
	}

	// 扣款
	wallet.Balance -= totalPrice
	if err := l.db.WithContext(ctx).Save(&wallet).Error; err != nil {
		return nil, fmt.Errorf("扣款失败: %w", err)
	}

	// 记录交易
	transaction := model.Transaction{
		WalletID: wallet.ID,
		Type:     "consume",
		Amount:   -totalPrice,
		Balance:  wallet.Balance,
		OrderID:  orderID,
		Remark:   "订单消费",
	}
	if err := l.db.WithContext(ctx).Create(&transaction).Error; err != nil {
		return nil, fmt.Errorf("记录交易失败: %w", err)
	}

	// 更新订单状态
	order.Status = "paid"
	if err := l.db.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 加载关联数据
	if err := l.db.WithContext(ctx).Preload("OrderItems").Preload("User").Preload("Plate").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	return &order, nil
}

// GetUserOrders 获取用户订单列表
func (l *RestaurantLogic) GetUserOrders(ctx context.Context, userID string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * pageSize

	query := l.db.WithContext(ctx).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询订单总数失败: %w", err)
	}

	// 获取列表
	if err := query.Preload("OrderItems").Preload("Plate").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("查询订单列表失败: %w", err)
	}

	return orders, total, nil
}

// GetOrderInfo 获取订单信息
func (l *RestaurantLogic) GetOrderInfo(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order
	err := l.db.WithContext(ctx).Preload("OrderItems").Preload("User").Preload("Plate").
		Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}
	return &order, nil
}

// GetPlateDepot 获取餐盘托管处信息
func (l *RestaurantLogic) GetPlateDepot(ctx context.Context, depotID string) (*model.PlateDepot, error) {
	var depot model.PlateDepot
	err := l.db.WithContext(ctx).Where("id = ?", depotID).First(&depot).Error
	if err != nil {
		return nil, fmt.Errorf("查询餐盘托管处失败: %w", err)
	}
	return &depot, nil
}

// HandleException 处理异常
func (l *RestaurantLogic) HandleException(ctx context.Context, workerID string, plateID string, exception string, action string) error {
	// 检查工作人员是否存在
	var worker model.Worker
	if err := l.db.WithContext(ctx).Where("id = ?", workerID).First(&worker).Error; err != nil {
		return fmt.Errorf("工作人员不存在: %w", err)
	}

	// 记录异常
	exceptionLog := model.ExceptionLog{
		WorkerID:  workerID,
		PlateID:   plateID,
		Exception: exception,
		Action:    action,
		Status:    "pending",
	}

	if err := l.db.WithContext(ctx).Create(&exceptionLog).Error; err != nil {
		return fmt.Errorf("记录异常失败: %w", err)
	}

	// 如果有餐盘ID，更新餐盘状态
	if plateID != "" {
		var plate model.Plate
		if err := l.db.WithContext(ctx).Where("id = ?", plateID).First(&plate).Error; err == nil {
			plate.Status = "maintenance"
			l.db.WithContext(ctx).Save(&plate)
		}
	}

	return nil
}

// ProcessGC 处理GC
func (l *RestaurantLogic) ProcessGC(ctx context.Context, plateID string, gcType string) error {
	// 检查餐盘是否存在
	var plate model.Plate
	if err := l.db.WithContext(ctx).Where("id = ?", plateID).First(&plate).Error; err != nil {
		return fmt.Errorf("餐盘不存在: %w", err)
	}

	// 创建GC处理记录
	gcLog := model.GCProcessLog{
		PlateID: plateID,
		Type:    gcType,
		Status:  "processing",
	}

	if err := l.db.WithContext(ctx).Create(&gcLog).Error; err != nil {
		return fmt.Errorf("创建GC处理记录失败: %w", err)
	}

	// 根据类型处理
	if gcType == "plate" {
		// 餐盘处理：清理、重置状态
		plate.Weight = 0
		plate.Status = "available"
		plate.IsBound = false
		plate.BoundUserID = ""
		plate.BoundAt = nil
		if err := l.db.WithContext(ctx).Save(&plate).Error; err != nil {
			return fmt.Errorf("处理餐盘失败: %w", err)
		}
	} else if gcType == "food_waste" {
		// 厨余垃圾处理：重置重量
		plate.Weight = 0
		if err := l.db.WithContext(ctx).Save(&plate).Error; err != nil {
			return fmt.Errorf("处理厨余垃圾失败: %w", err)
		}
	}

	// 更新GC处理状态
	gcLog.Status = "completed"
	if err := l.db.WithContext(ctx).Save(&gcLog).Error; err != nil {
		return fmt.Errorf("更新GC处理状态失败: %w", err)
	}

	return nil
}

// OrderFood 订单食物
type OrderFood struct {
	FoodID string
	Weight float64
}
