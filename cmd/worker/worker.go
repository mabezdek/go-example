package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{})

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	w := worker.New(c, "rubumo", worker.Options{})

	w.RegisterWorkflow(addEnvironment.AddEnvironmentWorkflow)
	w.RegisterActivity(addEnvironment.AddEnvironment)

	err = w.Run(worker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}
