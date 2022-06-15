package el

import (
	"fmt"
	"strconv"
	"bytes"
	"io/ioutil"
	"gopkg.in/ini.v1"
	dbt "tasktask/src/sqlitem"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)



func Handle() {
	db := newdb()
	data := db.List()
	for i := 0; i < len(data); i++ {
		go send(data[i])
		db.Del(strconv.Itoa(data[i].ID))
	}
	defer db.DB.Close()
	return
}

func send(el dbt.El) bool {
	fmt.Println(el.Txt)
	req, err := http.NewRequest(http.MethodPost, el.Url, bytes.NewReader([]byte(el.Txt)))
	if err != nil {
		fmt.Println("post err", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("send do err", err1)
		return false
	}
	defer resp.Body.Close()
	// respBody, err2 := ioutil.ReadAll(resp.Body)
	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("read err", err2)
		return false
	}
	// fmt.Println("res:",string(respBody))
	return true
}

func Bit() {
	// send heart bit to main server
	cfg, _ := ini.Load("conf.ini")
	var b = dbt.El{}
	b.Url = cfg.Section("server").Key("url").String()
	b.Txt = "{\"bit\":\"1\"}"
	go send(b)
	// is := send(b)
	// if !is {
		// fmt.Println("bit send err")
	// }
}

//New create a new element into database
func New(c *gin.Context) (result string) {
	db := newdb()
	el := formEl(c, db)
	db.DB.Begin()
	res, newid := db.New(el)
	result = strconv.Itoa(int(newid))
	if !res {
		db.DB.MustBegin().Rollback()
		result = "mis"
	}
	db.DB.MustBegin().Commit()
	defer db.DB.Close()
	return
}

//Del delete el from db
func Del(id string) (result string) {
	db := newdb()
	db.DB.Begin()
	res := db.Del(id)
	if !res {
		db.DB.MustBegin().Rollback()
		result = "mis"
	} 
	db.DB.MustBegin().Commit()
	defer db.DB.Close()
	result = "done"
	return
}

func formEl(c *gin.Context, db *dbt.Con) dbt.El {
	var el = dbt.El{}
	el.Url = c.PostForm("url")
	el.Txt = c.PostForm("txt")
	el.State = "1"
	el.Time = time.Now().Format("2006-1-2 15:04:05")
	// fmt.Printf("%#v", el)
	return el
}

//GetEl ...
func GetEl(id string) dbt.El {
	db := newdb()
	res := db.Get(id)
	defer db.DB.Close()
	return res
}
//Save submit saving element
func Save(id, val, col string) string {
	db := newdb()
	res := db.Update(id, val, col)
	defer db.DB.Close()
	if res {
		return "done"
	}
	return "mis"
}

func Count() []dbt.Elc {
// func Count() int {
	db := newdb()
	res := db.Count()
	return res
}

func newdb() *dbt.Con {
	db := dbt.NewCon()
	return db
}
