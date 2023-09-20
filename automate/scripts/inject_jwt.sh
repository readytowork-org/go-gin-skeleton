#!/bin/bash

set -e

printf "\n *** Go Gin GORM Scaffold Generator *** \n"
printf "This scaffolder injects jwt authentication service to enable authenticating user with jwt token.\n"
printf "\n* Injecting Jwt authtication *\n\n"

create_file() {
    #pass first parameter as entity name
    #pass second parameter as file to write
    entity_name=$1
    file_to_write=$2

    cat "${ROOT}/automate/services/jwt/${entity_name}.txt" >>$file_to_write

    for item in "${placeholder_value_hash[@]}"; do
        placeholder="${item%%:*}"
        value="${item##*:}"
        if [[ $os_name == "Darwin" ]]; then
            sed -i "" "s/$placeholder/$value/g" $file_to_write
            continue
        fi
        sed -i "s/$placeholder/$value/g" $file_to_write

    done
    echo $file_to_write "created."
}

os_name=$(uname)

read -r _ project_name _ <go.mod
project_name=$(echo $project_name | tr -d '\r')
service_name="jwt_auth"
ROOT=$(pwd)

placeholder_value_hash=(
    "{{project_name}}:$project_name"
)

entity_path_hash=(
  "controllers:${ROOT}/api/controllers"
  "services:${ROOT}/api/services"
  "middlewares:${ROOT}/api/middlewares"
  "routes:${ROOT}/api/routes"
  "dtos:${ROOT}/dtos"

)

# #### creating file
for entity in "${entity_path_hash[@]}"; do
    entity_name="${entity%%:*}"
    entity_path="${entity##*:}"
    file_to_write="$entity_path/${service_name}.go"

    ## create file if not exists
    if ! test -f "$file_to_write"; then
        create_file $entity_name $file_to_write
    else
        echo "${file_to_write} exists."
    fi
done

### inject constructor to fx module
BASE_CONSTRUCTOR="NewJwtAuth"

fx_path_hash=(
  "Controller:${ROOT}/api/controllers/controllers.go"
  "Service:${ROOT}/api/services/services.go"
  "MiddleWare:${ROOT}/api/middlewares/middlewares.go"
)

fx_init_string="var Module = fx.Options("
for deps_value in "${fx_path_hash[@]}"; do
  deps_name="${deps_value%%:*}"
  deps_path="${deps_value##*:}"
    echo "${deps_name} name"
  echo "${deps_path} path"
  if [[ $os_name == "Darwin" ]]; then
    sed -i "" "s/${fx_init_string}/${fx_init_string}\n\t  fx.Provide(${BASE_CONSTRUCTOR}${deps_name}),/g" $deps_path
    continue
  fi
  sed -i "s/${fx_init_string}/${fx_init_string}\n\t  fx.Provide(${BASE_CONSTRUCTOR}${deps_name}),/g" $deps_path
  echo $deps_path "updated."
done

