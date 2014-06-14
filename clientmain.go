package main

import (
	"github.com/phaikawl/brogpal/model"
	wd "github.com/phaikawl/wade"
	"github.com/phaikawl/wade/services/http"
	"github.com/phaikawl/wade/services/pdata"
)

type UserInfo struct {
	Name string
	Age  int
}

type AuthedStat struct {
	AuthGened bool
}

type UsernamePassword struct {
	Username string
	Password string
}

type RegUser struct {
	wd.Validated
	Data UsernamePassword

	Reset  func()
	Submit func()
}

func main() {
	//js.Global.Call("test", jquery.NewJQuery("title"))
	wade := wd.WadeUp("pg-home", "/web")
	wade.Pager().RegisterPages(map[string]string{
		"/home":          "pg-home",
		"/post":          "pg-post",
		"/post/new":      "pg-post-new",
		"/user":          "pg-user",
		"/user/login":    "pg-user-login",
		"/user/register": "pg-user-register",
		"/404":           "pg-not-found",
	})
	wade.Pager().SetNotFoundPage("pg-not-found")

	wade.RegisterNewTag("t-userinfo", UserInfo{})

	wade.RegisterNewTag("t-errorlist", wd.ErrorsBinding{})

	wade.Pager().RegisterHandler("pg-user-login", func() interface{} {
		req := http.Service().NewRequest(http.MethodGet, "/auth")
		promise := wd.NewPromise(&AuthedStat{false}, req.DoAsync())
		promise.OnSuccess(func(r *http.Response) wd.ModelUpdater {
			u := new(model.User)
			r.DecodeDataTo(u)
			pdata.Service().Set("authToken", u.Token)
			return func(au *AuthedStat) {
				au.AuthGened = true
			}
		})
		return promise.Model()
	})

	wade.Pager().RegisterHandler("pg-user-register", func() interface{} {
		ureg := new(RegUser)
		ureg.Validated.Init(ureg.Data)
		//ureg.Errors = map[string]map[string]interface{}{
		//	"Username": make(map[string]interface{}),
		//	"Password": make(map[string]interface{}),
		//}
		ureg.Reset = func() {
			ureg.Data.Password = ""
			ureg.Data.Username = ""
		}

		ureg.Submit = func() {
			//validate here...
			//then send
			wd.SendFormTo("/api/user/register", ureg.Data, &ureg.Validated).OnSuccess(
				func(r *http.Response) wd.ModelUpdater {
					//ureg.Errors = js.Global.Call("createObj")
					//println(ureg.Errors)
					//err := r.DecodeDataTo(ureg.Errors)
					//if err != nil {
					//	panic(err.Error())
					//}
					//println(ureg.errors["Password"])
					return nil
				})
		}
		return ureg
	})

	http.Service().AddHttpInterceptor(func(req *http.Request) {
		token, ok := pdata.Service().GetStr("authToken")
		if !ok {
			return
		}
		req.Headers.Set("AuthToken", token)
	})

	wade.Start()
}
