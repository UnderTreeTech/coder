package sanitizer

import (
	"reflect"
	"unicode"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xslice"
	"github.com/gin-gonic/gin"

	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	p *bluemonday.Policy
}

var skipUrls = []string{}

func NewSanitizer() *Sanitizer {
	sanitizer := &Sanitizer{
		p: bluemonday.UGCPolicy(),
	}

	//todo: customize allow elements
	//sanitize.p.AllowElements("style")

	return sanitizer
}

func (s *Sanitizer) Sanitize(ctx *gin.Context, req interface{}) {
	if xslice.ContainString(skipUrls, ctx.Request.URL.Path) {
		return
	}
	s.sanitize(ctx, reflect.ValueOf(req))
}

// sanitize user input
func (s *Sanitizer) sanitize(ctx *gin.Context, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		s.sanitize(ctx, v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			s.sanitize(ctx, v.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			s.sanitize(ctx, v.Index(i))
		}
	case reflect.String:
		// html.EscapeString is lighter, bluemonday is more comprehensive
		// Both sanitize will replace &
		//v.SetString(html.EscapeString(v.String()))
		if !v.CanSet() {
			return
		}

		if len(v.String()) == 0 {
			return
		}

		if isNumeric(v.String()) {
			return
		}

		v.SetString(s.p.Sanitize(v.String()))
	}
}

// isNumeric Checks if the string contains only digits. A decimal point is not a digit and returns false.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
