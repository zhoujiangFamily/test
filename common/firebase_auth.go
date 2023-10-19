package common

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"git.in.codoon.com/Overseas/runbox/first-test/conf"
	"git.in.codoon.com/Overseas/runbox/first-test/firebase_conf"
)

//验证token

func VerifyIDTokenFireBase(ctx context.Context, token string) *auth.Token {
	if token == "" {
		return nil
	}
	return firebase_conf.VerifyIDTokenAndCheckRevoked(ctx, conf.Firebaer_app, token)
}

//获取用户账号
func GetUserByUid(ctx context.Context, uid string) *auth.UserRecord {
	return firebase_conf.GetUser(ctx, conf.Firebaer_app, uid)
}
