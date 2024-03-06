#!/bin/bash

save_logs() {
    timestamp=$(date +'%Y-%m-%d %H:%M:%S')
    docker-compose logs --timestamps >"/home/ec2-user/auth_app/test/logs/logfile_${timestamp}.log"
}

echo "Saving logs..."
save_logs

echo "Stopping and removing containers..."
docker-compose down -v

echo "Removing images..."
docker image rm 210742672709.dkr.ecr.us-east-1.amazonaws.com/test

echo "Logging into ECR..."
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 210742672709.dkr.ecr.us-east-1.amazonaws.com/test

echo "Starting containers..."
docker-compose up -d

echo "Checking containers..."
docker-compose ps

echo "Checking logs..."
docker-compose logs

echo "Done!"
