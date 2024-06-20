#!/bin/bash
ROOTDIR=$(pwd)
GITHUBID="AadithyanRaju"
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
    echo "7. Enter Bash/Zsh Shell (Default: Bash, enter zsh for zsh)"
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
                ./install-fabric.sh d s b
            fi

            if [ ! -d "$ROOTDIR/go/src/github.com/$GITHUBID/chaincode/banking/go" ]; then
                echo "Setting up Chaincode..."
                mkdir -p "$ROOTDIR/go/src/github.com/$GITHUBID/chaincode/banking/go"
                cd "$ROOTDIR/go/src/github.com/$GITHUBID/chaincode/banking/go"
                go mod init banking
                cp $ROOTDIR/banking.go .
                go get github.com/hyperledger/fabric-contract-api-go/contractapi
                go mod tidy
                go mod vendor
                cd $ROOTDIR
            else
                echo "Chaincode already exists."
            fi

            if [ ! -d "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system" ]; then
                echo "Setting up Application..."
                mkdir -p "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system"
                cd "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system"
                npx hardhat init
                npm install fabric-network fabric-ca-client
                cd $ROOTDIR
            else
                echo "Banking-system already exists."
                echo "If you want to reset the workspace, please delete the 'go' directory and run this script again."
            fi

            echo "Updating the banking-system scripts..."
            if [ ! -d "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system/scripts" ]; then
                mkdir -p "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system/scripts"
            fi
            cp $ROOTDIR/Scripts/* $ROOTDIR/go/src/github.com/$GITHUBID/banking-system/scripts/

            echo "Updating the banking-system contracts..."
            cp $ROOTDIR/banking.go $ROOTDIR/go/src/github.com/$GITHUBID/chaincode/banking/go/banking.go

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
                ./network.sh up createChannel -ca
            fi
            if [ $(docker ps | grep hyperledger | wc -l) -eq 6 ]; then
                ./network.sh deployCC -ccn banking -ccp $ROOTDIR/go/src/github.com/$GITHUBID/chaincode/banking/go -ccl go
            fi
            cd -
            ;;
        4)
            echo "Adding Chaincode to the Blockchain..."
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/fabric-samples/test-network"
            ./network.sh deployCC -ccn banking -ccp $ROOTDIR/go/src/github.com/$GITHUBID/chaincode/banking/go -ccl go
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
        7)
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system/scripts"
            echo "Entering Bash Shell..."
            bash
            ;;
        "zsh")
            cd "$ROOTDIR/go/src/github.com/$GITHUBID/banking-system/scripts"
            echo "Entering Zsh Shell..."
            zsh
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