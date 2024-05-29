package utils

import (
	"github.com/mitchellh/mapstructure"
	"github.com/punpeo/punpeo-lib/utils/jzcrypto"
	"net/url"
	"strings"
)

type SkToUser struct {
	UserId     uint   `mapstructure:"user_id"`
	Phone      string `mapstructure:"phone"`
	ExpireTime int64  `mapstructure:"expire_time"`
	Openid     string `mapstructure:"openid"`
	SessionKey string `mapstructure:"session_key"`
}

// ParseSk 将sk解析为 SkToUser 结构体 返回结构体指针.
func ParseSk(sk, key, iv string) *SkToUser {
	parse, err := jzcrypto.TripleDesDecrypt(sk, key, iv)
	if err != nil {
		return nil
	}
	parseArr := strings.Split(parse, "&")
	skMap := make(map[string]interface{})
	for _, value := range parseArr {
		item := strings.Split(value, "=")
		if len(item) != 2 {
			return nil
		}
		skMap[item[0]] = item[1]
	}
	var user SkToUser
	if err := mapstructure.WeakDecode(skMap, &user); err != nil {
		return nil
	}
	return &user
}

// ParseSkUrl 解析经过urlEncode后的sk.
// 实际开发中php项目中生成的sk会被urlEncode.
func ParseSkUrl(skUrlEncode, key, iv string) *SkToUser {
	if sk, err := url.QueryUnescape(skUrlEncode); err == nil {
		return ParseSk(sk, key, iv)
	}
	return nil
}
