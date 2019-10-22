#!/bin/bash
# This script was adapted from https://github.com/openfaas/cli.openfaas.com/blob/master/get.sh

export OWNER=inlets
export REPO=inlets
export SUCCESS_CMD="$REPO version"
export BINLOCATION="/usr/local/bin"

version=$(curl -sI https://github.com/$OWNER/$REPO/releases/latest | grep Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')

if [ ! $version ]; then
    echo "Failed while attempting to install $REPO. Please manually install:"
    echo ""
    echo "1. Open your web browser and go to https://github.com/$OWNER/$REPO/releases"
    echo "2. Download the latest release for your platform. Call it '$REPO'."
    echo "3. chmod +x ./$REPO"
    echo "4. mv ./$REPO $BINLOCATION"
    exit 1
fi

hasCli() {

    has=$(which $REPO)

    if [ "$?" = "0" ]; then
        echo
        echo "You already have the $REPO cli!"
        export n=1
        echo "Overwriting in $n seconds.. Press Control+C to cancel."
        echo
        sleep $n
    fi

    hasCurl=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        exit 1
    fi
}

getPackage() {
    uname=$(uname)
    userid=$(id -u)

    suffix=""
    case $uname in
    "Darwin")
    suffix="-darwin"
    ;;
    "Linux")
        arch=$(uname -m)
        echo $arch
        case $arch in
        "aarch64")
        suffix="-arm64"
        ;;
        esac
        case $arch in
        "armv6l" | "armv7l")
        suffix="-armhf"
        ;;
        esac
    ;;
    esac

    targetFile="/tmp/$REPO$suffix"

    if [ "$userid" != "0" ]; then
        targetFile="$(pwd)/$REPO$suffix"
    fi

    if [ -e $targetFile ]; then
        rm $targetFile
    fi

    url=https://github.com/inlets/$REPO/releases/download/$version/$REPO$suffix
    echo "Downloading package $url as $targetFile"

    curl -sSLf $url --output $targetFile

    if [ "$?" = "0" ]; then

    chmod +x $targetFile

    echo "Download complete."

        if [ ! -w "$BINLOCATION" ]; then

            echo
            echo "============================================================"
            echo "==   The script was run as a user who is unable to write  =="
            echo "==   to $BINLOCATION. To complete the installation the  =="
            echo "==   following commands may need to be run manually.      =="
            echo "============================================================"
            echo
            echo "  sudo cp $REPO$suffix $BINLOCATION/$REPO"
            echo
            ./$REPO$suffix version
        else

            echo
            echo "Running with sufficient permissions to attempt to move $REPO to $BINLOCATION"

            mv $targetFile $BINLOCATION/$REPO

            if [ "$?" = "0" ]; then
                echo "New version of $REPO installed to $BINLOCATION"
            fi

            if [ -e $targetFile ]; then
                rm $targetFile
            fi

           ${SUCCESS_CMD}
        fi
    fi
}

hasCli
getPackage

