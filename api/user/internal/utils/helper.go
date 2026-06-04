package utils

import (
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"user/internal/i18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	defaultLocale = "zh-cn"
	_letters      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+-_"
)

// GetLocaleLng get locale language
func GetLocaleLng(lng string) (locale string) {
	// Multiple types, weighted with the quality value syntax:
	// Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5, en-US,en;q=0.5, en, *
	locale = strings.Split(lng, ";")[0]
	locale = strings.TrimSpace(locale)
	if strings.Contains(locale, ",") {
		elements := strings.Split(locale, ",")
		locale = elements[0]
		if strings.Contains(locale, "-") { // 英文系统一返回en
			if strings.Split(locale, "-")[0] == "en" {
				locale = "en"
			}
		}
	}

	// lng default to en if it doesn't have Accept-Language header or accept any language
	if "" == locale || "*" == locale || "zh" == locale {
		locale = defaultLocale
	}

	locale = strings.ToLower(locale)
	return
}

func TranslateError(ctx *gin.Context, err error) (errmsg string) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		lng := GetLocaleLng(ctx.Request.Header.Get("Accept-Language"))
		trans, ok := i18n.GetTranslator().Uni.GetTranslator(lng)
		if !ok {
			trans, _ = i18n.GetTranslator().Uni.GetTranslator("zh")
		}

		errmsg = errs[0].Translate(trans)
		if lngMap, ok := i18n.GetFieldMapping()[lng]; ok {
			if pathMap, ok := lngMap[ctx.Request.URL.Path]; ok {
				if fieldName, ok := pathMap[strings.ToLower(errs[0].Field())]; ok {
					errmsg = strings.Replace(errmsg, errs[0].Field(), fieldName, 1)
				}
			}
		}
	} else {
		errmsg = err.Error()
	}

	return
}

// RandomString generate random string by len
func RandomString(length int) string {
	sb := strings.Builder{}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		sb.WriteByte(_letters[rand.Intn(len(_letters))])
	}

	return sb.String()
}

// StripContentType strip content-type
// application/json;charset=utf-8
func StripContentType(contentType string) string {
	i := strings.Index(contentType, ";")
	if i != -1 {
		contentType = contentType[:i]
	}
	return contentType
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ParseDSNAddr 解析DSN中的IP+PORT，该方法仅适用于@(IP+PORT)/格式的DSN，如非此类型的请自定义解析方法
func ParseDSNAddr(dsn string) (addr string) {
	atIdx := strings.Index(dsn, "@")
	if atIdx == -1 {
		return
	}
	slot := strings.Trim(dsn[atIdx:], "@")
	addrs := strings.Split(slot, "/")
	if len(addrs) == 0 {
		return
	}
	addr = addrs[0]
	return
}
