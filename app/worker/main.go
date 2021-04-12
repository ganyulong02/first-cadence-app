package main

import (
	"fmt"
	"github.com/ganyulong02/first-cadence-app/app/adapters/cadenceAdapter"
	"github.com/ganyulong02/first-cadence-app/app/config"
	"github.com/ganyulong02/first-cadence-app/app/worker/workflows"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func startWorkers(h *cadenceAdapter.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()
	fmt.Printf("Hostname: %v\n", appConfig.Cadence.HostPort)

	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	startWorkers(&cadenceClient, workflows.TaskListName)
	// The workers are supposed to be long running process that should not exit.
	select {}
}