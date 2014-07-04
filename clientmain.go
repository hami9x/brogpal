package main

import (
	"reflect"

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
	Validated
	Data UsernamePassword
}

func (r *RegUser) Reset() {
	r.Data.Password = ""
	r.Data.Username = ""
}

func (r *RegUser) Submit() {
	//validate here...
	//then send
	wd.SendFormTo("/api/user/register", r.Data, &r.Errors)
}

type PostView struct {
	PostId int
}

type ErrorListModel struct {
	Errors map[string]interface{}
}

type ErrorMap map[string]map[string]interface{}

type Validated struct {
	Errors ErrorMap
}

func (v *Validated) Init(dataModel interface{}) {
	m := make(ErrorMap)
	typ := reflect.TypeOf(dataModel)
	if typ.Kind() != reflect.Struct {
		panic("Validated data model passed to Init() must be a struct.")
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		m[f.Name] = make(map[string]interface{})
	}
	v.Errors = m
}

func main() {
	//js.Global.Call("test", jquery.NewJQuery("title"))
	wade := wd.WadeUp("pg-home", "/web", "wade-content", "wpage-container", func(wade *wd.Wade) {
		wade.Pager().RegisterPages(map[string]string{
			"/home":          "pg-home",
			"/posts":         "pg-post",
			"/posts/new":     "pg-post-new",
			"/post/:postid":  "pg-post-view",
			"/user":          "pg-user",
			"/user/login":    "pg-user-login",
			"/user/register": "pg-user-register",
			"/404":           "pg-not-found",
		})
		wade.Pager().SetNotFoundPage("pg-not-found")

		//wade.Custags().RegisterNew("t-userinfo", UserInfo{})
		wade.Custags().RegisterNew("t-errorlist", ErrorListModel{})
		wade.Custags().RegisterNew("t-test", UsernamePassword{})

		wade.Pager().RegisterController("pg-user-login", func(p *wd.PageData) interface{} {
			req := http.Service().NewRequest(http.MethodGet, "/auth")
			as := &AuthedStat{false}
			req.DoAsync().Done(func(r *http.Response) {
				u := new(model.User)
				r.DecodeDataTo(u)
				pdata.Service().Set("authToken", u.Token)
				as.AuthGened = true
			})
			return as
		})

		wade.Pager().RegisterController("pg-user-register", func(p *wd.PageData) interface{} {
			ureg := new(RegUser)
			ureg.Validated.Init(ureg.Data)
			return ureg
		})

		wade.Pager().RegisterController("pg-post-view", func(p *wd.PageData) interface{} {
			pv := new(PostView)
			p.ExportParam("postid", &pv.PostId)
			return pv
		})
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
