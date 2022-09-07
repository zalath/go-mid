package el

import (
	"fmt"
	"strconv"
	"bytes"
	"gopkg.in/ini.v1"
	dbt "tasktask/src/sqlitem"
	dbs "tasktask/src/mysqls"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)



func Handle(con *dbt.Con) {
	data := con.List()
	for i := 0; i < len(data); i++ {
		// fmt.Println(data[i])
		go send(data[i])
		con.Del(strconv.Itoa(data[i].ID))
	}
	return
}

func send(el dbt.El) bool {
	req, err := http.NewRequest(http.MethodPost, el.Url, bytes.NewReader([]byte(el.Txt)))
	if err != nil {
		fmt.Println("post err", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("send do err", err1)
		return false
	}
	resp.Body.Close()
	return true
}

func Bit() {
	cfg, _ := ini.Load("conf.ini")
	var b = dbt.El{}
	b.Url = cfg.Section("server").Key("url").String()
	b.Txt = "{\"bit\":\"1\"}"
	go send(b)
}

//New create a new element into database
func New(c *gin.Context) (result string) {
	el := formEl(c)
	con := new(dbt.Con)
	res, newid := con.New(el)
	result = strconv.Itoa(int(newid))
	if !res {
		result = "mis"
	}
	return
}
// 用于mysql写入测试
// func T (i int, c *dbs.Con) (r string) {
func T (i int) (r string) {

	c := new(dbs.Con)
	c.Opendb()

	el := dbs.El{}
	// fmt.Println(i)
	el.Url = "http://test.sljj.com/sa/as"
	el.Code = strconv.Itoa(i)
	el.State = "1"
	el.Ctime =  time.Now().Format("2006-1-2 15:04:05")
	res, newid := c.New(i, el)
	r = strconv.Itoa(int(newid))
	c.DB.Close()
	if !res {
		r = "mis"
	}
	return
}

func formEl(c *gin.Context) dbt.El {
	var el = dbt.El{}
	el.Url = c.PostForm("url")
	el.Txt = c.PostForm("txt")
	el.State = "1"
	el.Time = time.Now().Format("2006-1-2 15:04:05")
	el.Req = c.PostForm("req")
	return el
}