package handler

import (
	"net"

	"github.com/gin-gonic/gin"
)

// getClientIP 获取客户端真实IP，将IPv6回环地址转为IPv4
func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	// 将IPv6回环地址 ::1 转为 127.0.0.1
	if ip == "::1" {
		return "127.0.0.1"
	}
	// 将IPv6映射的IPv4地址（如 ::ffff:192.168.1.1）转为纯IPv4
	parsed := net.ParseIP(ip)
	if parsed != nil {
		if v4 := parsed.To4(); v4 != nil {
			return v4.String()
		}
	}
	return ip
}
