package http

import (
	"net/http"

	"user/internal/server/http/model"
	"user/internal/utils/reply"
	"github.com/gin-gonic/gin"
)

func getUserInfo(ctx *gin.Context) {
	req := &model.GetUserInfoReq{}
	if err := ShouldBind(ctx, &req); err != nil {
		return
	}

	resp, err := svc.GetUserInfo(ctx.Request.Context(), req.UserId)
	ctx.JSON(http.StatusOK, reply.Reply(ctx, resp, err))
}
