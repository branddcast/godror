#!/bin/bash

Custom_machine(){

    new_machine=$1

    export ENV_BAP_GO_ACTUAL_MACHINE=$(echo $HOSTNAME)

    echo "Change machine from $ENV_BAP_GO_ACTUAL_MACHINE to $new_machine"


    file=/etc/hostname

    if ! sed -i 's/'"$ENV_BAP_GO_ACTUAL_MACHINE"'/'"$new_machine"'/g' $file ; then
        printf "\nsed /etc/hostname returned an error."
    else
        printf "\nsed /etc/hostname successfully"
    fi

}

while getopts c:r: flag
do
    case "${flag}" in
        c) Custom_machine $OPTARG ;;
    esac
done