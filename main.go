package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	time2 "time"
)

const db_type = "mysql"
const db_name = "test"
const db_user = "root"
const db_pass = "root"

func getDb() (db *sql.DB) {
	db, err := sql.Open(db_type, db_user+":"+db_pass+"@/"+db_name+"?charset=utf8")
	checkerr(err)
	return db
}

type User struct {
	Name  string
	Infos string
}

type Title struct {
	Head string
	Type string
}

type Index struct {
	P    User
	Info Title
}

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数
	cookie, err := r.Cookie("Cookie")
	var username string
	if err == nil {
		db := getDb()
		db.QueryRow("select username from userinfo where cookie=?", cookie.Value).Scan(&username)
	}
	p := User{Name: username, Infos: "logout"}
	if username == "" {
		p.Name = "Anonymous"
		p.Infos = "login"
		http.SetCookie(w, &http.Cookie{Name: "Cookie", Path: "/", MaxAge: -1})
	}
	m := Title{Head: "ssss", Type: "ssssss"}
	s := Index{P: p, Info: m}
	t, _ := template.ParseFiles("/html/index.html")
	log.Println(t.Execute(w, s))
}

func logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//db, err :=sql.Open("mysql","root:root@/test?charset=utf8")
	//checkerr(err)
	db := getDb()
	cookie, err := r.Cookie("Cookie")
	var cookies string
	if err != nil {
		db.QueryRow("select cookie from userinfo where email=?", cookie).Scan(&cookie)
		if cookies != "" {
			mts, err := db.Prepare("update userinfo set cookie=\"\" where cookie=?")
			checkerr(err)
			mts.Exec(template.HTMLEscapeString(cookie.Value))
		}
		fmt.Fprintf(w, "<script>alert(\"logouted\");window.location.href=\"/login\"</script>")
	} else {
		mts, err := db.Prepare("update userinfo set cookie=\"\" where cookie=?")
		mts.Exec(template.HTMLEscapeString(cookie.Value))
		http.SetCookie(w, &http.Cookie{Name: "Cookie", Path: "/", MaxAge: -1})
		checkerr(err)
		defer mts.Close()
		fmt.Fprintf(w, "<script>alert(\"logout success\");window.location.href=\"/index\"</script>")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//db, err :=sql.Open("mysql","root:root@/test?charset=utf8")
	db := getDb()
	log.Println("method:", r.Method)
	var cookie string
	fmt.Println(cookie)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./Html/login.html")
		log.Println(t.Execute(w, nil))
	} else {
		email := template.HTMLEscapeString(r.Form["email"][0])
		password := template.HTMLEscapeString(r.Form["password"][0])
		db.QueryRow("select cookie from userinfo where email=?", email).Scan(&cookie)
		if cookie == "" {
			var passwords string = ""
			err := db.QueryRow("SELECT password FROM userinfo where email=?", email).Scan(&passwords)
			if err != nil {
				fmt.Println(err)
			}
			if passwords == password && password != "" {
				//t,_:= template.ParseFiles("./html/index.html")
				//t.Execute(w,nil)
				http.SetCookie(w, &http.Cookie{
					Name:  "Cookie",
					Value: GetRandomSalt(),
				})
				stmt, err := db.Prepare("update userinfo set cookie=? where email=?")
				checkerr(err)
				stmt.Exec(GetRandomSalt(), email)
				fmt.Fprintf(w, "<script>alert(\"login success!\");window.location.href=\"/\"</script>")
				//http.Redirect(w,r,"/",http.StatusFound)
			} else {
				fmt.Fprintf(w, "<script>alert(\"login faild\");window.location.href=\"/login\"</script>")
			}
		} else {
			cookies, err := r.Cookie("Cookie")
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:  "Cookie",
					Value: cookie,
				})
			}
			fmt.Println(cookies)
			fmt.Fprintf(w, "<script>alert(\"logined\");window.location.href=\"/\"</script>")
		}
	}
}

func reg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("../html/register.html")
		log.Println(t.Execute(w, nil))
	} else {
		//db, err :=sql.Open("mysql","root:root@/test?charset=utf8")
		//if err != nil{
		//	fmt.Println(err)
		//}
		db := getDb()
		stmt, err := db.Prepare("INSERT userinfo SET username=?,password=?,email=?,time=?")
		checkerr(err)
		username := template.HTMLEscapeString(r.Form["username"][0])
		password := template.HTMLEscapeString(r.Form["password"][0])
		email := template.HTMLEscapeString(r.Form["email"][0])
		time := time2.Now().Format("2006-01-02 15:04:05")
		log.Println(time)
		var flag bool = true
		che, err := db.Query("select email from userinfo")
		defer che.Close()
		for che.Next() {
			var emaild string
			che.Scan(&emaild)
			if emaild == email {
				flag = false
				break
			}
		}
		if flag {
			res, err := stmt.Exec(username, password, email, time)
			checkerr(err)
			id, err := res.LastInsertId()
			checkerr(err)
			fmt.Println(id)
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			fmt.Fprintf(w, "<script>alert(\"faild,same email\");window.location.href=\"/reg\"</script>")
		}
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(32)
}

//生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time2.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", reg)
	http.HandleFunc("/logout", logout)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Listen: ", err)
	}
}
