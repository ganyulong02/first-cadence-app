package main

/**

With this simple http server, we can now easily start and signal the workflow using 2 endpoints! To start the workflow,
use your favourite REST client and make a post request to http://localhost:3030/api/start-hello-world.
The response should be something like this
{
    "ID": "8e1e251a-6603-46f9-b0bf-05c0fcb8baf1",
    "RunID": "c2d68c3b-4786-4148-a240-2e467f57c9d2"
}
Now check the Cadence web UI to see if the workflow has successfully started! If so we can now signal the workflow
by making a post request to http://localhost:3030/api/signal-hello-world?workflowId=<wid>&age=25

 */

import (
	"context"
	"encoding/json"
	"github.com/ganyulong02/first-cadence-app/app/adapters/cadenceAdapter"
	"github.com/ganyulong02/first-cadence-app/app/config"
	"github.com/ganyulong02/first-cadence-app/app/worker/workflows"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	cadenceAdapter *cadenceAdapter.CadenceAdapter
	logger *zap.Logger
}

func (h *Service) triggerHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountId := r.URL.Query().Get("accountId")

		wo := client.StartWorkflowOptions{
			TaskList: workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.Workflow, accountId)
		if err != nil {
			http.Error(w, "Error starting workflow!", http.StatusBadRequest)
			return
		}
		h.logger.Info("Started workflow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) signalHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST"{
		workflowId := r.URL.Query().Get("workflowId")
		age, err := strconv.Atoi(r.URL.Query().Get("age"))
		if err != nil {
			h.logger.Error("Failed to parse age from request!")
		}
		err = h.cadenceAdapter.CadenceClient.SignalWorkflow(context.Background(), workflowId, "", workflows.SignalName, age)
		if err != nil {
			http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
			return
		}
		h.logger.Info("Signaled workflow iwth the following params!", zap.String("WorkflowId", workflowId), zap.Int("Age", age))
		js, _ := json.Marshal("Success")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func main() {
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	service := Service{&cadenceClient, appConfig.Logger}
	http.HandleFunc("/api/start-hello-world", service.triggerHelloWorld)
	http.HandleFunc("/api/signal-hello-world", service.signalHelloWorld)

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}


