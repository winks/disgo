package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ungerik/go-gravatar"
	"html/template"
	"log"
	"time"
)

var m *martini.Martini

func main() {
	m = martini.New()
	m.Map(initDb())
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"http*"},
	}))
	m.Use(martini.Logger())
	m.Use(MapView)
	m.Use(render.Renderer(render.Options{
		Funcs: []template.FuncMap{
			{
				"formatTime": func(args ...interface{}) string {
					t1 := time.Unix(args[0].(int64), 0)
					return t1.Format(time.Stamp)
				},
				"gravatar": func(args ...interface{}) string {
					return gravatar.Url(args[0].(string))
				},
			},
		},
	}))
	m.Use(martini.Static("public"))
	r := martini.NewRouter()
	r.Get(`/comments/:id`, GetComment)
	r.Post(`/comments`, binding.Bind(Comment{}), CreateComment)
	r.Get(`/comments`, GetComments)
	r.Delete(`/comments/:id`, DestroyComment)
	m.Action(r.Handle)
	m.Run()
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "/tmp/test_db.bin")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Comment{}, "comments").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
