# 系统架构

当前采用 `Nginx → Vue 3 → Go Gin → MySQL` 的前后端分离结构。后端按 `router → api → service → repository → entity` 分层，医疗业务模块不得在 Handler 中直接访问数据库。
