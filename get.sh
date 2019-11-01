#!/bin/bash

# Copyright OpenFaaS Author(s) 2019
#########################
# Repo specific content #
#########################

export VERIFY_CHECKSUM=0
export ALIAS=""
export OWNER=inlets
export REPO=inlets
export SUCCESS_CMD="$REPO version"
export BINLOCATION="/usr/local/bin"

###############################
# Content common across repos #
###############################

version=$(curl -sI https://github.com/$OWNER/$REPO/releases/latest | grep Location | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')
if [ ! $version ]; then
    echo "Failed while attempting to install $REPO. Please manually install:"
    echo ""
    echo "1. Open your web browser and go to https://github.com/$OWNER/$REPO/releases"
    echo "2. Download the latest release for your platform. Call it '$REPO'."
    echo "3. chmod +x ./$REPO"
    echo "4. mv ./$REPO $BINLOCATION"
    if [ -n "$ALIAS_NAME" ]; then
        echo "5. ln -sf $BINLOCATION/$REPO /usr/local/bin/$ALIAS_NAME"
    fi
    exit 1
fi

hasCli() {

    hasCurl=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        exit 1
    fi
}

checkHash(){

    sha_cmd="sha256sum"

    if [ ! -x "$(command -v $sha_cmd)" ]; then
    sha_cmd="shasum -a 256"
    fi

    if [ -x "$(command -v $sha_cmd)" ]; then

    targetFileDir=${targetFile%/*}

    (cd $targetFileDir && curl -sSL $url.sha256|$sha_cmd -c >/dev/null)
   
        if [ "$?" != "0" ]; then
            rm $targetFile
            echo "Binary checksum didn't match. Exiting"
            exit 1
        fi   
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

    targetFile="/tmp/$REPO$suffix"
    
    if [ "$userid" != "0" ]; then
        targetFile="$(pwd)/$REPO$suffix"
    fi

    if [ -e $targetFile ]; then
        rm $targetFile
    fi

    url=https://github.com/$OWNER/$REPO/releases/download/$version/$REPO$suffix
    echo "Downloading package $url as $targetFile"

    curl -sSL $url --output $targetFile

    if [ "$?" = "0" ]; then

        if [ "$VERIFY_CHECKSUM" = "1" ]; then
            checkHash
        fi

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
            echo "  sudo cp $REPO$suffix $BINLOCATION/$REPO"
            
            if [ -n "$ALIAS_NAME" ]; then
                echo "  sudo ln -sf $BINLOCATION/$REPO $BINLOCATION/$ALIAS_NAME"
            fi
            
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

            if [ -n "$ALIAS_NAME" ]; then
                if [ ! -L $BINLOCATION/$ALIAS_NAME ]; then
                    ln -s $BINLOCATION/$REPO $BINLOCATION/$ALIAS_NAME
                    echo "Creating alias '$ALIAS_NAME' for '$REPO'."
                fi
            fi

            ${SUCCESS_CMD}
        fi
    fi
}

hasCli
getPackage
