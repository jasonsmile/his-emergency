# 患者查询 Demo

基于 Vue 3、TypeScript、Element Plus，参考 Pure Admin 的浅色后台布局。页面使用真实的 `GET /api/v1/patients` 接口，不内置模拟患者。

## 本地启动

在当前 `web` 目录执行：

```powershell
npm install
npm run dev
```

打开终端显示的本地地址。默认将 `/api` 代理到 `http://127.0.0.1:8080`。
如果后端位于另一台机器，将 `.env.example` 复制为 `.env`，修改 `API_PROXY_TARGET` 后重启开发服务。

## 查询行为

- 页面打开时自动查询；支持回车查询、重置和刷新。
- 姓名模糊匹配，患者编号或手机号精确匹配；不支持卡号搜索。
- 空条件返回按同步表 ID 倒序排列的最多 100 条记录；这不代表总患者数，也不代表就诊时间顺序。
- 当前接口没有分页参数和总数，页面不展示虚假的分页。
- “查看详情”展示列表记录的基础字段，不请求额外接口。
- 查询失败会清除旧结果，显示错误和重试按钮。

## 代码位置

- `src/App.vue`：侧栏、顶栏和页面外框。
- `src/views/PatientList.vue`：查询表单、表格、详情抽屉与状态处理。
- `src/api/patients.ts`：请求封装和接口类型。
- `src/styles.css`：全局布局和页面样式。
- `vite.config.ts`：本地 API 代理。

## 构建

```powershell
npm run build
```

生成 `dist` 静态文件。Vite 开发代理不适用于生产环境；正式部署需由 Nginx 等将 `/api` 转发到 Go 服务。当前仓库的 Compose 只启动后端，此 Demo 未修改生产部署配置。
