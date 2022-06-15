package sqlitem

import (
	// "database/sql"
	// "fmt"
	"log"
	"strings"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //sqlite3
)

/*Con ...*/
type Con struct {
	DB *sqlx.DB
}

//Opendb ...
func (c *Con) Opendb() {
	db, err := sqlx.Connect("sqlite3", "./db.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	c.DB = db
}

//NewCon ...
func NewCon() *Con {
	var c = new(Con)
	c.Opendb()
	return c
}

//El ...
type El struct {
	ID      int    `db:"id" json:"id"`
	Url     string `db:"url" json:"url"`
	Txt     string    `db:"txt" json:"txt"`
	State		string `db:"state" json:"state"`//1新，2处理中
	Time		string `db:"time" json:"time"`
}
type Elc struct {
	Count int `db:"count" json:"count"`
}
func (c *Con) Count() []Elc {
	db := c.DB
	var data = []Elc{}
	err := db.Select(&data, "select count(*) as count from e where state=1;")
	if err != nil {
		c.haveErr(err)
		return data
	}
	return data
}
//List a test
func (c *Con) List() []El {
	db := c.DB
	var err error
	var data = []El{}
	err = db.Select(&data, "select * from e where state=1;")
	if err != nil {
		c.haveErr(err)
		return nil
	}
	for i := 0; i < len(data); i++ {
		// 更新数据表数据状态
		_ = c.Update(strconv.Itoa(data[i].ID), "2", "state")
	}
	return data
}

//Get select online from database on id
func (c *Con) Get(id string) El {
	db := c.DB
	el := El{}
	err := db.Get(&el, "select * from e where id = ?", id)
	if err != nil {
		c.haveErr(err)
	}
	return el
}

//New create a new element
func (c *Con) New(el El) (isdone bool, newid int64) {
	isdone = true
	db := c.DB
	stmt, err := db.Prepare("insert into e (url,txt,state,time) values(?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		c.haveErr(err)
		isdone = false
		return
	}
	res, er1 := stmt.Exec(el.Url, el.Txt, el.State, el.Time)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	newid, _ = res.LastInsertId()
	return
}

//Del delete an element
func (c *Con) Del(id string) (isdone bool) {
	isdone = true
	db := c.DB
	_, er1 := db.Exec("delete from e where id=$1", id)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}

//Update ...
func (c *Con) Update(id, val, col string) (isdone bool) {
	// fmt.Println("update el :", id, col, val)
	isdone = true
	db := c.DB

	var sb strings.Builder
	sb.WriteString("update e set ")
	sb.WriteString(col)
	sb.WriteString("=? where id=?")

	stmt, err := db.Prepare(sb.String())
	defer stmt.Close()
	if err != nil {
		c.haveErr(err)
	}
	_, er1 := stmt.Exec(val, id)
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}
func (c *Con) haveErr(err error) {
	if err.Error() == "no such table: e" {
		db := c.DB
		sql := `CREATE TABLE "e" (
			"id"  INTEGER NOT NULL,
			"url"  TEXT NOT NULL,
			"txt"  TEXT NOT NULL,
			"state"  INTEGER NOT NULL default 1,
			"time"  TEXT,
			PRIMARY KEY ("id" ASC)
			);
			
			`
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("database error")
			return
		}
	} else {
		// fmt.Println(err)
	}
}
