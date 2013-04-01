/**
 * User: Leon
 * Date: 30/03/2013
 * Time: 09:04
 */
package balls

import (
	"database/sql"
	"log"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/kisielk/sqlstruct"
	"time"
	"fmt"
	"errors"
)

type Ball struct {
	Id int 					`sql:"id"`
	URL string 				`sql:"url"`
	Title string 			`sql:"title"`
	Description string 		`sql:"description"`
	Timestamp string 		`sql:"timestamp"`
	Date time.Time
	Public bool 			`sql:"public"`
	ImageURL string 		`sql:"imageurl"`
	ThumbURL string 		`sql:"thumburl"`
}

var _db *sql.DB

func db() (*sql.DB) {
	if (_db != nil) {
		return _db
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
	if (err != nil) {
		log.Fatal(err)
	}
	_db = db
	return _db
}

func close() {
	if (_db == nil) {
		return
	}
	err := _db.Close()
	if (err != nil) {
		log.Fatal(err)
	}
	_db = nil
}

func timeFromSQLTimestamp(ts string) (t time.Time, err error) {
	t, err = time.Parse("2006-01-02 15:04:05", ts)
	return
}

func unmarshalBalls(rows *sql.Rows) (balls []Ball, err error) {
	for rows.Next() {
		var ball Ball
		err := sqlstruct.Scan(&ball, rows)
		if (err != nil) {
			log.Fatal(err)
		}
		ball.Date, err = timeFromSQLTimestamp(ball.Timestamp)
		if (err != nil) {
			log.Fatal(err)
		}
		balls = append(balls, ball)
	}
	return
}

func All() []Ball {
	rows, err := db().Query(
	`SELECT b.id, b.url, m.title, m.description, m.timestamp, m.public, i.url as imageurl, i.thumbnailurl as thumburl
	FROM pb_balls.balls b, pb_balls.metadata m, pb_images.images i
	WHERE b.url = m.parenturl AND b.url = i.parenturl;`)
	if (err != nil) {
		log.Fatal(err)
	}

	balls, err := unmarshalBalls(rows)
	if (err != nil) {
		log.Fatal(err)
	}
	return balls
}

func ByURL(url string) Ball {
	fmt.Println("ball url", url)
	rows, err := db().Query(
	`SELECT b.id, b.url, m.title, m.description, m.timestamp, m.public, i.url as imageurl, i.thumbnailurl as thumburl
	FROM pb_balls.balls b, pb_balls.metadata m, pb_images.images i
	WHERE b.url = m.parenturl AND b.url = i.parenturl AND b.url = ?;`, url)
	if (err != nil) {
		log.Fatal(err)
	}

	arr, err := unmarshalBalls(rows)
	if (err != nil) {
		log.Fatal(err)
	}
	if (len(arr) == 0) {
		return Ball{}
	}
	return arr[0]
}

func Random() Ball {
	rows, err := db().Query("SELECT url FROM pb_balls.balls ORDER BY RAND() LIMIT 1;")
	if (err != nil) {
		log.Fatal(err)
	}
	var url string
	if (rows.Next()) {
		rows.Scan(&url)
	}
	return ByURL(url)
}

func Count() (cnt int64) {
	rows, err := db().Query("SELECT COUNT(*) AS cnt FROM pb_balls.balls;")
	if (err != nil) {
		log.Fatal(err)
	}
	if (rows.Next()) {
		rows.Scan(&cnt)
	}
	return
}

func ImagePathByURL(imageURL string) (path string, err error) {
	fmt.Println("imageURL:", imageURL)
	rows, err := db().Query("SELECT path FROM pb_images.images WHERE url LIKE '%"+imageURL+"%';")
	if (err != nil) {
		log.Fatal(err)
	}
	if (rows.Next()) {
		rows.Scan(&path)
	}
	if len(path) == 0 {
		err = errors.New("Fail!")
	}
	return
}

func ThumbPathByURL(thumbURL string) (path string, err error) {
	fmt.Println("thumbURL:", thumbURL)
	rows, err := db().Query("SELECT thumbnailpath FROM pb_images.images WHERE thumbnailurl LIKE '%"+thumbURL+"%';")
	if (err != nil) {
		log.Fatal(err)
	}
	if (rows.Next()) {
		rows.Scan(&path)
	}
	if len(path) == 0 {
		err = errors.New("Fail!")
	}
	return
}
