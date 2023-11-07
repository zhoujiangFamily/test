package controller

import (
	"git.in.codoon.com/Overseas/runbox/first-test/common"
	"git.in.codoon.com/Overseas/runbox/first-test/http_util"
	"github.com/gin-gonic/gin"
	"log"
)

type PostGpsReq struct {
	UserId string `json:"user_id"`
}

type PostGpsRsp struct {
	// in: body
	Rsp struct {
		Code   int            `json:"code"`
		Status string         `json:"status"`
		Data   PostGpsRspData `json:"data"`
		Desc   string         `json:"description"`
	}
}
type PostGpsRspData struct {
	//路线ID
	RouteId string `json:"route_id"`
	//奖章
	Medals interface{} `json:"medals"`
	//现在等级
	Grade int `json:"grade"`
	//老等级
	OldGrade int `json:"old_grade"`
}

func PostGps(c *gin.Context) {

	userId := common.GetUserId(c.Request)

	gpsDto := PostGpsReq{
		UserId: userId,
	}
	rsp := PostGpsRsp{}.Rsp
	rsp.Status = http_util.FAILED
	rsp.Code = http_util.HTTP_CODE_SUCCESS
	//绑定参数
	if http_util.Bind(c.Request, gpsDto) {
		log.Printf("bing request failed ")
		rsp.Desc = "bing request param failed"
		rsp.Code = http_util.HTTP_CODE_PARAM_FAILE
		http_util.ReturnCompFunc(c, rsp)
		return
	}

	out := PostGpsRspData{}
	//todo

	rsp.Data = out
	rsp.Status = http_util.SUCCESS
	http_util.ReturnCompFunc(c, rsp)
	return

}
