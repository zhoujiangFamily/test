package http_util

import (
	"encoding/json"
	"git.in.codoon.com/Overseas/runbox/first-test/common"
	"log"
	"math"
	"net/http"
)

const (
	AbortIndex            = math.MaxInt8 / 2
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEPOSTForm2B        = "application/x-www-form-urlencode" // be compatible with codoon Android. WTF!
	MIMEMultipartPOSTForm = "multipart/form-data"
)

const (
	HTTP_CODE_SUCCESS = 200

	HTTP_CODE_PARAM_FAILE    = 601
	HTTP_CODE_BUSINESS_FAILE = 602

	//token 校验失败
	HTTP_CODE_AUTH_TOKEN_FAILED = 701
	//UID 与TOKEN中不一致
	HTTP_CODE_AUTH_UID_FAILED = 702
)

const (
	FAILED  = "failed"
	SUCCESS = "success"
)

func RenderJson(w http.ResponseWriter, code int, data ...interface{}) error {
	writeHeader(w, code, "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data[0])
}
func Render(w http.ResponseWriter, code int, data ...interface{}) error {
	err := RenderJson(w, code, data)
	if err != nil {
		log.Printf("Render failed parse form: %v ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	return nil
}

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
}

func Bind(c *http.Request, obj interface{}) bool {
	var b common.Binding
	ctype := filterFlags(c.Header.Get("Content-Type"))
	switch {
	case c.Method == "GET" || c.Method == "DELETE" || ctype == MIMEPOSTForm || ctype == MIMEPOSTForm2B:
		b = common.Form
	case ctype == MIMEMultipartPOSTForm:
		b = common.MultipartForm
	case ctype == MIMEJSON:
		b = common.JSON
	case ctype == MIMEXML || ctype == MIMEXML2:
		b = common.XML
	default:
		log.Printf("Render failed parse form:  ")
		return false
	}
	//b = JSON
	return BindWith(c, obj, b)
}
func BindWith(c *http.Request, obj interface{}, b common.Binding) bool {
	if err := b.Bind(c, obj); err != nil {
		log.Printf("xaxsa %v", err)
		return false
	}
	return true
}
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

type CommonRsp struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Desc   string      `json:"description"`
	Code   int         `json:"code"`
}
