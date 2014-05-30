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
	wade.RegisterElement("t-userinfo", &UserInfo{
		Name: "Hai Thanh Nguyen",
		Age:  18,
	})
	wade.RegisterPageHandler("login_handler", func() interface{} {
		req := http.Service().NewRequest(http.MethodGet, "/auth")
		promise := wd.NewPromise(&AuthedStat{false}, req.DoAsync())
		promise.OnSuccess(func(r *http.Response) wd.ModelUpdater {
			return func(as *AuthedStat) {
				as.AuthGened = true
			}
		})
		return promise.Model()
	})

	http.Service().AddHttpInterceptor(func(req *http.Request) {
	})

	wade.Start()
}
