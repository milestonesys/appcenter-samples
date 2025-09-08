package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"apigateway-webserver/src/pkg/constants"
	"apigateway-webserver/src/pkg/entities/vms"
	handlers_context "apigateway-webserver/src/pkg/handlers/context"
	"apigateway-webserver/src/pkg/view"
)

type ViewHandler struct {
	mu sync.Mutex
}

func NewViewHandler() *ViewHandler {
	return &ViewHandler{}
}

func (vh *ViewHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vh.mu.Lock()
	defer vh.mu.Unlock()
	log.Println("ViewHandler.Handle() called")

	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	if username == "" {
		http.Error(w, "Missing required fields: username.", http.StatusBadRequest)
		return
	}

	path := "templates/view_events.html"
	tmpl, err := template.ParseFS(view.TemplateFS, path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parsing template file %s: %v", path, err), http.StatusInternalServerError)
		return
	}

	appCtx, exists := handlers_context.GetAppContextsInstance().GetAppContext(username)
	if !exists {
		http.Error(w, "App context not found.", http.StatusBadRequest)
		return
	}

	cameras, eventTypes, err := setupPageData(appCtx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not read data from the VMS: %v", err), http.StatusInternalServerError)
		return
	}

	camerasJson, err := cameras.ToJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Converting cameras to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	eventTypesJson, err := eventTypes.ToJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Converting event types to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	sessionJson, err := appCtx.GetWsCommandResponse().ToJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Converting session info to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// data to be passed to the template
	pageData := struct {
		AppName    string
		Cameras    string
		EventTypes string
		Username   string
		Session    string
	}{
		AppName:    constants.AppName,
		Cameras:    camerasJson,
		EventTypes: eventTypesJson,
		Username:   username,
		Session:    sessionJson,
	}
	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, fmt.Sprintf("Executing template: %v", err), http.StatusInternalServerError)
	}
}

func setupPageData(appCtx handlers_context.AppContext) (*vms.CamerasList, *vms.AnalyticEventTypes, error) {
	cameras, err := appCtx.GatewayService().RequestEnabledCameras(context.Background(), appCtx.Server(), appCtx.Token())
	if err != nil {
		return nil, nil, err
	}

	eventTypes, err := appCtx.GatewayService().RequestAnalyticEventTypes(context.Background(), appCtx.Server(), appCtx.Token())
	if err != nil {
		return nil, nil, err
	}
	return cameras, eventTypes, nil
}
