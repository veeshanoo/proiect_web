package server

import (
	"net/http"
	"time"
)

const Key = "secret-cookie"
const duration = 360 * time.Second

func SetCookie(res http.ResponseWriter, value string) {
	expires := time.Now().Add(duration)
	cookie := &http.Cookie{
		Name:    Key,
		Value:   value,
		Path:    "/",
		Expires: expires,
	}
	http.SetCookie(res, cookie)
}

func UnsetCookie(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    Key,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(res, cookie)
}
