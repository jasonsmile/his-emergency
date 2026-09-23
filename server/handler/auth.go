package handler

import (
	"database/sql"

	"emergency-his/server/errors"
	"emergency-his/server/logger"
	"emergency-his/server/middleware"
	"emergency-his/server/response"
	"github.com/gin-gonic/gin"
)

// Login 当前阶段只验证账号是否存在且启用；密码字段保留在协议中，待认证方案确定后补充校验。
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			LoginName string `json:"login_name" binding:"required"`
			Password  string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "登录名不能为空")
			return
		}
		var id, name, dept string
		query := `SELECT login_name, COALESCE(name, ''), COALESCE(dept_code, '') FROM sync_user WHERE login_name = ? AND status = 1`
		logger.SQL(query, req.LoginName)
		err := db.QueryRowContext(c.Request.Context(), query, req.LoginName).Scan(&id, &name, &dept)
		if err == sql.ErrNoRows {
			be := errors.NewBusinessError(40101, "登录名不存在或已停用")
			response.Error(c, 401, be.Code, be.Message)
			return
		}
		if err != nil {
			logger.Error.Printf("login database query failed login_name=%q error=%v", req.LoginName, err)
			response.InternalError(c, "登录失败")
			return
		}
		token, err := middleware.GenerateToken(id, name, dept)
		if err != nil {
			logger.Error.Printf("login token generation failed login_name=%q error=%v", req.LoginName, err)
			response.InternalError(c, "登录失败")
			return
		}
		response.Success(c, gin.H{"token": token, "user_id": id, "user_name": name, "dept_code": dept})
	}
}
