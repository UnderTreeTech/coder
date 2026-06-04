package ecode

import (
	"user/internal/utils"

	"github.com/gin-gonic/gin"
)

type errors map[int]string

var ecodes = make(map[string]errors)

// Register register error msg
func Register(codes map[string]errors) {
	for locale, msgs := range codes {
		if existLocaleCodes, ok := ecodes[locale]; ok {
			for code, msg := range msgs {
				existLocaleCodes[code] = msg
			}
			ecodes[locale] = existLocaleCodes
		} else {
			ecodes[locale] = msgs
		}
	}
}

// Message return the errmsg associate with the errcode
func Message(ctx *gin.Context, code int) string {
	locale := utils.GetLocaleLng(ctx.Request.Header.Get("Accept-Language"))

	if localCodes, ok := ecodes[locale]; ok {
		return localCodes[code]
	}

	return ""
}

func init() {
	// 注册自定义错误码，相同错误码会覆盖全局，达到更换文案目的
	Register(map[string]errors{
		"en":    enRequestErrMsg,
		"zh-cn": zhRequestErrMsg,
	})

}
