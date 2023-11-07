package http_util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.in.codoon.com/Overseas/runbox/first-test/common"
	"git.in.codoon.com/Overseas/runbox/first-test/string_util"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	TOKEN_ID = "tokenId"
	UID      = "uid"
)

type CommentRsp struct {
	// in: body
	Rsp struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		Data   string `json:"data"`
		Desc   string `json:"desc"`
	}
}

func checkToken(ID string) (err error, userId string) {
	//ctx := context.WithValue(nil, "svc_name", "")
	token := common.VerifyIDTokenFireBase(context.Background(), ID)
	if token == nil || token.UID == "" {
		err = errors.New("token check failed")
		return
	}

	userId = token.UID
	return

}
func checkUser(c *gin.Context) error {
	token := c.Request.Header.Get(TOKEN_ID)
	user_id := c.Request.Header.Get(UID)
	err, uid := checkToken(token)
	if err == nil && uid != "" {
		if user_id != uid {
			log.Printf("runboxServer check uid failed ")
			return errors.New("user failed ; illegality uid")
		}
	} else {
		//token 校验失败
		log.Printf("runboxServer check token failed err: %v ", err)
		return errors.New("token failed")
	} //校验token结束
	return nil

}

func TakeToken() gin.HandlerFunc {
	//获取token：
	return func(c *gin.Context) {
		start := time.Now()
		token := c.Request.Header.Get(TOKEN_ID)
		user_id := c.Request.Header.Get(UID)

		log.Printf("runboxServer Started[token:%s] ", token)
		log.Printf("runboxServer Started[uid:%s] %s %s", user_id, c.Request.Method, c.Request.URL.Path)

		ck_err := checkUser(c)
		if ck_err != nil {
			//token校验失败
			rsp := CommentRsp{}.Rsp
			rsp.Status = FAILED
			rsp.Code = HTTP_CODE_AUTH_TOKEN_FAILED
			rsp.Desc = ck_err.Error()
			ReturnCompFunc(c, rsp)
		} else {
			//执行接口
			c.Next()
		}
		//name := fmt.Sprintf("HTTP %s %s%s", c.Request.Method, c.Request.Host, c.Request.URL.Path)
		log.Printf("runboxServer Completed[uid:%s] %s in %v", user_id, c.Request.URL.Path, time.Since(start))

	}
}

func ReqData2Form() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("MiddleWare") == "ON" || c.Request.Header.Get("WM") == "ON" {
			reqData2Form(c)
		}
	}
}

func reqData2Form(c *gin.Context) {
	userId := c.Request.Header.Get(UID)
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("read request body error:%v", err)
		return
	}
	// fmt.Printf("raw body:%s\n", data)
	var v map[string]interface{}
	if len(data) == 0 {
		v = make(map[string]interface{})
		err = nil
	} else {
		v, err = loadJson(bytes.NewReader(data))
	}
	if err != nil {
		// if request data is NOT json format, restore body
		// log.Printf("ReqData2Form parse as json failed. restore [%s] to body", string(data))
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
	} else {
		// if user_id in request is not empty, move it to req_user_id
		if uid, ok := v[UID]; ok {
			v["req_user_id"] = uid
		}
		// inject use_id into form
		v[UID] = userId
		form := map2Form(v)
		s := form.Encode()
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			c.Request.Body = ioutil.NopCloser(strings.NewReader(s))
		} else if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Request.Header.Del("Content-Type")
			// append url values
			urlValues := c.Request.URL.Query()
			for k, vv := range urlValues {
				if _, ok := form[k]; !ok {
					form[k] = vv
				}
			}
			c.Request.URL.RawQuery = form.Encode()
		} else {
			c.Request.Body = ioutil.NopCloser(strings.NewReader(s))
		}
	}
}

func map2Form(v map[string]interface{}) url.Values {
	form := url.Values{}
	var vStr string
	for key, value := range v {
		switch value.(type) {
		case string:
			vStr = value.(string)
		case float64, int, int64:
			vStr = fmt.Sprintf("%v", value)
		default:
			if b, err := json.Marshal(&value); err != nil {
				vStr = fmt.Sprintf("%v", value)
			} else {
				vStr = string(b)
			}
		}
		form.Set(key, vStr)
	}
	return form
}

func loadJson(r io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var v map[string]interface{}
	err := decoder.Decode(&v)
	if err != nil {
		// log.Printf("loadJson decode error:%v", err)
		return nil, err
	}
	return v, nil
}

func ReturnCompFunc(c *gin.Context, obj interface{}) {
	GinRsp(c, 200, obj)
}

func GinRsp(c *gin.Context, statusCode int, obj interface{}) {
	requestData := GetRequestData(c)
	objData := fmt.Sprintf("%+v", obj)

	clientIP := c.ClientIP()
	method := c.Request.Method
	userId := c.Request.Header.Get("uid")
	log.Printf("[GIN-RSP] [%s] [req_data: %s] [ip:%s]  [user_id:%s] [rsp:%s]",
		method,
		string_util.Cuts(requestData, 1024),
		clientIP,
		userId,
		string_util.Cuts(objData, 1024),
	)
	c.JSON(statusCode, obj)
}

func GetRequestData(c *gin.Context) string {
	var requestData string
	method := c.Request.Method
	if method == "GET" || method == "DELETE" {
		requestData = c.Request.RequestURI
	} else {
		c.Request.ParseForm()
		requestData = fmt.Sprintf("%s [%s]", c.Request.RequestURI, c.Request.Form.Encode())
	}
	return requestData
}
