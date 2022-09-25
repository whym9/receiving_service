# PCAP processing service

Consists of services: receiving_service, pcap_statistics, saving_service

Step 1. Cloning repositories

Clone files in your preferred local repository by entering this commands in CMD:

```
git clone https://github.com/whym9/receiving_service.git

git clone https://github.com/whym9/pcap_statistics.git

git clone https://github.com/whym9/saving_service.git

```

Step 2. Building Docker Images

Start building docker images subsequentially. Go to each repository's Dockerfile and run these commands:

```
.../receiving_service:~ sudo docker build -t my-app . (Same with other services)

```

Step 4. Creating .env files

Create a file in your cloned repositories and call file.env
In this file you need to declare environmental variables in this format:
```
VAR_NAME1=VAR_VALUE1
VAR_NAME=VAR_VALUE2
...
```
List of environmental valriable names are: HTTP_RECEIVER, GRPC_SENDER, PROMETHEUS_ADDRESS - receiving_service; GRPC_RECEIVER, GRPC_SENDER, PROMETHEUS_ADDRESS, DIR - pcap_statistics; GRPC_RECEIVER, PROMETHEUS_ADDRESS, DIR, DSN - saving_service.

Step 5. Running docker images as containers 

We use the command -pd to run the image in a detached mode and give it some parameters of host_port:docker_port to connect your host port to docker's. For example:

```
sudo docker run -pd 8000:80 image-name
```

Step 6. Testing

Now all services are running in a detached mode and are connected with each other. To test the whole system we need to make a request to the receiving_service. 
You can do it by:
1. making a POST request through POSTMAN to the port you specified. 
2. Similarly, make a curl request.
3. Or runnin a client service [like client.go](github.com/whym9/receiving_service/client) by doing  go run main.go

Step 7. Possible answers

