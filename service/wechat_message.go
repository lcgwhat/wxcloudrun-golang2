package service

import "net/http"

func WechatMessageHandle(w http.ResponseWriter, r *http.Request) {
	r.PostForm.Get("")
}
