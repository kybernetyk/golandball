/**
 * User: Leon
 * Date: 01/04/2013
 * Time: 21:17
 */
package ballsite

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
	Id          int64 		`sql:"id"`
	Title       string		`sql:"title"`
	Description string 		`sql:"description"`
	Timestamp   string 		`sql:"timestamp"`
	Date        time.Time
	Public      bool 		`sql:"public"`
	ImagePath   string 		`sql:"imagepath"`
	ThumbPath   string 		`sql:"thumbpath"`
	ImageId		int64 		`sql:"imageid"`
}

var _db *sql.DB

func db() (*sql.DB) {
	if (_db != nil) {
		return _db
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/polandball")
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

func AllBalls() []Ball {
	rows, err := db().Query(
	`SELECT b.id, m.title, m.description, m.timestamp, m.public,
	i.path as imagepath, i.thumbnailpath as thumbpath, i.id as imageid
	FROM balls b, metadata m, images i
	WHERE b.id = m.parentid AND b.id = i.parentid;`)
	if (err != nil) {
		log.Fatal(err)
	}

	balls, err := unmarshalBalls(rows)
	if (err != nil) {
		log.Fatal(err)
	}
	return balls
}

func BallByID(ballid int) Ball {
	fmt.Println("ball id:", ballid)
	rows, err := db().Query(
		`SELECT b.id, m.title, m.description, m.timestamp, m.public,
		i.path as imagepath, i.thumbnailpath as thumbpath, i.id as imageid
		FROM balls b, metadata m, images i
		WHERE b.id = m.parentid AND b.id = i.parentid AND b.id = ?;`, ballid)
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

func InsertBall(ball Ball) Ball {
	result, err := db().Exec("INSERT INTO balls () VALUES ();")
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	ball.Id = id

	result, err = db().Exec(`INSERT INTO metadata (parentid, title, description) VALUES (?, ?, ?);`,
		ball.Id, ball.Title, ball.Description)
	if err != nil {
		log.Fatal(err)
	}

	result, err = db().Exec(`INSERT INTO images (parentid, path, thumbnailpath) VALUES (?, ?, ?);`,
							ball.Id, ball.ImagePath, ball.ThumbPath)
	if err != nil {
		log.Fatal(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	ball.ImageId = id
	return ball
}

func RandomBall() Ball {
	rows, err := db().Query("SELECT id FROM balls ORDER BY RAND() LIMIT 1;")
	if (err != nil) {
		log.Fatal(err)
	}
	var id int
	if (rows.Next()) {
		rows.Scan(&id)
	}
	return BallByID(id)
}

func BallCount() (cnt int64) {
	rows, err := db().Query("SELECT COUNT(*) AS cnt FROM balls;")
	if (err != nil) {
		log.Fatal(err)
	}
	if (rows.Next()) {
		rows.Scan(&cnt)
	}
	return
}

func ImagePathById(imageId int) (path string, err error) {
	fmt.Println("imageId:", imageId)
	rows, err := db().Query("SELECT path FROM images WHERE id = ?;", imageId)
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

func ThumbPathById(imageId int) (path string, err error) {
	fmt.Println("thumbID:", imageId)
	rows, err := db().Query("SELECT thumbnailpath FROM images WHERE id = ?;", imageId)
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

