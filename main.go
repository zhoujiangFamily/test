package main

import (
	"git.in.codoon.com/Overseas/runbox/first-test/conf"
	"git.in.codoon.com/Overseas/runbox/first-test/service"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	TOKEN_ID = "tokenId"
	UID      = "uid"

	//token 校验失败
	HTTP_CODE_AUTH_TOKEN_FAILED = 2001
	//UID 与TOKEN中不一致
	HTTP_CODE_AUTH_UID_FAILED = 2002
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("runboxServer Listening on port %s", port)

	/*serverName := "runboxServer"*/
	conf.InitBase()

	mux := http.NewServeMux()

	http.Handle("/", midHandler(http.HandlerFunc(service.Votes)))
	http.Handle("/gps", midHandler(http.HandlerFunc(service.Gps)))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func midHandler(next http.Handler) http.Handler {
	//获取token：

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		token := r.Header.Get(TOKEN_ID)
		user_id := r.Header.Get(UID)

		user_id = "abc"

		log.Printf("runboxServer Started[token:%s] ", token)
		log.Printf("runboxServer Started[uid:%s] %s %s", user_id, r.Method, r.URL.Path)

		//校验token开始

		err, uid := checkToken(token)
		if err == nil && uid != "" {
			if user_id != uid {
				log.Printf("runboxServer check uid failed ")
				w.WriteHeader(HTTP_CODE_AUTH_UID_FAILED)
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			//token 校验失败
			log.Printf("runboxServer check token failed ")
			w.WriteHeader(HTTP_CODE_AUTH_TOKEN_FAILED)
		} //校验token结束

		log.Printf("runboxServer Completed[uid:%s] %s in %v", user_id, r.URL.Path, time.Since(start))
		//之后
	})

	return nil
}
func checkToken(ID string) (err error, userId string) {
	/*	ctx := context.WithValue(nil, "svc_name", "")
		token := common.VerifyIDTokenFireBase(ctx, ID)
		if token == nil || token.UID == "" {
			err = errors.New("token check failed")
			return
		}

		userId = token.UID*/
	userId = "abc"
	return

}
