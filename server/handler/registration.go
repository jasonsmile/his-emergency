package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegistrationNotImplemented 保留挂号模块统一路由和响应契约。
// 尚未接入的功能明确返回 501，避免被误认为已经完成业务办理。
func RegistrationNotImplemented(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    50101,
			"message": feature + "功能待实现",
			"data":    nil,
		})
	}
}
