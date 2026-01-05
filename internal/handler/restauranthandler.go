package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zeusro.com/gotemplate/internal/logic"
	"zeusro.com/gotemplate/internal/svc"
)

type RestaurantHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRestaurantHandler(svcCtx *svc.ServiceContext) *RestaurantHandler {
	return &RestaurantHandler{
		svcCtx: svcCtx,
	}
}

// HealthCheck 健康检查
func (h *RestaurantHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
	})
}

// WalletCharge 钱包充值
func (h *RestaurantHandler) WalletCharge(w http.ResponseWriter, r *http.Request) {
	var req logic.WalletChargeRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	wallet, err := l.ChargeWallet(r.Context(), req.UserID, req.Amount)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "充值成功",
		"data": map[string]interface{}{
			"user_id": wallet.UserID,
			"balance": wallet.Balance,
		},
	})
}

// GetUserInfo 获取用户信息
func (h *RestaurantHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("用户ID不能为空"))
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	user, err := l.GetUserInfo(r.Context(), userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	var balance float64
	if user.Wallet != nil {
		balance = user.Wallet.Balance
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"balance":  balance,
		},
	})
}

// BindPlate 绑定餐盘
func (h *RestaurantHandler) BindPlate(w http.ResponseWriter, r *http.Request) {
	var req logic.BindPlateRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	plate, err := l.BindPlate(r.Context(), req.UserID, req.PlateID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "绑定成功",
		"data": map[string]interface{}{
			"plate_id":      plate.ID,
			"qr_code":       plate.QRCode,
			"is_bound":      plate.IsBound,
			"bound_user_id": plate.BoundUserID,
		},
	})
}

// UnbindPlate 解绑餐盘
func (h *RestaurantHandler) UnbindPlate(w http.ResponseWriter, r *http.Request) {
	var req logic.UnbindPlateRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	if err := l.UnbindPlate(r.Context(), req.UserID, req.PlateID); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "解绑成功",
	})
}

// GetPlateInfo 获取餐盘信息
func (h *RestaurantHandler) GetPlateInfo(w http.ResponseWriter, r *http.Request) {
	plateID := r.PathValue("plate_id")
	if plateID == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("餐盘ID不能为空"))
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	plate, err := l.GetPlateInfo(r.Context(), plateID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"plate_id":      plate.ID,
			"qr_code":       plate.QRCode,
			"weight":        plate.Weight,
			"is_bound":      plate.IsBound,
			"bound_user_id": plate.BoundUserID,
			"status":        plate.Status,
		},
	})
}

// GetPlateList 获取餐盘列表
func (h *RestaurantHandler) GetPlateList(w http.ResponseWriter, r *http.Request) {
	var isBound *bool
	if boundStr := r.URL.Query().Get("is_bound"); boundStr != "" {
		if bound, err := strconv.ParseBool(boundStr); err == nil {
			isBound = &bound
		}
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	plates, err := l.GetPlateList(r.Context(), isBound)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	var plateList []map[string]interface{}
	for _, plate := range plates {
		plateList = append(plateList, map[string]interface{}{
			"plate_id":      plate.ID,
			"qr_code":       plate.QRCode,
			"weight":        plate.Weight,
			"is_bound":      plate.IsBound,
			"bound_user_id": plate.BoundUserID,
			"status":        plate.Status,
		})
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": plateList,
	})
}

// CreateOrder 创建订单
func (h *RestaurantHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req logic.OrderRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)

	// 转换请求数据
	var orderFoods []logic.OrderFood
	for _, food := range req.Foods {
		orderFoods = append(orderFoods, logic.OrderFood{
			FoodID: food.FoodID,
			Weight: food.Weight,
		})
	}

	order, err := l.CreateOrder(r.Context(), req.UserID, req.PlateID, orderFoods)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 转换订单明细
	var foods []map[string]interface{}
	for _, item := range order.OrderItems {
		foods = append(foods, map[string]interface{}{
			"food_id": item.FoodID,
			"name":    item.FoodName,
			"weight":  item.Weight,
			"price":   item.Price,
		})
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "订单创建成功",
		"data": map[string]interface{}{
			"order_id":    order.ID,
			"user_id":     order.UserID,
			"plate_id":    order.PlateID,
			"foods":       foods,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"created_at":  order.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// GetUserOrders 获取用户订单列表
func (h *RestaurantHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	var req logic.UserOrderListRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	orders, total, err := l.GetUserOrders(r.Context(), req.UserID, req.Page, req.PageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	var orderList []map[string]interface{}
	for _, order := range orders {
		var foods []map[string]interface{}
		for _, item := range order.OrderItems {
			foods = append(foods, map[string]interface{}{
				"food_id": item.FoodID,
				"name":    item.FoodName,
				"weight":  item.Weight,
				"price":   item.Price,
			})
		}

		orderList = append(orderList, map[string]interface{}{
			"order_id":    order.ID,
			"user_id":     order.UserID,
			"plate_id":    order.PlateID,
			"foods":       foods,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"created_at":  order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	httpx.OkJson(w, map[string]interface{}{
		"code":  0,
		"msg":   "success",
		"data":  orderList,
		"total": total,
	})
}

// GetOrderInfo 获取订单信息
func (h *RestaurantHandler) GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	if orderID == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("订单ID不能为空"))
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	order, err := l.GetOrderInfo(r.Context(), orderID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	var foods []map[string]interface{}
	for _, item := range order.OrderItems {
		foods = append(foods, map[string]interface{}{
			"food_id": item.FoodID,
			"name":    item.FoodName,
			"weight":  item.Weight,
			"price":   item.Price,
		})
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"order_id":    order.ID,
			"user_id":     order.UserID,
			"plate_id":    order.PlateID,
			"foods":       foods,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"created_at":  order.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// GetPlateDepot 获取餐盘托管处信息
func (h *RestaurantHandler) GetPlateDepot(w http.ResponseWriter, r *http.Request) {
	depotID := r.PathValue("depot_id")
	if depotID == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("托管处ID不能为空"))
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	depot, err := l.GetPlateDepot(r.Context(), depotID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"depot_id":  depot.ID,
			"name":      depot.Name,
			"location":  depot.Location,
			"capacity":  depot.Capacity,
			"available": depot.Available,
		},
	})
}

// HandleException 处理异常
func (h *RestaurantHandler) HandleException(w http.ResponseWriter, r *http.Request) {
	var req logic.WorkerExceptionRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	if err := l.HandleException(r.Context(), req.WorkerID, req.PlateID, req.Exception, req.Action); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "异常处理记录成功",
	})
}

// ProcessGC 处理GC
func (h *RestaurantHandler) ProcessGC(w http.ResponseWriter, r *http.Request) {
	var req logic.GCProcessRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewRestaurantLogic(h.svcCtx.DB)
	if err := l.ProcessGC(r.Context(), req.PlateID, req.Type); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJson(w, map[string]interface{}{
		"code": 0,
		"msg":  "GC处理成功",
	})
}
