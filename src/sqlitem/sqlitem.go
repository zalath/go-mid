package sqlitem

import (
	// "database/sql"
	"fmt"
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
	if err != nil {
		log.Fatal(err)
	}
	c.DB = db
}

//El ...
type El struct {
	ID      int    `db:"id" json:"id"`
	Url     string `db:"url" json:"url"`
	Txt     string    `db:"txt" json:"txt"`
	State		string `db:"state" json:"state"`//1新，2处理中
	Time		string `db:"time" json:"time"`
}
//List a test
func (c *Con) List() []El {
	// c.Opendb()
	var err error
	var data = []El{}
	err = c.DB.Select(&data, "select * from e where state=1;")
	// c.DB.Close()
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

//New create a new element
func (c *Con) New(el El) (isdone bool, newid int64) {
	isdone = true
	c.Opendb()
	stmt, err := c.DB.Prepare("insert into e (url,txt,state,time) values(?,?,?,?)")
	if err != nil {
		stmt.Close()
		fmt.Println(err)
		c.haveErr(err)
		isdone = false
		return
	}
	res, er1 := stmt.Exec(el.Url, el.Txt, el.State, el.Time)
	stmt.Close()
	if er1 != nil {
		c.haveErr(er1)
		c.DB.Close()
		isdone = false
		return
	}
	newid, _ = res.LastInsertId()
	c.DB.Close()
	return
}

//Del delete an element
func (c *Con) Del(id string) (isdone bool) {
	isdone = true
	// c.Opendb()
	_, er1 := c.DB.Exec("delete from e where id=$1", id)
	// c.DB.Close()
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}

//Update ...
func (c *Con) Update(id, val, col string) (isdone bool) {
	isdone = true
	// c.Opendb()
	var sb strings.Builder
	sb.WriteString("update e set ")
	sb.WriteString(col)
	sb.WriteString("=? where id=?")

	stmt, err := c.DB.Prepare(sb.String())
	if err != nil {
		stmt.Close()
		c.DB.Close()
		c.haveErr(err)
	}
	_, er1 := stmt.Exec(val, id)
	stmt.Close()
	// c.DB.Close()
	if er1 != nil {
		c.haveErr(er1)
		isdone = false
		return
	}
	return
}
func (c *Con) haveErr(err error) {
	if err.Error() == "no such table: e" {
		c.Opendb()
		sql := `CREATE TABLE "e" (
			"id"  INTEGER NOT NULL,
			"url"  TEXT NOT NULL,
			"txt"  TEXT NOT NULL,
			"state"  INTEGER NOT NULL default 1,
			"time"  TEXT,
			PRIMARY KEY ("id" ASC)
			);
			
			`
		_, err := c.DB.Exec(sql)
		c.DB.Close()
		if err != nil {
			log.Fatal("database error")
			return
		}
	}
}
