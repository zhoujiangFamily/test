package common

import (
	"net/http"
)

func GetUserId(r *http.Request) string {
	user_id := r.Header.Get("uid")
	return user_id

}
