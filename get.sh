#!/bin/bash
# This script was adapted from https://github.com/openfaas/cli.openfaas.com/blob/master/get.sh

version=$(curl -sI https://github.com/alexellis/inlets/releases/latest | grep Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')

if [ ! $version ]; then
    echo "Failed while attempting to install inlets. Please manually install:"
    echo ""
    echo "1. Open your web browser and go to https://github.com/alexellis/inlets/releases"
    echo "2. Download the latest release for your platform. Call it 'inlets'."
    echo "3. chmod +x ./inlets"
    echo "4. mv ./inlets /usr/local/bin"
    exit 1
fi

hasCli() {

    has=$(which inlets)

    if [ "$?" = "0" ]; then
        echo
        echo "You already have the inlets cli!"
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

    targetFile="/tmp/inlets$suffix"

    if [ "$userid" != "0" ]; then
        targetFile="$(pwd)/inlets$suffix"
    fi

    if [ -e $targetFile ]; then
        rm $targetFile
    fi

    url=https://github.com/alexellis/inlets/releases/download/$version/inlets$suffix
    echo "Downloading package $url as $targetFile"

    curl -sSL $url --output $targetFile

    if [ "$?" = "0" ]; then

    chmod +x $targetFile

    echo "Download complete."

        if [ "$userid" != "0" ]; then

            echo
            echo "========================================================="
            echo "==    As the script was run as a non-root user the     =="
            echo "==    following commands may need to be run manually   =="
            echo "========================================================="
            echo
            echo "  sudo cp inlets$suffix /usr/local/bin/inlets"
            echo

        else

            echo
            echo "Running as root - Attempting to move inlets to /usr/local/bin"

            mv $targetFile /usr/local/bin/inlets

            if [ "$?" = "0" ]; then
                echo "New version of inlets installed to /usr/local/bin"
            fi

            if [ -e $targetFile ]; then
                rm $targetFile
            fi

            inlets version
        fi
    fi
}

hasCli
getPackage
