package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

//Config service configuration
type Config struct {
	SlackWebhookUrl string `envconfig:"SLACK_WEBHOOK_URL" required:true"`
	Port            string `envconfig:"PORT" default:"8080"`
}

func main() {
	conf := Config{}
	if cErr := envconfig.Process("", &conf); cErr != nil {
		fmt.Println(cErr)
		panic("failed to load service config")
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

	log.Println("listening on", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
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
	url := "https://hooks.slack.com/services/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"

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
