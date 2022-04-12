package workflows

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func Connect() client.Client {
	client, err := client.NewClient(client.Options{})

	if err != nil {
		log.Fatal(err)
	}

	return client
}

var Client client.Client = Connect()

func ExecuteWorkflow(w interface{}, args ...interface{}) (client.WorkflowRun, error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "rubumo",
	}

	return Client.ExecuteWorkflow(context.Background(), workflowOptions, w, args...)
}
