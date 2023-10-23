package common

import (
	"encoding/json"
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
	var b Binding
	ctype := filterFlags(c.Header.Get("Content-Type"))
	switch {
	case c.Method == "GET" || c.Method == "DELETE" || ctype == MIMEPOSTForm || ctype == MIMEPOSTForm2B:
		b = Form
	case ctype == MIMEMultipartPOSTForm:
		b = MultipartForm
	case ctype == MIMEJSON:
		b = JSON
	case ctype == MIMEXML || ctype == MIMEXML2:
		b = XML
	default:
		return false
	}
	return BindWith(c, obj, b)
}
func BindWith(c *http.Request, obj interface{}, b Binding) bool {
	if err := b.Bind(c, obj); err != nil {
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
