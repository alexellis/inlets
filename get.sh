#!/bin/bash
# This script was adapted from https://github.com/openfaas/cli.openfaas.com/blob/master/get.sh

export OWNER=inlets
export REPO=inlets
export SUCCESS_CMD="$REPO version"
export BINLOCATION="/usr/local/bin"

version=$(curl -sI https://github.com/$OWNER/$REPO/releases/latest | grep -i Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')

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
    "MINGW"*)
    suffix=".exe"
    BINLOCATION="$HOME/bin"
    mkdir -p $BINLOCATION

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

    targetFile="/tmp/$REPO"

    if [ "$userid" != "0" ]; then
        targetFile="$(pwd)/$REPO"
    fi

    if [ -e $targetFile ]; then
        rm $targetFile
    fi

    url=https://github.com/$OWNER/$REPO/releases/download/$version/$REPO$suffix
    echo "Downloading package $url as $targetFile"

    curl -sSLf $url --output $targetFile

    if [ "$?" = "0" ]; then

    chmod +x $targetFile

    echo "Download complete."

        if [ ! -w "$BINLOCATION" ]; then

            echo
            echo "============================================================"
            echo "  The script was run as a user who is unable to write"
            echo "  to $BINLOCATION. To complete the installation the"
            echo "  following commands may need to be run manually."
            echo "============================================================"
            echo
            echo "  sudo cp $REPO $BINLOCATION/$REPO"
            echo

        else

            echo
            echo "Running with sufficient permissions to attempt to move $REPO to $BINLOCATION"

            if [ ! -w "$BINLOCATION/$REPO" ] && [ -f "$BINLOCATION/$REPO" ]; then

            echo
            echo "================================================================"
            echo "  $BINLOCATION/$REPO already exists and is not writeable"
            echo "  by the current user.  Please adjust the binary ownership"
            echo "  or run sh/bash with sudo." 
            echo "================================================================"
            echo
            exit 1

            fi

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
