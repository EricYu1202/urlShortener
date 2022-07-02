package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const prefix = "http://localhost:8080/r?r="

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	fmt.Println(s)
	return s
}

func inputUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //print request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/index.html")
		log.Println(t.Execute(w, nil))

	} else if r.Method == "POST" {
		//邏輯判斷
		r.ParseForm()

		fmt.Println("url:", template.HTMLEscapeString(r.Form.Get("url"))) //輸出到伺服器端

		//check url legal

		//check url exits in redis or not

		//create random string
		outputUrl := randomString(5)
		//check if it is unique in redis

		//url concate
		outputUrl = fmt.Sprintf("%s%s", prefix, outputUrl)
		//template.HTMLEscape(w, []byte(r.Form.Get("url"))) //輸出到客戶端
		t, _ := template.ParseFiles("template/result.html")
		log.Println(t.Execute(w, outputUrl))

	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		// handle parse error
		http.Redirect(w, r, "/", 404)
	} else {
		shortUrl := r.Form.Get("r")
		fmt.Println("short url : ", shortUrl)
		//select url from redis
		newUrl := "/"

		//redirect
		http.Redirect(w, r, newUrl, 301)
	}

}

func main() {
	fmt.Println("Hello World!")

	//http
	http.HandleFunc("/index", inputUrl)
	http.HandleFunc("/", inputUrl)
	http.HandleFunc("/r", redirect)

	fileServer := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fileServer))

	js_fileServer := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", js_fileServer))

	err := http.ListenAndServe(":8080", nil) //設定監聽的埠

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
