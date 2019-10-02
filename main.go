package main

import (
	"flag"
	"net/http"

	"net/http/pprof"
	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gosagawa/isucon/controller"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func main() {
	flag.Set("bind", ":8080")
	rooter(goji.DefaultMux)
	goji.Serve()

	pprofMux := web.New()
	pprofMux.Get("/debug/pprof/", http.HandlerFunc(pprof.Index))
}

func rooter(m *web.Mux) http.Handler {

	m.Get("/user/index", controller.UserIndex)
	m.Get("/user/new", controller.UserNew)
	m.Post("/user/new", controller.UserCreate)
	m.Get("/user/edit/:id", controller.UserEdit)
	m.Post("/user/update/:id", controller.UserUpdate)
	m.Get("/user/delete/:id", controller.UserDelete)

	return m
}
