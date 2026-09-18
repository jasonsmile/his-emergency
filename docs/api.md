# API 约定

## 健康检查

`GET /health`，返回 HTTP 200 和 `{"status":"ok"}`，不查询数据库。

## 查询患者

`GET /api/v1/patients?keyword=查询条件`

- `keyword` 可省略，前后空白会去除。
- 患者编号、手机号精确匹配；姓名使用 SQL LIKE 模糊匹配。
- 按 `sync_patient.id` 倒序，最多返回 100 条；当前不支持分页，不返回总记录数。
- 成功返回 `{"code":0,"message":"success","data":[...]}`；无匹配记录时 `data` 为 `[]`。
- 每条记录包含 `patientId`、`name`，以及可选的 `gender`、`birthDate`、`phone`、`cardNo`。出生日期格式为 `YYYY-MM-DD`。
- 数据库查询或读取失败时返回 HTTP 500，响应包含 `code` 和 `message`。

前端 Demo 通过 Vite 将 `/api` 代理到本地 `http://127.0.0.1:8080`，启动方法见 [前端说明](../web/README.md)。

## 服务启动

服务启动时加载配置、初始化 MySQL 连接池并检查连接，然后初始化路由并监听端口；数据库初始化失败时退出。
