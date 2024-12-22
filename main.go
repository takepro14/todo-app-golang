package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID   int
	Name string
}

// Global variable for task list and generate id
var todos = struct {
	sync.Mutex
	Items  []Task
	NextID int
}{NextID: 1} // Initial ID

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
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
					<li id=task-{{.ID}}>{{.ID}}. {{.Name}}
						<form action="/delete" method="POST" style="display:inline">
							<input type="hidden" name="id" value="{{.ID}}">
							<button type="submit">Delete</button>
						</form>
					</li>
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
		taskName := r.FormValue("task")
		todos.Lock()
		todos.Items = append(todos.Items, Task{
			ID:   todos.NextID,
			Name: taskName,
		})
		todos.NextID++
		todos.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		todos.Lock()
		defer todos.Unlock()

		for i, task := range todos.Items {
			if task.ID == id {
				todos.Items = append(todos.Items[:i], todos.Items[i+1:]...)
				break
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
