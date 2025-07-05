#!/bin/bash

# shellcheck disable=1111

# https://stackoverflow.com/questions/59895/how-do-i-get-the-directory-where-a-bash-script-is-located-from-within-the-script
SOURCE=${BASH_SOURCE[0]}
while [ -L "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
    DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )
    SOURCE=$(readlink "$SOURCE")
    [[ $SOURCE != /* ]] && SOURCE=$DIR/$SOURCE # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )

cd "$DIR/../" || exit;

SUDO_NOTE="ERROR: This script must be ran as root!"
INSTALL_PATH="%INSTALL_PATH%"




if [ -z "$1" ]; then
    echo "Subcommand missing!";
    exit 1;
fi


if [ "$1" = "install" ]; then
    if [ "$EUID" -ne 0 ]; then
        echo "$SUDO_NOTE";
        exit 1;
    fi
    if [ -e $INSTALL_PATH ]; then
        echo "Already installed!";
        exit 1;
    fi
    
    echo "Downloading Linux binaries...";
    mkdir -p "./bin/";
    wget 'https://github.com/Team-GhostLand/Grinch/releases/download/mvp-1rv/linux.zip'
    if [ $? -ne 0 ]; then
        echo "Download failed!"
        exit 1;
    fi
    mv "linux.zip" "./bin/bin.zip";
    
    echo "Unzipping Linux binaries...";
    cd "./bin/" || exit
    unzip "bin.zip";
    if [ $? -ne 0 ]; then
        echo "Unzip failed!"
        exit 1;
    fi
    cd ..;
    
    echo "Checking if you're on non-x64 (will assume ARM if so)..."
    BINNAME="grinch"
    UNAME=$(uname -m)
    if [ "$UNAME" != "x86_64" ]; then
        BINNAME="grinch-arm"
    fi
    echo "Archname: $UNAME  -  will use bin: $BINNAME"
    
    echo "Removing trash..."
    rm -v "./bin/bin.zip";
    rm -rv "./cmd/";
    rm -rv "./testing_assets/";
    rm -rv "./trans/";
    rm -rv "./util/";
    rm -v "./.gitignore";
    rm -v "./go.mod";
    rm -v "./go.sum";
    rm -v "./main.go";
    
    echo "Linking binaries and scripts...";
    chmod 755 --verbose "$(pwd)/bin/$BINNAME";
    ln --symbolic --verbose "$(pwd)/bin/$BINNAME" "$INSTALL_PATH";
    ln --symbolic --verbose "$(pwd)/scripts/installer.sh" "$INSTALL_PATH-manager";
    ln --symbolic --verbose "$(pwd)/scripts/make_serverpack.sh" "$INSTALL_PATH-serverpack";
    
    echo "Installed!!!  🎉"
    exit 0;
fi


if [ "$1" = "uninstall" ]; then
    if [ "$EUID" -ne 0 ]; then
        echo "$SUDO_NOTE";
        exit 1;
    fi

    echo "Unlinking binaries and scripts...";
    unlink "$INSTALL_PATH";
    unlink "$INSTALL_PATH-serverpack";
    unlink "$INSTALL_PATH-manager";

    echo "Removing this Spectre project altogether...";
    PROJECT_LOCATION=$(pwd)
    cd ..
    rm -rv "$PROJECT_LOCATION"
    exit;
fi


if [ "$1" = "update" ]; then
    if [ "$EUID" -ne 0 ]; then
        echo "$SUDO_NOTE";
        exit 1;
    fi
    
    echo "  ---- UNINSTALLING ----";
    ./scripts/installer.sh uninstall;
    if [ $? -ne 0 ]; then
        echo "Uninstall failed!"
        exit 1;
    fi
    cd .. #CDs don't carry over between shell excutions - we need to once again exit from the (now-deleted) project folder into Specre's main dir, or sudo -E bash will complain about starting from a non-existent directory.
    
    echo "  ---- CALLING UPON SPECTRE TO INSTALL THE NEWEST VERSION ----";
    if [ -z "$PROJECT_NAME" ]; then
        PROJECT_NAME="grinch"
    fi
    sleep 3;
    export PROJECT_NAME=$PROJECT_NAME INSTALL_PATH=$INSTALL_PATH SCRIPT_NAME="scripts/installer" GIT="https://github.com/Team-GhostLand/Grinch.git" && curl -fsSL https://raw.githubusercontent.com/Team-GhostLand/Spectre/master/universal-installer-scaffolding.sh | sudo -E bash
    exit;
fi


if [ "$1" = "help" ]; then
    echo "Manages your $INSTALL_PATH installation.";
    echo "USAGE: $INSTALL_PATH-manager <help|update|uninstall>"; #install is for Spectre - let's not mention it here
    echo "Note: When updating, please add Specte's project and path variables, if you don't want to use the defaults, as this script isn't aware of its own project name, so we can't pass it automatically."
    exit 0;
fi


echo "Unknown sub-command: $1. Use „help” for help.";
exit 1;