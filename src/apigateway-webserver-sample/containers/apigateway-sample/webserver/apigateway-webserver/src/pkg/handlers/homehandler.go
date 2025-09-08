package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"apigateway-webserver/src/pkg/constants"
	"apigateway-webserver/src/pkg/constants/enums"
	"apigateway-webserver/src/pkg/view"
)

type HomeHandler struct {
	mu sync.Mutex
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (hh *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	hh.mu.Lock()
	defer hh.mu.Unlock()
	log.Println("HomeHandler.Handle() called")

	path := "templates/index.html"
	tmpl, err := template.ParseFS(view.TemplateFS, path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parsing template file %s: %v", path, err), http.StatusInternalServerError)
		return
	}

	pageData := struct {
		AppName              string
		CredentialsFlowTypes []string
	}{
		AppName:              constants.AppName,
		CredentialsFlowTypes: enums.GetCredentialsFlowTypes(),
	}
	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, fmt.Sprintf("Executing template: %v", err), http.StatusInternalServerError)
	}
}
