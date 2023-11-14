#!/bin/bash

set -e

os_name=$(uname)
ROOT=$(pwd)
read -r _ project_name _ <go.mod
project_name=$(echo $project_name | tr -d '\r')
app_directory="${ROOT}/apps"


printf "Available authentication systems: *\n"
printf "1. Jwt Authentication system *\n"
printf "2. Firebase Authentication system *\n"
echo "*Enter Service number you want to inject :"
read service_id

app_name="auth"

router_package_name="${app_name}_router"

jwt_placeholder_value_hash=(
  "{{app_name}}:$app_name"
  "{{route_package}}:$router_package_name" 
  "{{project_name}}:$project_name"
)

jwt_entity_path_hash=(
  "controllers:${ROOT}/apps/${app_name}"
  "dtos:${ROOT}/apps/${app_name}"
  "fx:${ROOT}/apps/${app_name}"
  "routers:${ROOT}/apps/${app_name}/${app_name}_router"
  "services:${ROOT}/apps/${app_name}"
  "jwt_auth:${ROOT}/middlewares/"
)

#create file
create_file() {
    #pass first parameter as entity name
    #pass second parameter as file to write
    #pass third parameter as service name -> jwt / firebase

    # Check if the app directory exists
    if [ ! -d "$app_directory/${app_name}" ]; then
        mkdir -p "$app_directory/${app_name}"
    fi
    entity_name=$1
    file_to_write=$2
    service_name=$3
    cat "${ROOT}/automate/templates/services/${service_name}/${entity_name}.txt" >>$file_to_write

    for item in "${jwt_placeholder_value_hash[@]}"; do
        placeholder="${item%%:*}"
        value="${item##*:}"
        if [[ $os_name == "Darwin" ]]; then
            sed -i "" "s/$placeholder/$value/g" $file_to_write
            continue
        fi
        sed -i "s/$placeholder/$value/g" $file_to_write

    done
}


inject_jwt() {
    echo "Injecting JWT authentication system"
    for entity in "${jwt_entity_path_hash[@]}"; do
   
        entity_name="${entity%%:*}"
        entity_path="${entity##*:}"
         echo $entity_name
        if [[ $entity_name == "routers" ]]; then
            cd ${app_directory}/${app_name}
            mkdir "auth_router"
        fi
        file_to_write="$entity_path/${entity_name}.go"
            create_file $entity_name $file_to_write "jwt"
    done
}

inject_firebase() {
    printf "Injecting Firebase authentication system"
    printf "Yet to implement"
}


# check if auth type is empty or not
if [ -z "$service_id" ]; then
    echo "Service Name  name is empty. Please enter a non-empty name."
    exit
fi

if [ "$service_id" == "1" ]; then
    inject_jwt
elif [ "$service_id" == "2" ]; then
    inject_firebase
else
    echo "Invalid service id"
    exit
fi



