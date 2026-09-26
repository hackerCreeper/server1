package main

import (
	"html/template"
	"net/http"
)

type ServerHandler struct {
	terminal *Terminal
}

func NewServerHandler(t *Terminal) *ServerHandler {
	return &ServerHandler{terminal: t}
}

func (sh *ServerHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "HTML source template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (sh *ServerHandler) HandleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	inputCommand := r.FormValue("command")
	responseSignal := sh.terminal.Execute(inputCommand)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(responseSignal))
}
