#!/usr/bin/env bash

function check_prereqs() {
    terraform=$(which terraform)
    if [[ $? -ne 0 ]]; then
        echo "Error: cannot locate the program 'terraform' in your path!"
        exit 1
    fi
    sha256sum=$(which sha256sum)
    if [[ $? -ne 0 ]]; then
        echo "Error: cannot locate the program 'sha256sum' in your path!"
        exit 1
    fi
}

function git_clone() {
    if [ -d inlets-aws-ec2-terraform ]; then
        echo "Removing git repository directory inlets-aws-ec2-terraform"
        rm -rf inlets-aws-ec2-terraform
    fi
    echo
    echo "-----------------------------------------------------------------------------------"
    echo
    echo "Running 'git clone' command..."
    git clone https://github.com/mbacchi/inlets-aws-terraform.git >/dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        echo "Error: cannot clone the GitHub repository https://github.com/mbacchi/inlets-aws-terraform.git"
        exit 1
    else
        echo
        echo "Successfully cloned GitHub repository https://github.com/mbacchi/inlets-aws-terraform.git"
    fi
}

function setup_exec_terraform() {
    if [ ! -f ${LIGHTSAIL_KEYPAIR} ]; then
        ssh-keygen -t rsa -b 4096 -N '' -f ${LIGHTSAIL_KEYPAIR}
        if [[ $? -ne 0 ]]; then
            echo "Error: cannot create keypair ${LIGHTSAIL_KEYPAIR}"
            exit 1
        fi
    fi

    if [[ "${LIGHTSAIL_KEYPAIR}" != "lightsail_keypair" ]]; then
        sed -i "s/lightsail_keypair/${LIGHTSAIL_KEYPAIR}/" main.tf
    fi

    grep ${LIGHTSAIL_KEYPAIR} main.tf >/dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        echo "Error: keypair ${LIGHTSAIL_KEYPAIR} not found in main.tf! Exiting..."
        exit 1
    fi

    echo
    echo "-----------------------------------------------------------------------------------"
    echo
    echo "Running 'terraform init' command..."
    $terraform init 2>&1 | tee -a terraform.log
    if [[ $? -ne 0 ]]; then
        echo "Error: cannot run 'terraform init'. Look at log file 'terraform.$$.log' for errors."
        exit 1
    else
        echo
        echo "SUCCESS: 'terraform init' command successful."
    fi

    echo
    echo "-----------------------------------------------------------------------------------"
    echo
    echo "Running 'terraform plan' command..."
    token=$(head -c 16 /dev/urandom | sha256sum | cut -d" " -f1); \
        $terraform plan -out=terraform_plan_inlets.$$.out  -var "token=$token" 2>&1 | tee -a terraform.log
    if [[ $? -ne 0 ]]; then
        echo "Error: cannot run 'terraform plan'. Look at log file 'terraform.$$.log' for errors."
        exit 1
    else
        echo
        echo "SUCCESS: 'terraform plan' command successful."
    fi

    echo
    echo "-----------------------------------------------------------------------------------"
    echo
    echo "Running 'terraform apply' command..."
    $terraform apply terraform_plan_inlets.$$.out 2>&1 | tee -a terraform.log
    if [[ $? -ne 0 ]]; then
        echo "Error: cannot run 'terraform apply'. Look at log file 'terraform.$$.log' for errors."
        exit 1
    else
        echo
        echo "SUCCESS: 'terraform apply' command successful."
    fi
}


if [[ ${#} -lt "4" ]]; then
    echo
    echo "usage: $0 <lightsail_keypair> <aws_profile> <aws_region> <upstream_with_port>"
    echo
    echo "where:"
    echo -e "<lightsail_keypair> the name of an existing Lightsail keypair, if you haven't created it \c"
    echo "yet we will create it using the command:"
    echo "      ssh-keygen -t rsa -b 4096 -N '' -f <lightsail_keypair>"
    echo "<aws_profile> is the AWS profile previously configured in your ~/.aws/credentials file"
    echo "<aws_region> is the AWS region you will be deploying this EC2 instance"
    echo "<upstream_with_port> is the application with the port running on 127.0.0.1 (example: https://127.0.0.1:4000)"
    exit 1
fi

LIGHTSAIL_KEYPAIR="${1}"
AWS_PROFILE="${2}"
AWS_REGION="${3}"
UPSTREAM="${4}"
export AWS_PROFILE AWS_REGION

check_prereqs
git_clone

save_dir=$(pwd)
cd inlets-aws-terraform/terraform/lightsail

setup_exec_terraform

TOKEN=$(cat terraform.log | grep -E "^inlets_token =" | sed -r 's/^inlets_token = (\w+)$/\1/')
IP_ADDRESS=$(cat terraform.log | grep -E "^public_ip_address =" | sed -r 's/^public_ip_address = (.*)$/\1/')

echo "-----------------------------------------------------------------------------------"
echo
echo
echo "To connect to Inlets run the command:"
echo
echo "-----------------------------------------------------------------------------------"
echo "inlets client  --remote=$IP_ADDRESS:80 --upstream=$UPSTREAM --token $TOKEN"
echo "-----------------------------------------------------------------------------------"

echo
echo
echo "To tear down the terraform infrastructure, run the following commands:"
echo
echo "cd ${save_dir}/inlets-aws-terraform/terraform/lightsail"
echo "terraform destroy -var token=${TOKEN}"
echo
echo "You will be prompted to answer \"yes\" in order to actually remove the terraform resources."
echo

cd $save_dir
