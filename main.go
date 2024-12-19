package main

import (
	"html/template"
	"net/http"
	"sync"
)

var todos = struct {
	sync.Mutex
	Items []string
}{}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(`
        <html>
        <head><title>TODO App</title></head>
        <body>
            <h1>TODO App</h1>
            <ul>
                {{range .}}
                    <li>{{.}}</li>
                {{end}}
            </ul>
            <form action="/add" method="POST">
                <input type="text" name="task" required>
                <button type="submit">Add</button>
            </form>
        </body>
        </html>
        `))
	todos.Lock()
	defer todos.Unlock()
	tmpl.Execute(w, todos.Items)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		todos.Lock()
		todos.Items = append(todos.Items, task)
		todos.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
