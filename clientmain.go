package main

import (
	wd "github.com/phaikawl/wade"
	"github.com/phaikawl/wade/services/http"
)

type UserInfo struct {
	Name string
	Age  int
}

type AuthedStat struct {
	AuthGened bool
}

func main() {
	wade := wd.WadeUp("_wade")
	wade.RegisterPages(map[string]string{
		"/home":       "pg-home",
		"/post":       "pg-post",
		"/post/new":   "pg-post-new",
		"/user":       "pg-user",
		"/user/login": "pg-user-login",
		"/404":        "pg-not-found",
	})
	wade.SetNotFoundPage("pg-not-found")
	wade.RegisterElement("t-userinfo", &UserInfo{
		Name: "Hai Thanh Nguyen",
		Age:  18,
	})
	wade.RegisterPageHandler("login_handler", func() interface{} {
		req := http.Service().NewRequest(http.MethodGet, "/auth")
		promise := wd.NewPromise(&AuthedStat{false}, req.DoAsync())
		promise.OnSuccess(func(r *http.Response) wd.ModelUpdater {
			return func(au *AuthedStat) {
				au.AuthGened = true
			}
		})
		return promise.Model()
	})

	http.Service().AddHttpInterceptor(func(req *http.Request) {
	})

	wade.Start()
}
