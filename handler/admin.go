package handler

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/pascalj/disgo/models"
	"net/http"
	"strconv"
)

// AdminIndex shows the overview of the admin interface with latest comments.
func AdminIndex(req *http.Request, ren render.Render, dbmap *gorp.DbMap) {
	qry := req.URL.Query()
	page, err := strconv.Atoi(qry.Get("page"))

	if err != nil {
		page = 0
	}
	comments := paginatedComments(dbmap, page)
	ren.HTML(200, "admin/index", comments, render.HTMLOptions{
		Layout: "admin/layout",
	})
}

// UnapprovedComments will only list unapproved comments, else it behaves like AdminIndex.
func UnapprovedComments(ren render.Render, dbmap *gorp.DbMap, cfg models.Config) {
	count, err := dbmap.SelectInt("select count(*) from comments where approved<>1")
	if err == nil && count > 0 {
		var comments []models.Comment
		dbmap.Select(&comments, "select * from comments where approved<>1 order by created desc")
		ren.HTML(200, "admin/unapproved", comments, render.HTMLOptions{
			Layout: "admin/layout",
		})
	} else {
		ren.Redirect(cfg.General.Prefix + "/admin")
	}

}

// GetRegister shows the register form.
func GetRegister(ren render.Render) {
	ren.HTML(200, "admin/register", nil, render.HTMLOptions{
		Layout: "admin/layout",
	})
}

// PostUser will create a new user when no other users are in the database.
// If users are present, it will redirect to the login.
func PostUser(ren render.Render, req *http.Request, dbmap *gorp.DbMap, cfg models.Config) {
	if models.UserCount(dbmap) == 0 {
		email, password := req.FormValue("email"), req.FormValue("password")
		user := models.NewUser(email, password)
		err := dbmap.Insert(&user)
		if err != nil {
			ren.Redirect(cfg.General.Prefix + "/register")
		} else {
			ren.Redirect(cfg.General.Prefix + "/login")
		}
	}
}

// RegireLogin is a middleware that ensures that only an admin call the following
// handler(s).
func RequireLogin(rw http.ResponseWriter, req *http.Request,
	s sessions.Session, dbmap *gorp.DbMap, c martini.Context, cfg models.Config) {
	obj, err := dbmap.Get(models.User{}, s.Get("userId"))

	if err != nil || obj == nil {
		http.Redirect(rw, req, cfg.General.Prefix+"/login", http.StatusFound)
		return
	}

	user := obj.(*models.User)
	c.Map(user)
}

// GetLogin shows the login form for the backend.
func GetLogin(ren render.Render, dbmap *gorp.DbMap, cfg models.Config) {
	if models.UserCount(dbmap) > 0 {
		ren.HTML(200, "admin/login", nil, render.HTMLOptions{
			Layout: "layout",
		})
	} else {
		ren.Redirect(cfg.General.Prefix + "/register")
	}

}

// PostLogin takes the email and password parameter and logs the user in if they are correct.
func PostLogin(ren render.Render,
	req *http.Request,
	session sessions.Session,
	dbmap *gorp.DbMap,
	cfg models.Config) {
	var user models.User

	email, password := req.FormValue("email"), req.FormValue("password")
	err := dbmap.SelectOne(&user, "select * from users where email = :email", map[string]interface{}{"email": email})
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		ren.Redirect(cfg.General.Prefix + "/login")
		return
	}

	session.Set("userId", user.Id)
	ren.Redirect(cfg.General.Prefix + "/admin/")
}

// PostLogout logs the user out and redirects to the login page.
func PostLogout(ren render.Render, session sessions.Session, cfg models.Config) {
	session.Clear()
	ren.Redirect(cfg.General.Prefix + "/login")
}

func paginatedComments(dbmap *gorp.DbMap, page int) *models.PaginatedComments {
	var comments []models.Comment
	pages, err := dbmap.SelectInt("select ceil(count(*)/10) from comments")
	if err != nil {
		pages = 1
	}
	dbmap.Select(&comments, "select * from comments order by created desc limit 10 offset ?", page*10)
	return &models.PaginatedComments{int(pages), page, 10, comments}
}
