#!/bin/bash

set -e



os_name=$(uname)
ROOT=$(pwd)
read -r _ project_name _ <go.mod
project_name=$(echo $project_name | tr -d '\r')
app_directory="${ROOT}/apps"


printf " Enter Service id you want to inject:  *\n"
printf "1. Jwt Authentication system *\n"
printf "2. Firebase Authentication system *\n"
echo "*Enter id :"
read service_id



inject_jwt() {
 echo "Jwt called"
    
}

inject_firebase() {
    echo "Firebase called"

}


# check if auth type is empty or not
if [ -z "$service_id" ]; then
    echo "Service Name  name is empty. Please enter a non-empty name."
    exit
fi

# Check the value of the variable
if [ "$service_id" == "1" ]; then
    inject_jwt
elif [ "$service_id" == "2" ]; then
    inject_firebase
else
    echo "Invalid service id"
    exit
fi


# Check if the app directory exists
if [ ! -d "$app_directory/auth" ]; then
    mkdir -p "$app_directory/auth"
fi


