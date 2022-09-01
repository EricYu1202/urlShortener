package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/go-redis/redis/v9"
)

//const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const prefix = "http://localhost:8080/r?r="

var ctx = context.Background()

func ExampleClient() *redis.Client {
	fmt.Println("redis before connect!!!")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,
	})

	return rdb
}

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	fmt.Println(s)
	return s
}

func inputUrl(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	fmt.Println("method:", r.Method) //print request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/index.html")
		log.Println(t.Execute(w, nil))

	} else if r.Method == "POST" {
		//cover func

		//邏輯判斷
		err_parse_form := r.ParseForm()
		if err_parse_form != nil {
			fmt.Printf("parse error%s", err_parse_form)
		}

		fmt.Println("url:", template.HTMLEscapeString(r.Form.Get("url"))) //輸出到伺服器端

		outputUrl := "網址錯誤"

		//check url legal
		u, err := url.ParseRequestURI(r.Form.Get("url"))
		if err != nil || u.Scheme == "" || u.Host == "" {
			//panic(err)
			outputUrl = "網址格式錯誤"
		} else {

			for {
				//create random string
				outputUrl = randomString(5)

				val, err := client.Get(ctx, outputUrl).Result()

				if err == redis.Nil {
					fmt.Printf("%s does not exist\n", outputUrl)
					//insert value into redis
					err := client.Set(ctx, outputUrl, u, 0).Err()
					if err != nil {
						panic(err)
					}

					break
				} else if err != nil {
					outputUrl = "資料庫錯誤"
					panic(err)

				} else {
					//continue
					fmt.Printf("%s exists", val)
					continue
				}

			}

			//url concate
			outputUrl = fmt.Sprintf("%s%s", prefix, outputUrl)

		}

		//template.HTMLEscape(w, []byte(r.Form.Get("url"))) //輸出到客戶端
		t, _ := template.ParseFiles("template/result.html")
		log.Println(t.Execute(w, outputUrl))

	}
}

func redirect(w http.ResponseWriter, r *http.Request, client *redis.Client) {

	if err := r.ParseForm(); err != nil {
		// handle parse error
		http.Redirect(w, r, "/", 404)
	} else {
		shortUrl := r.Form.Get("r")
		fmt.Println("query : ", shortUrl)
		//select url from redis
		newUrl, err := client.Get(ctx, shortUrl).Result()
		fmt.Println("long url : ", newUrl)
		if err == redis.Nil {
			fmt.Println("404  url ")
			http.Redirect(w, r, "/", 404)

		} else if err != nil {

			//panic(err)
			t, _ := template.ParseFiles("template/result.html")
			log.Println(t.Execute(w, "資料庫錯誤"))

		} else {
			//redirect
			fmt.Println("long url : ", newUrl)
			http.Redirect(w, r, newUrl, 301)
		}

	}

}

func main() {
	fmt.Println("Hello World!")

	//rdb
	client := ExampleClient()

	//http
	//http.HandleFunc("/index", inputUrl)
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		inputUrl(w, r, client)
	})

	//http.HandleFunc("/r", redirect)
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		redirect(w, r, client)
	})

	//http.HandleFunc("/", inputUrl)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		inputUrl(w, r, client)
	})

	file_server := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset/", file_server))

	js_file_server := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", js_file_server))

	//http
	err := http.ListenAndServe(":8080", nil) //設定監聽的埠

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
