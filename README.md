# Emergency HIS

医院 HIS 扩展与应急保障平台，采用前后端分离的 B/S 架构建设。

项目定位是补充现有 C/S HIS 未覆盖或开发成本较高的日常业务功能，并在后续阶段提供 HIS 故障时限定范围内的应急业务能力。

> 当前阶段：第一期优先开发日常补充功能，暂不以完整应急接管作为验收目标。

## 建设目标

- 补充现有 HIS 未覆盖或迭代困难的业务功能。
- 通过浏览器访问，降低客户端安装和升级成本。
- 建立独立的权限、审计、接口和业务扩展能力。
- 后续在 HIS 故障时，接管经过验证的基础业务流程。
- 在 HIS 恢复后，支持应急数据核对、补录和回写。

## 技术栈

| 层次 | 技术 |
| --- | --- |
| 后端 | Go、Gin |
| 数据访问 | database/sql 或 sqlx、MySQL Driver |
| 前端 | Vue 3、TypeScript、Vite |
| 前端状态 | Pinia |
| 前端路由 | Vue Router |
| UI 组件 | Element Plus |
| HTTP 请求 | Axios |
| 数据库 | MySQL 8.4 |
| 反向代理 | Nginx |
| 部署 | Docker、Docker Compose |

Node.js 用于前端依赖安装、开发服务器和构建，Vue 3 构建后的静态文件由 Nginx 提供。

## 系统架构

```text
浏览器
   │
   ▼
Nginx
   ├── Vue 3 静态文件
   └── /api 反向代理
          │
          ▼
       Go Gin API
          │
          ▼
       MySQL 8.4
          │
          ▼
       HIS / 集成平台接口
```

## 推荐目录结构

```text
emergency-his/
├── README.md
├── backend/
├── frontend/
├── deploy/
├── database/
├── configs/
└── docs/
```

## 开发环境

- Go 1.24 或更高版本
- Node.js LTS 和 npm
- MySQL 8.4
- Docker Engine 和 Docker Compose Plugin
- Git

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

前端：

```bash
cd frontend
npm install
npm run dev
npm run build
```

## 配置与数据库

配置模板位于 `server/config/config.example.yaml`。部署时复制为 `server/config/config.yaml`，再填写服务器端口、MySQL 连接信息、JWT 密钥和 HIS 接口配置：

```bash
cp server/config/config.example.yaml server/config/config.yaml
```

`server/config/config.yaml` 已加入 `.gitignore`，不会提交到 Git。`.env.example` 用于 Docker 部署时的环境变量参考；两种配置方式最终应保持同一套参数，不能把真实密码写入仓库。

生产环境使用专用数据库账号，不使用 MySQL root 运行应用。数据库密码、JWT 密钥和 HIS 接口密钥不得提交到 Git，应通过受限配置文件或环境变量注入。

数据库脚本位于 `database/his_emergency_init.sql`，执行示例：

```bash
sudo mysql --default-character-set=utf8mb4 < database/his_emergency_init.sql
```

脚本负责创建数据库、33 张表、索引和默认配置；数据同步、HIS 回写、定时任务和应急切换由应用实现。

## Docker 与 Nginx

```bash
cd deploy
docker compose build
docker compose up -d
docker compose ps
docker compose logs -f backend
```

Nginx 对外提供统一入口：`/` 提供 Vue 静态文件，`/api/` 转发到 Go Gin，`/health` 用于健康检查。正式环境应配置 HTTPS、访问日志、请求大小和超时限制。

## 建设计划

第一期优先建设日常补充功能：

1. 访谈科室，确定一个高频且边界明确的业务痛点。
2. 完成患者、科室、医生等基础数据查询。
3. 建立登录、岗位权限、审计和健康检查。
4. 完成首个业务模块的前后端闭环。
5. 小范围试用并迭代。

后续建设应急模式、数据副本、应急流水、限定业务接管、恢复回写、重复防护和故障演练。

## 安全与运维

- 按医生、护士、收费员、药房和管理员划分权限。
- 对患者查询、处方、收费、退费、发药和配置变更保留审计记录。
- 应用账号只授予运行所需的数据库权限。
- 定期备份数据库并验证可恢复。
- 监控 HIS 接口失败、同步延迟和回写失败。
- 配置、日志和备份文件不得提交到公开仓库。

## 当前状态

数据库初始化脚本已准备；后端和前端工程尚未完成初始化。首个日常补充功能、HIS 接口、统一认证、部署网络和备份策略尚待确认。

## Git 初始化

```bash
cd emergency-his
git init
git add README.md
git commit -m "docs: add project README"
```
