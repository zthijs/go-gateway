package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

type Service struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
	Port   string `json:"port"`
}

func main() {

	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	r := mux.NewRouter()
	services := GetServices()

	for _, service := range services {

		target, err := url.Parse(fmt.Sprintf("http://localhost:%s/", service.Port))

		if err != nil {
			log.Printf("%s", red("Whoops something went wrong, try again"))
			os.Exit(0)
		}

		r.PathPrefix(service.Prefix).Handler(http.StripPrefix(service.Prefix, httputil.NewSingleHostReverseProxy(target)))

		log.Printf("Loaded %s", blue(service.Name))

	}

	r.Handle("/services.json", ServerInfo())
	r.Handle("/ping", Ping())

	var dir string

	flag.StringVar(&dir, "dir", "./cdn", "the directory to serve files from.")
	flag.Parse()

	r.PathPrefix("/cdn/").Handler(http.StripPrefix("/cdn/", http.FileServer(http.Dir(dir))))

	http.Handle("/", Logging(Headers(r)))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func GetServices() []Service {
	services := make([]Service, 3)
	raw, err := ioutil.ReadFile("./services.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &services)
	return services
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		yellow := color.New(color.FgYellow).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		start := time.Now()
		next.ServeHTTP(w, req)

		log.Printf("[ %s ] -> %s (%s)", yellow(req.Method), blue(req.RequestURI), green(time.Since(start)))
	})
}

func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, req)
	})
}

func ServerInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(string(fmt.Sprintf(`{"count" : %d}`, len(GetServices())))))
		return
	})
}

func Ping() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`pong`))
		return
	})
}

