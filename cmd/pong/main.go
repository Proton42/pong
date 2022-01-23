package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	counter := 0

	go runPinger(&counter)
	//postToSlack()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]int{
			"Count": counter,
		}

		t.ExecuteTemplate(w, "index.html.tmpl", data)
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runPinger(cnt *int) {
	for {
		resp, httpErr := http.Get("https://expeditions.com")

		fmt.Println(resp)
		fmt.Println(httpErr)
		*cnt++
		time.Sleep(5 * time.Second)
	}
}

func postToSlack() {
	url := ""

	var jsonStr = []byte(`
		{"text":"Hello! My name is Pong. The friendly monitoring App"}
	`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
