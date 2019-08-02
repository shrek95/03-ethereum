### MATIC TASK

## Requirements to spin up the project

# 1. Install Go version 1.12

    sudo apt-get update
    sudo apt-get -y upgrade

    curl -O https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz

    tar -C /usr/local -xvf go1.11.2.linux-amd64.tar.gz

    sudo nano ~/.profile

    Paste the below lines in the ~/.profile
        export GOPATH=$HOME/go
        export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
# 2. Install Go-Ethereum
    https://github.com/ethereum/go-ethereum/wiki/Installation-Instructions-for-Ubuntu
# 3. Install Docker 
    https://docs.docker.com/install/linux/docker-ce/ubuntu/

# 4. go get --save github.com/go-sql-driver/mysql

## Steps to spin up the project

    This project has two phases:
        1. Populate the data base with 10,000 recent blocks
        2. Retrieve transactions based on user Address.

For Phase 1:
    ## STEP 1:
       # Create a repo and clone the repo form : 
        
       # Start the docker using the following command.
           $  docker run --name mysqleth -v pwd:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7
           $  docker inspect mysqleth | grep IPAddr
           $  set -a
           $  export docker_ip=<copy and paste the ip form result of grep>
           $  set +a

         Your docker for mysql server will start running 
     ## STEP 2
       # Access the block folder
            $ cd block
            $ go run block.go
NOTE: First phase will take some time to complete as it will download 10000 latest blocks 

For Phase 2:
    ## STEP 1:
        # Access app folder
           $ cd app
    ## STEP 2:
        To run the server 
        $ go run app.go
        Open desired browser
        $ localhost:8080/


        

        




