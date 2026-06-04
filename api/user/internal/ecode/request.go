package ecode

import "github.com/UnderTreeTech/waterdrop/pkg/status"

var (
	InternalError = status.New(10000001, "网络开小差，请稍后重试")
)

var enRequestErrMsg = errors{
	InternalError.Code(): "Network is unavailable, please have a try later",
}

var zhRequestErrMsg = errors{
	InternalError.Code(): InternalError.Message(),
}
