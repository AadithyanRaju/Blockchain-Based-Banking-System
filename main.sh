#!/bin/bash
ROOTDIR=$(pwd)
GITHUBID="AadithyanRaju"
PROJECTNAME="BankingSystem"
CONTRACTNAME="banking"

#check if 'go' directory exists or not
if [ ! -f "$ROOTDIR/main.sh" ]; then
    echo "Go directory does not exist."
    echo "Run this script in the root directory of the project."
    exit 1
fi

install(){
    packages="curl git nodejs npm docker docker-compose"
    aptpkgs="golang"
    pacmanpkgs="go"
    
    if [ -x "$(command -v apt-get)" ]; then
        sudo apt update
        sudo apt install $packages $aptpkgs
    elif [ -x "$(command -v pacman)" ]; then
        sudo pacman -Syy
        sudo pacman -S $packages $pacmanpkgs
    else
        echo "Invalid package manager."
        echo "Please install the required packages manually."
        echo "Packages: $packages $aptpkgs"
        exit 1
    fi
}

while true; do
    clear
    echo "Menu:"
    echo "1. Install Required Packages"
    echo "2. Setup the workspace"
    echo "3. Start the Blockchain"
    echo "4. Update Chaincode to the Blockchain"
    echo "5. Stop the Blockchain"
    echo "6. List Docker Containers"
    echo "8. Reset the workspace"
    echo "9. Exit"
    read -p "Enter your choice: " choice

    case $choice in
        1)
            echo "Docker, Docker Compose, GoLang, and NodeJS needs be installed."
            echo "Do you want to install the required packages? (y/n)"
            read -p "Enter your choice: " install_choice
            if [ "$install_choice" == "y" ]; then
                install
            fi
            
            ;;
        2)
            echo "Setting up the workspace..."
            if [ -d "$ROOTDIR/go/src/github.com/$GITHUBID" ]; then
                echo "Workspace already exists."
            else
                mkdir -p "$ROOTDIR/go/src/github.com/$GITHUBID"
            fi
            cd "$ROOTDIR/go/src/github.com/$GITHUBID"
            
            if [ -f "install-fabric.sh" ]; then
                echo "Hyperledger Fabric Installer already exists."
            else
                curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
            fi
            
            if [ -d "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples" ]; then
                echo "Fabric Samples already exists."
            else
                ./install-fabric.sh d s b # -f 2.5.9 -c 1.5.12

            fi

            if [ ! -d "$ROOTDIR/go/src/github.com/$GITHUBID/chaincode/$CONTRACTNAME/go" ]; then
                echo "Setting up Chaincode..."
                mkdir -p "$ROOTDIR/go/src/github.com/$GITHUBID/chaincode/$CONTRACTNAME/go"
                cd "$ROOTDIR/go/src/github.com/$GITHUBID/chaincode/$CONTRACTNAME/go"
                go mod init $CONTRACTNAME
                cp $ROOTDIR/$CONTRACTNAME.go .
                go get github.com/hyperledger/fabric-contract-api-go/contractapi
                go mod tidy
                go mod vendor
                cd $ROOTDIR
            else
                echo "Chaincode already exists."
            fi

            if [ ! -d "$ROOTDIR/go/src/github.com/$GITHUBID/$PROJECTNAME" ]; then
                echo "Setting up Application..."
                mkdir -p "$ROOTDIR/go/src/github.com/$GITHUBID/$PROJECTNAME"
                cd "$ROOTDIR/go/src/github.com/$GITHUBID/$PROJECTNAME"
                # uncomment the below line if you want to use hardhat for ethereum development
                #npx hardhat init 
                npm install fabric-network fabric-ca-client
                cd $ROOTDIR
            else
                echo "$PROJECTNAME already exists."
                echo "If you want to reset the workspace, please delete the 'go' directory and run this script again."
            fi

            echo "Updating the $PROJECTNAME contracts..."
            cp $ROOTDIR/$CONTRACTNAME.go $ROOTDIR/go/src/github.com/$GITHUBID/chaincode/$CONTRACTNAME/go/$CONTRACTNAME.go

            ;;
        3)
            echo "Starting the Blockchain..."
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples/test-network"
            export GO111MODULE="on"
            export GOPROXY="https://proxy.golang.org"
            export GOMODCACHE="$ROOTDIR/go/pkg/mod"
            export GOPATH="$ROOTDIR/go"
            export GOROOT="/usr/lib/go"
            if [ $(docker ps | grep hyperledger | wc -l) -eq 0 ]; then
                ./network.sh up -ca -s couchdb
                ./network.sh createChannel -c mychannel
            fi
            cd -
            ;;
        4)
            echo "Adding Chaincode to the Blockchain..."
            cp $ROOTDIR/$CONTRACTNAME.go $ROOTDIR/go/src/github.com/$GITHUBID/chaincode/$CONTRACTNAME/go/$CONTRACTNAME.go -v
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples/test-network"
            ./network.sh deployCC -ccn $CONTRACTNAME -ccp $ROOTDIR/go/src/github.com/$GITHUBID/chaincode/$CONTRACTNAME/go -ccl go
            cd -
            ;;

        5)
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples/test-network"
            ./network.sh down
            cd -
            ;;
        6)
            echo "Listing Docker Containers..."
            docker ps -a
            ;;
        8)
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples/test-network"
            ./network.sh down
            echo "Resetting the workspace..."
            sudo rm -rf "$ROOTDIR/go"
            echo "Workspace has been reset."
            ;;
        9)
            echo "Exiting..."
            exit 0
            ;;
        *)
            echo "Invalid choice. Please try again."
            ;;
    esac

    read -p "Press Enter to continue..."
done
cd $ROOTDIR
# ln -s "$ROOTDIR/go/src/github.com/$GITHUBID/$PROJECTNAME/fabric-samples" "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples"