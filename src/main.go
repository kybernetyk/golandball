/**
 * User: Leon
 * Date: 29/03/2013
 * Time: 22:25
 */
package main
import (
	"net/http"
	"fmt"
	"balls"
	"log"
	"html/template"
	"code.google.com/p/gorilla/sessions"
)

type Context struct {
	Session *sessions.Session
}

var store sessions.Store

func NewContext(req *http.Request) (*Context, error) {
	sess, err := store.Get(req, "gostbook")
	return &Context{
		//Database: session.Clone().DB(database),
		Session: sess,
	}, err
}

type rapeHandler func(http.ResponseWriter, *http.Request, *Context) error

func makeHandler(hf rapeHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := NewContext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = hf(w, r, ctx)
		if (err != nil) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func random(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	var page = template.Must(template.ParseFiles(
		"templates/_base.html",
		"templates/ball.html",
	))

	var ball = balls.Random()
	return page.Execute(w, ball)
}

func ball(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	var page = template.Must(template.ParseFiles(
		"templates/_base.html",
		"templates/ball.html",
	))
	var url = "http://localhost:8080" + r.RequestURI

	fmt.Println("ball:", url)
	var ball = balls.ByURL(url)
	return page.Execute(w, ball)
}


func index(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	var index = template.Must(template.ParseFiles(
		"templates/_base.html",
		"templates/index.html",
	))
//	var s = ctx.Session
//
//	flashes := s.Flashes()
//	for _, v := range flashes {
//		fmt.Println(v)
//	}

//	fmt.Println("ficken:", s.Values["ficken"])
//	var x = s.Values["ficken"].(int)
//	x++
//	s.Values["ficken"] = x
//	s.AddFlash("Hurrenson!")
//
//	s.Save(r,w)
	var b = balls.All()
	return index.Execute(w, b)
}


func image(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	var url = r.RequestURI
	path, err := balls.ImagePathByURL(url)
	if (err != nil) {
		return err;
	}
	fmt.Println("path:", path)
	http.ServeFile(w, r, path);
	return err
}

func thumb(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	var url = r.RequestURI
	path, err := balls.ThumbPathByURL(url)
	if (err != nil) {
		return err;
	}
	fmt.Println("path:", path)
	http.ServeFile(w, r, path);
	return err
}

func main() {
	store = sessions.NewCookieStore([]byte("rape"))

	var c = balls.Count()
	fmt.Println("Ball count:", c)

	http.Handle("/", makeHandler(index))
	http.Handle("/ball/", makeHandler(ball))
	http.Handle("/rand", makeHandler(random))
	http.Handle("/view/image/", makeHandler(image))
	http.Handle("/view/thumbnail/", makeHandler(thumb))
	err := http.ListenAndServe(":8080", nil)
	if (err != nil) {
		log.Fatal(err)
	}
}
