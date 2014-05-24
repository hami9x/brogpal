package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmcvetta/randutil"
	_ "github.com/lib/pq"
	"github.com/phaikawl/brogpal/model"
	"github.com/pilu/fresh/runner/runnerutils"
)

const (
	mySigningKey = "n0t9r34t6cz9r34tn0t1na9r34tw4y"
)

var (
	g *Environment = environment()
)

type Environment struct {
	devMode bool
	dbMap   *gorp.DbMap
}

func environment() *Environment {
	return &Environment{
		dbMap:   nil,
		devMode: true,
	}
}

func (g *Environment) IsDevMode() bool {
	return g.devMode
}

func (g *Environment) InitDb() {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", "postgres://postgres:thebest1om@localhost:5432/prognet")
	errPanic(err, "sql.Open failed")

	// construct a gorp DbMap
	g.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	g.dbMap.AddTableWithName(model.Post{}, "posts").SetKeys(true, "Id")
	g.dbMap.AddTableWithName(model.User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = g.dbMap.CreateTablesIfNotExists()
	errPanic(err, "Create tables failed")
}

func (g *Environment) Db() *gorp.DbMap {
	if g.dbMap == nil {
		panic("Trying to get an uninitialized dbmap")
	}
	return g.dbMap
}

func errPanic(err error, message string) {
	if err != nil {
		if message == "" {
			message = err.Error()
		}
		log.Printf(message)
		if g.IsDevMode() {
			panic(err.Error())
		}
	}
}

func checkErr(err error) {
	errPanic(err, "")
}

func makeRandomUserToken() (username string, tokenString string) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	username, err := randutil.AlphaStringRange(5, 10)
	errPanic(err, "Cannot random string, wtf?")
	token.Claims["username"] = username
	token.Claims["secret"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err = token.SignedString([]byte(mySigningKey))
	errPanic(err, "Cannot sign string, wtf?")
	return
}

func runnerMiddleware(w http.ResponseWriter, r *http.Request) {
	if runnerutils.HasErrors() {
		runnerutils.RenderError(w)
	}
}

func main() {
	m := martini.Classic()
	if os.Getenv("MARTINI_ENV") == "production" {
		g.devMode = false
	}

	if g.IsDevMode() {
		m.Use(runnerMiddleware)
	}
	g.InitDb()

	m.Use(martini.Static("../public", martini.StaticOptions{
		Prefix: "/public",
	}))
	m.Use(func(resp http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/api/") {
			token := req.Header.Get("AuthToken")
			if token != "" {
				n, err := g.Db().SelectInt(`select count(*) from users where token=$1`, token)
				if err != nil && err != sql.ErrNoRows {
					log.Fatalf(err.Error())
				}
				if n > 0 {
					return
				}
			}

			resp.WriteHeader(http.StatusUnauthorized)
			resp.Write([]byte("You're not allowed to do this, sorry."))
		} else if strings.HasPrefix(req.URL.Path, "/web/") {
			f, err := os.Open("../public/index.html")
			checkErr(err)
			conts, err := ioutil.ReadAll(f)
			checkErr(err)
			resp.Write(conts)
		}
	})

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
	})

	m.Get("/auth", func(r http.ResponseWriter) {
		username, token := makeRandomUserToken()
		user := &model.User{
			Username: username,
			Token:    token,
			Role:     model.RoleUser,
		}
		checkErr(g.Db().Insert(user))
		resp, err := json.Marshal(map[string]interface{}{
			"username": user.Username,
			"token":    user.Token,
		})
		checkErr(err)
		r.Write(resp)
	})

	m.Get("/api/test", func(r http.ResponseWriter) {
		r.Write([]byte("CLGT"))
	})

	m.Run()
}
