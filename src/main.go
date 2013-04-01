/**
 * User: Leon
 * Date: 29/03/2013
 * Time: 22:25
 */
package main
import (
	"net/http"
	"fmt"
	"ballsite"
	"log"
	"html/template"
	"code.google.com/p/gorilla/sessions"
	"strings"
	"strconv"
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

	var ball = ballsite.RandomBall()
	return page.Execute(w, ball)
}

func ball(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	var page = template.Must(template.ParseFiles(
		"templates/_base.html",
		"templates/ball.html",
	))
	id_a := strings.Split(r.RequestURI, "/")
	id_s := id_a[len(id_a)-1]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		return err
	}

	fmt.Println("ball id:", id)
	var ball = ballsite.BallByID(id)
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
	var b = ballsite.AllBalls()
	return index.Execute(w, b)
}


func image(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	id_a := strings.Split(r.RequestURI, "/")
	id_s := id_a[len(id_a)-1]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		return err
	}

	path, err := ballsite.ImagePathById(id)
	if (err != nil) {
		return err;
	}
	fmt.Println("path:", path)
	http.ServeFile(w, r, path);
	return err
}

func thumb(w http.ResponseWriter, r *http.Request, ctx *Context) (err error) {
	id_a := strings.Split(r.RequestURI, "/")
	id_s := id_a[len(id_a)-1]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		return err
	}
	path, err := ballsite.ThumbPathById(id)
	if (err != nil) {
		return err;
	}
	fmt.Println("path:", path)
	http.ServeFile(w, r, path);
	return err
}

func main() {
	store = sessions.NewCookieStore([]byte("rape"))

	var c = ballsite.BallCount()
	fmt.Println("Ball count:", c)

	/*
	var ball = ballsite.Ball{}
	ball.Title = "First Ball"
	ball.Description = "Dies ist der erste Ball mit dem storm clouds and rainbow image"
	ball.ImagePath = "/var/www/polandball/media/images/storm_clouds_and_rainbow-wallpaper-1920x1200.jpg"
	ball.ThumbPath = "/var/www/polandball/media/thumbnails/storm_clouds_and_rainbow-wallpaper-1920x1200.jpg.png"

	fmt.Println("before insert:", ball)
	ball = ballsite.InsertBall(ball)
	fmt.Println("after insert:", ball)


	ball = ballsite.Ball{}
	ball.Title = "Second Ball"
	ball.Description = "Dies ist der zweite Ball mit dem hacker redux image"
	ball.ImagePath = "/var/www/polandball/media/images/hacker_redux_by_hashbox.png"
	ball.ThumbPath = "/var/www/polandball/media/thumbnails/hacker_redux_by_hashbox.png.png"

	fmt.Println("before insert:", ball)
	ball = ballsite.InsertBall(ball)
	fmt.Println("after insert:", ball)

	return
	*/


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
