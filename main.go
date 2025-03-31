package main

import (
	"fmt"
	"log"
	"net/http"
	"the-autoscaler/docker"
	"the-autoscaler/utils"
)

const (
	maxReplicas = 2
	minReplicas = 1
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func createInstanceHandler(w http.ResponseWriter, r *http.Request) {
	containerList, err := docker.GetAllContainerIDs()
	if err != nil {
		log.Fatal("Error listing all nodes: ", err)
	}

	if len(containerList) == maxReplicas {
		http.Error(w, "Maximum number of replicas reached", http.StatusNotAcceptable)
		return
	}

	newNodeName := fmt.Sprintf("node-%v", utils.RandomString(10))

	node, err := docker.CreateInstance(newNodeName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Node %s created!", node.ID)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Node %s created!", node.ID)))
}

func deleteInstanceHandler(w http.ResponseWriter, r *http.Request) {
	containerList, err := docker.GetAllContainerIDs()
	if err != nil {
		log.Fatal("Error listing all nodes: ", err)
	}

	if len(containerList) == minReplicas {
		http.Error(w, "Minimum number of replicas reached", http.StatusNotAcceptable)
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		http.Error(w, "Container ID is required", http.StatusBadRequest)
		return
	}

	if err := docker.DeleteInstance(containerID); err != nil {
		log.Fatal(err)
	}
	log.Printf("Node %s deleted!", containerID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Node %s deleted!", containerID)))
}

func main() {
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/create", createInstanceHandler)
	http.HandleFunc("/delete", deleteInstanceHandler)

	log.Println("Starting Orchestrator Server on :8080!")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
