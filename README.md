### MATIC TASK

# Requirements to spin up the project

## 1. Install Go version 1.12

    sudo apt-get update
    sudo apt-get -y upgrade

    curl -O https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz

    tar -C /usr/local -xvf go1.11.2.linux-amd64.tar.gz

    sudo nano ~/.profile

    Paste the below lines in the ~/.profile
        export GOPATH=$HOME/go
        export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

## 2. Install Docker 
    Use the following link:
    https://docs.docker.com/install/linux/docker-ce/ubuntu/

## 3. go get --save github.com/go-sql-driver/mysql

# Steps to spin up the project

    This project has two phases:
        1. Populate the data base with 10,000 recent blocks
        2. Retrieve transactions based on user Address.

### For Phase 1:
    STEP 1:
       # Open a terminal window
       # Create a repo and clone the repo form : 
       # Start the docker using the following command.
           $  docker run --name mysqleth -p 3306:3306 -v pwd:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7
           $  docker inspect mysqleth | grep IPAddr
           $  set -a
           $  export DOCKER_IP=<copy and paste the ip form result of grep>
           $  export DB_NAME=ethBlock
           $  set +a
           $  docker exec -it ethBlock bash
           $  mysql -uroot -proot
           $  CREATE DATABASE ethBlock

         Your docker for mysql server will start running
         
    STEP 2
       # Open another terminal window 
       # Access the block folder
            $ cd block
            $ go mod download
            $ go run block.go

### NOTE: 
First phase will take some time to complete as it will download 10000 latest blocks 


## For Phase 2:
    STEP 1:
        # Access app folder
           $ cd app
    STEP 2:
        To run the server 
        $ go run app.go
        Open desired browser
        $ localhost:8080/

     

        




