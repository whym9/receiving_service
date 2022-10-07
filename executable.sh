sudo docker build -t receiver . 
sudo docker run --name=rec --env-file .env -p 8080:8080 -d receiver
