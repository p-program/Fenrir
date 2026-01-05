# 智慧大学食堂 API 文档

基于 go-zero 框架实现的智慧食堂管理系统。

## 项目结构

```
Fenrir/
├── api/
│   └── restaurant.api          # API 定义文件
├── cmd/
│   └── restaurant/
│       └── main.go             # 服务入口
├── etc/
│   └── restaurant-api.yaml    # 配置文件
├── internal/
│   ├── config/
│   │   └── config.go          # 配置结构
│   ├── handler/
│   │   ├── restauranthandler.go  # 请求处理器
│   │   └── routes.go          # 路由注册
│   ├── logic/
│   │   ├── restaurant.go      # 业务逻辑
│   │   └── types.go           # 请求类型定义
│   └── svc/
│       └── servicecontext.go  # 服务上下文
└── model/
    ├── module.go              # 业务对象（保留原有设计）
    └── restaurant.go          # 数据模型（GORM）
```

## 核心功能

### 1. 用户管理
- 用户信息查询
- 钱包充值
- 钱包余额查询

### 2. 餐盘管理
- 餐盘绑定（用户与餐盘关联）
- 餐盘解绑（手动/自动）
- 餐盘信息查询
- 餐盘列表查询

### 3. 订单管理
- 创建订单（点餐）
- 订单查询
- 订单列表（分页）

### 4. 餐盘托管处
- 托管处信息查询
- 餐盘库存管理

### 5. 工作人员功能
- 异常处理记录
- 异常处理查询

### 6. GC 处理
- 餐盘清理
- 厨余垃圾处理

## API 接口

### 健康检查
```
GET /api/health
```

### 用户相关
```
POST /api/wallet/charge        # 钱包充值
GET  /api/user/info/:user_id   # 获取用户信息
```

### 餐盘相关
```
POST /api/plate/bind           # 绑定餐盘
POST /api/plate/unbind         # 解绑餐盘
GET  /api/plate/info/:plate_id # 获取餐盘信息
GET  /api/plate/list           # 获取餐盘列表（支持 ?is_bound=true/false 过滤）
```

### 订单相关
```
POST /api/order/create         # 创建订单
POST /api/order/list            # 获取用户订单列表
GET  /api/order/info/:order_id # 获取订单信息
```

### 餐盘托管处
```
GET /api/depot/info/:depot_id  # 获取托管处信息
```

### 工作人员
```
POST /api/worker/exception     # 处理异常
```

### GC 处理
```
POST /api/gc/process           # GC 处理（餐盘清理/厨余垃圾处理）
```

## 配置说明

配置文件：`etc/restaurant-api.yaml`

```yaml
Name: restaurant-api
Host: 0.0.0.0
Port: 8888

Database:
  Type: sqlite        # 支持 sqlite, mysql, postgres
  DSN: restaurant.db  # 数据库连接字符串

Log:
  ServiceName: restaurant-api
  Mode: file
  Path: logs
  Level: info
```

## 运行项目

### 1. 安装依赖
```bash
go mod tidy
```

### 2. 运行服务
```bash
go run cmd/restaurant/main.go -f etc/restaurant-api.yaml
```

### 3. 测试接口
```bash
# 健康检查
curl http://localhost:8888/api/health

# 钱包充值
curl -X POST http://localhost:8888/api/wallet/charge \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user123","amount":100.0}'
```

## 数据库模型

### 核心表结构
- `users` - 用户表
- `wallets` - 钱包表
- `transactions` - 交易记录表
- `plates` - 餐盘表
- `foods` - 食物表
- `orders` - 订单表
- `order_items` - 订单明细表
- `plate_depots` - 餐盘托管处表
- `workers` - 工作人员表
- `exception_logs` - 异常处理记录表
- `gc_process_logs` - GC处理记录表

## 业务逻辑说明

### 餐盘绑定流程
1. 用户扫描餐盘二维码或 RFID
2. 系统检查餐盘是否可用
3. 如果用户已有绑定餐盘，自动解绑旧餐盘
4. 绑定新餐盘，更新状态为 `in_use`

### 点餐流程
1. 用户选择食物和重量
2. 系统计算订单总价
3. 检查用户钱包余额
4. 创建订单并扣款
5. 记录交易记录
6. 更新订单状态为 `paid`

### 自动解绑机制
- 用户未使用餐盘超过 15-20 分钟，系统自动解绑
- 可通过定时任务实现

## 开发说明

### 添加新接口
1. 在 `api/restaurant.api` 中定义 API
2. 在 `internal/logic/restaurant.go` 中实现业务逻辑
3. 在 `internal/handler/restauranthandler.go` 中添加处理器
4. 在 `internal/handler/routes.go` 中注册路由

### 数据库迁移
项目启动时会自动执行数据库迁移，创建所有必要的表结构。

## 注意事项

1. 生产环境建议使用 MySQL 或 PostgreSQL
2. 需要实现自动解绑的定时任务
3. 建议添加 JWT 认证中间件
4. 建议添加请求限流和熔断保护
5. 建议添加日志记录和监控

