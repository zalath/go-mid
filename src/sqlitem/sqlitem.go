package sqlitem

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

/*Con ...*/
type Con struct {
	DB *sqlx.DB
}

//Opendb ...
func (c *Con) Opendb() {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8","root","root","localhost",3306,"go_mid")
	db, err := sqlx.Open("mysql", dsn)
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
	Req			string `db:"req" json:"req"`//发起访问的时间
}
//List a test
func (c *Con) List() []El {
	var err error
	var data = []El{}
	err = c.DB.Select(&data, fmt.Sprintf("select * from e where state=1 and (req=0 or req <%d) limit 0,400;",time.Now().Unix()))
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
	stmt, err := c.DB.Prepare("insert into e (url,txt,state,time,req) values(?,?,?,?,?)")
	if err != nil {
		stmt.Close()
		fmt.Println(err)
		c.haveErr(err)
		isdone = false
		return
	}
	res, er1 := stmt.Exec(el.Url, el.Txt, el.State, el.Time,el.Req)
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
	_, er1 := c.DB.Exec("delete from e where state=2")
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
			"id"  int primary key autoincrement,
			"url"  varchar(500) NOT NULL,
			"txt"  text NOT NULL,
			"state"  int(1) NOT NULL default 1,
			"time"  int(11) NOT NULL default 0,
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
