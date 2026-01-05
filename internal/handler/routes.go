package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"zeusro.com/gotemplate/internal/svc"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	handler := NewRestaurantHandler(serverCtx)

	// 健康检查
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/health",
				Handler: handler.HealthCheck,
			},
		},
	)

	// 用户相关
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/wallet/charge",
				Handler: handler.WalletCharge,
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/user/info/:user_id",
				Handler: handler.GetUserInfo,
			},
		},
	)

	// 餐盘相关
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/plate/bind",
				Handler: handler.BindPlate,
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/plate/unbind",
				Handler: handler.UnbindPlate,
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/plate/info/:plate_id",
				Handler: handler.GetPlateInfo,
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/plate/list",
				Handler: handler.GetPlateList,
			},
		},
	)

	// 订单相关
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/order/create",
				Handler: handler.CreateOrder,
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/order/list",
				Handler: handler.GetUserOrders,
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/order/info/:order_id",
				Handler: handler.GetOrderInfo,
			},
		},
	)

	// 餐盘托管处
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/depot/info/:depot_id",
				Handler: handler.GetPlateDepot,
			},
		},
	)

	// 工作人员
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/worker/exception",
				Handler: handler.HandleException,
			},
		},
	)

	// GC 处理
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/gc/process",
				Handler: handler.ProcessGC,
			},
		},
	)
}
