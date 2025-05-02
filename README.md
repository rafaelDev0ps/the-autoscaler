# The Autoscaler 
#### This is a lab project to learn how an autoscaler works and how gRPC works.
This project is in progress, check the [Future Improvements](#future-improvements) section to see the next steps to be done.  

## How gRPC Works
gRPC (gRPC Remote Procedure Call) is a high-performance, open-source universal RPC framework. It allows clients and servers to communicate seamlessly by defining services in a `.proto` file. The gRPC framework generates client and server code in multiple programming languages based on this definition.

In this project:
- The client sends system metrics (CPU and memory usage) to the server using the `CheckResourcesUsage` RPC method.
- The server evaluates the metrics and decides whether to scale up (create a new instance) or scale down (delete an instance) based on predefined thresholds.
- Communication between the client and server is handled over HTTP/2 using protocol buffers for serialization.

## How it works
The client checks CPU and memory usage every 2 minutes. If the resource usage exceeds the threshold, it requests the Orchestrator (server) to create a new instance. If the resources are underutilized, the client requests the Orchestrator to delete an instance.

## How to run this project
> _Make sure you have Go and Docker installed and enabled on your machine!_   

### Step 1: Build the server image
```sh
make build-server
```

### Step 2: Run the Orchestrator (server)
```sh
make go-run-server
# or 
make run-server
```

### Step 3: Build the client image
```sh
make build-client
```

### Step 4: Run the client
```sh
make go-run-client
# or
make run-client
```
  
This setup demonstrates how the autoscaler dynamically manages instances based on system resource usage.  

## Future Improvements
- Check ammount of connections to the client
- Check CPU trottling
- Create unit-tests and GH Workflows