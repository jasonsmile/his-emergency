// Package utils 提供服务端通用工具函数。
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateNo 生成单号。
// 格式为 prefix + YYYYMMDD + 4 位随机字符。
func GenerateNo(prefix string) string {
	return fmt.Sprintf("%s%s%s", prefix, time.Now().Format("20060102"), generateRandom(4))
}

// GenerateNoWithSeq 生成带日期和序列号的单号。
func GenerateNoWithSeq(prefix, dateStr string, seq int) string {
	return fmt.Sprintf("%s%s%04d", prefix, dateStr, seq)
}

func generateRandom(length int) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)[:length]
	}

	// crypto/rand 不可用时保留固定长度，避免单号生成流程崩溃。
	fallback := fmt.Sprintf("%0*x", length, time.Now().UnixNano())
	if len(fallback) < length {
		fallback = fmt.Sprintf("%0*s", length, fallback)
	}
	return fallback[:length]
}

// GetTodayDate 获取今天的日期，格式为 YYYY-MM-DD。
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// GetNowTime 获取当前本地时间。
func GetNowTime() time.Time {
	return time.Now()
}

// CalculateAge 根据出生日期计算年龄，返回如“18岁”“3月”“12天”的结果。
func CalculateAge(birthDate time.Time) string {
	now := time.Now()
	if birthDate.After(now) {
		return "0天"
	}

	years := now.Year() - birthDate.Year()
	anniversary := birthDate.AddDate(years, 0, 0)
	if anniversary.After(now) {
		years--
		anniversary = birthDate.AddDate(years, 0, 0)
	}
	if years > 0 {
		return fmt.Sprintf("%d岁", years)
	}

	months := 0
	monthiversary := birthDate
	for next := monthiversary.AddDate(0, 1, 0); !next.After(now); next = monthiversary.AddDate(0, 1, 0) {
		months++
		monthiversary = next
	}
	if months > 0 {
		return fmt.Sprintf("%d月", months)
	}

	days := int(now.Sub(birthDate).Hours() / 24)
	return fmt.Sprintf("%d天", days)
}
