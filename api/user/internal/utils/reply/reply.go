package reply

import (
	"strconv"

	"user/internal/ecode"

	"github.com/UnderTreeTech/waterdrop/pkg/trace"

	"github.com/gin-gonic/gin"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/UnderTreeTech/waterdrop/pkg/status"
)

// Response reply request result
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"msg"`
	TransId string      `json:"trans_id"`
	Datas   interface{} `json:"datas,omitempty"`
}

func Reply(ctx *gin.Context, data interface{}, err error) interface{} {
	estatus := status.ExtractStatus(err)
	reply := &Response{
		Code:    strconv.Itoa(estatus.Code()),
		Message: ecode.Message(ctx, estatus.Code()),
		TransId: trace.TraceID(ctx.Request.Context()),
		Datas:   data,
	}
	log.Debug(ctx.Request.Context(), "reply", log.Int("code", estatus.Code()), log.String("msg", reply.Message), log.Any("data", reply.Datas))
	return reply
}
