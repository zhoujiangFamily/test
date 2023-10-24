package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.in.codoon.com/Overseas/runbox/first-test/common"
	"git.in.codoon.com/Overseas/runbox/first-test/http_util"
	"git.in.codoon.com/Overseas/runbox/first-test/service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	TOKEN_ID = "tokenId"
	UID      = "uid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("runboxServer Listening on port %s ", port)

	//serverName := "runboxServer"
	//conf.InitBase()

	router := http.NewServeMux()

	router.Handle("/note", http.HandlerFunc(service.Votes))
	router.Handle("/v1/gps", midHandler(http.HandlerFunc(service.Gps)))
	router.Handle("/v1/test", midHandler(http.HandlerFunc(service.Test)))

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func midHandler(next http.Handler) http.Handler {
	//获取token：

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		token := r.Header.Get(TOKEN_ID)
		user_id := r.Header.Get(UID)

		log.Printf("runboxServer Started[token:%s] ", token)
		log.Printf("runboxServer Started[uid:%s] %s %s", user_id, r.Method, r.URL.Path)

		//校验token开始
		//checkUser(w, r)
		/*	if e1 != nil {
			return
		}*/
		//处理 http 相关
		delHttpBody(w, r, user_id)

		//之后
		next.ServeHTTP(w, r)
		log.Printf("runboxServer Completed[uid:%s] %s in %v", user_id, r.URL.Path, time.Since(start))
	})
	return nil
}
func checkUser(w http.ResponseWriter, r *http.Request) error {
	rsp := http_util.CommonRsp{
		Status: http_util.SUCCESS,
		Code:   http_util.CODE_SUCCESS,
		Desc:   "",
	}
	token := r.Header.Get(TOKEN_ID)
	user_id := r.Header.Get(UID)
	err, uid := checkToken(token)
	if err == nil && uid != "" {
		if user_id != uid {
			log.Printf("runboxServer check uid failed ")
			w.WriteHeader(http_util.HTTP_CODE_AUTH_TOKEN_FAILED)
			http.Error(w, "check uid failed ", http_util.HTTP_CODE_AUTH_UID_FAILED)
			rsp.Desc = "runboxServer check uid failed"
			rsp.Status = http_util.FAILED
			http_util.Render(w, 200, rsp)
			return errors.New("user failed")
		}
	} else {
		//token 校验失败
		log.Printf("runboxServer check token failed err: %v ", err)
		w.WriteHeader(http_util.HTTP_CODE_AUTH_TOKEN_FAILED)
		rsp.Desc = "check token failed"
		rsp.Status = http_util.FAILED
		http_util.Render(w, 200, rsp)
		return errors.New("token failed")
	} //校验token结束
	return nil

}

func delHttpBody(w http.ResponseWriter, r *http.Request, user_id string) {
	//处理http
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body error:%v", err)
		return
	}
	log.Printf(" user: %s , body : %s ", user_id, string(data))

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
		r.Body = ioutil.NopCloser(bytes.NewReader(data))
	} else {
		v["user_id"] = user_id
		form := map2Form(v)
		s := form.Encode()
		if r.Method == "POST" || r.Method == "PUT" {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Body = ioutil.NopCloser(strings.NewReader(s))
		} else if r.Method == "GET" || r.Method == "DELETE" {
			r.Header.Del("Content-Type")
			urlValues := r.URL.Query()
			for k, vv := range urlValues {
				if _, ok := form[k]; !ok {
					form[k] = vv
				}
			}
			r.URL.RawQuery = form.Encode()

		} else {
			r.Body = ioutil.NopCloser(strings.NewReader(s))

		}
	}
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
