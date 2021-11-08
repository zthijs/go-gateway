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


// Typings of a service.
type Service struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
	Port   string `json:"port"`
}


// Main function
func main() {

	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Create a new router.
	r := mux.NewRouter()
	services := GetServices()


	// For each service create a HTTP handler
	for _, service := range services {

		target, err := url.Parse(fmt.Sprintf("http://localhost:%s/", service.Port))

		if err != nil {
			log.Printf("%s", red("Whoops something went wrong, try again"))
			os.Exit(0)
		}

		// Strip the path prefix.
		r.PathPrefix(service.Prefix).Handler(http.StripPrefix(service.Prefix, httputil.NewSingleHostReverseProxy(target)))

		log.Printf("Loaded %s", blue(service.Name))

	}


	// Handle default server routes.
	r.Handle("/services.json", ServerInfo())
	r.Handle("/ping", Ping())

	var dir string


	// Handle CDN requests.
	flag.StringVar(&dir, "dir", "./cdn", "the directory to serve files from.")
	flag.Parse()

	r.PathPrefix("/cdn/").Handler(http.StripPrefix("/cdn/", http.FileServer(http.Dir(dir))))

	// Handle logging.
	http.Handle("/", Logging(Headers(r)))


	// Start the server.
	log.Fatal(http.ListenAndServe(":8080", nil))

}


// Function to retrieve al te services.
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


// Middleware for logging.
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


// Handler to set global HTTP headers.
func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, req)
	})
}


// Handler to display server info.
func ServerInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(string(fmt.Sprintf(`{"count" : %d}`, len(GetServices())))))
		return
	})
}


// Handler for the ping request.
func Ping() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`pong`))
		return
	})
}

