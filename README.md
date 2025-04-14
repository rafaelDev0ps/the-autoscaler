# The Autoscaler 
#### Basically this is a lab to learn how an autoscaler works and how gRPC works (tbd)

## How it works
The client checks CPU and memory usage every 2 minutes, if the resource usage exceeds the threshold it request to the Orchestrator create a new instance. If the resources are not being used as expected the client request the Orchestrator to delete the instance.

## How to run this thing
> _Make sure you have Go and Docker is enabled in you machine!_   
  
Build the server image  
```sh
docker buildx build -f Dockerfile . -t server --build-arg ORCHESTRATOR_URL=$(ifconfig | grep -w inet | awk '{print $2} ' | tail -n 1)
```
  
Run the Orchestrator (server)
```sh
go run main.go
```