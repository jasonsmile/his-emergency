# API 约定

## 健康检查

`GET /health`，返回 HTTP 200 和 `{"status":"ok"}`，不查询数据库。

- 数据库查询或读取失败时返回 HTTP 500，响应包含 `code` 和 `message`。

前端 Demo 通过 Vite 将 `/api` 代理到本地 `http://127.0.0.1:8080`，启动方法见 [前端说明](../web/README.md)。

## 服务启动

服务启动时加载配置、初始化 MySQL 连接池并检查连接，然后初始化路由并监听端口；数据库初始化失败时退出。

## 收费模块

以下接口均位于 `/api` 下，且需要 Bearer Token。

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| POST | `/charge/getPendingCharges` | 查询挂号费和已签名处方的待收费项目，请求体传入 `encounter_id` |
| POST | `/charge/createCharge` | 按项目收费，支付方式为 `CASH`、`WECHAT`、`ALIPAY` 或 `CARD` |
| POST | `/charge/refundCharge` | 整笔或按 `refund_items` 指定数量退费 |
| GET | `/charge/getChargeRecords` | 查询收费记录，支持就诊、患者、状态、日期和分页筛选 |
| GET | `/charge/getDailyReport` | 查询指定日期及可选收费员的日结汇总 |
