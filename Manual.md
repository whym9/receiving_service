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

Step 3. Running services.
