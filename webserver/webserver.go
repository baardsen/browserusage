package webserver 

import (
	"html/template"
    "net/http"
    "fmt"
    "io"
    "strings"
    "encoding/json"
    "browserusage/dao"
    resourcelocator "github.com/baardsen/resourcelocator"
    "strconv"
    "time"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "BrowserUsage.html", nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	fmt.Fprint(w, jsonify(dao.Query(parseDate(r, "from"), parseDate(r, "to"))))
}

const dateFormat = "2006-01-02"
func parseDate(r *http.Request, param string) time.Time {
	date, _ := time.Parse(dateFormat, r.URL.Query().Get(param))
	return date
}

func jsonify(obj interface{}) string{
	json, err := json.Marshal(obj)
	checkError(err)
	return string(json)
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	resource := string(resourcelocator.Locate(r.URL.Path))

	if strings.HasPrefix(r.URL.Path, "/resources/script") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasPrefix(r.URL.Path, "/resources/stylesheets") {
		w.Header().Set("Content-Type", "text/css")
	}
	fmt.Fprintf(w, resource) 
} 

func makeHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request: %s\n", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(w, "%v", err)
			}
		}()
		handler(w, r)
	}
}

func renderTemplate(writer io.Writer, name string, data interface{}){
	tmpl, err := template.New(name).Parse(string(resourcelocator.Locate("/resources/templates/" + name)))
	checkError(err)
	tmpl.Execute(writer, data)
}

func checkError(err error){
	if err != nil {
	    panic(err)
	}
}

func Start(port int){
    http.HandleFunc("/", makeHandler(defaultHandler))
    http.HandleFunc("/data/", makeHandler(dataHandler))
    http.HandleFunc("/resources/", makeHandler(resourceHandler))
    fmt.Println("Listening on port " + strconv.Itoa(port))
    http.ListenAndServe(":"+strconv.Itoa(port), nil)
}