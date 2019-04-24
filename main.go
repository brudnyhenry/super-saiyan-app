package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var bootTime = time.Now().Format("2006-01-02 15:04:05")
var tpl *template.Template

const gitURL string = "https://api.github.com/repos/"

type data struct {
	OS       string
	Hostname string
	Time     string
	Rating   string
	Repo     string
}

func init() {
	tpl = template.Must(template.ParseFiles("template/index_new.gohtml"))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	h, _ := os.Hostname()
	d := data{
		OS:       runtime.GOOS,
		Hostname: h,
		Time:     bootTime,
	}
	project, ok := r.Form["project"]
	if ok {
		p := strings.TrimPrefix(project[0], "https://github.com/")
		rating := generateRating(getStars(p))
		d.Rating = rating
		d.Repo = p
	}

	err = tpl.ExecuteTemplate(w, "index_new.gohtml", d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateRating(r int) string {
	if r < 10 {
		return "weak"
	} else if r >= 100 && r < 5000 {
		return "strong"
	} else {
		return "legendary"
	}
}

func getStars(p string) int {
	var result map[string]interface{}
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}

	resp, err := client.Get(fmt.Sprintf("%s%s", gitURL, p))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &result)
		return int(result["stargazers_count"].(float64))
	}
	return 0
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", RootHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
