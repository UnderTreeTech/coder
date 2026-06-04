package model

type GetUserInfoReq struct {
	UserId string `form:"uid" json:"uid" validate:"required"`
}
