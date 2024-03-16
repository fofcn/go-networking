package util

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUIDNoDash() string {
	// 生成UUID并去掉所有短横线
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
