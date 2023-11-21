#!/bin/bash

set -e

os_name=$(uname)
ROOT=$(pwd)
read -r _ project_name _ <go.mod
project_name=$(echo $project_name | tr -d '\r')
app_directory="${ROOT}/apps"


printf "\nAvailable authentication systems:\n"
printf "1. Jwt Authentication system \n"
printf "2. Firebase Authentication system \n"
printf "\n*Enter Service number you want to inject: "
read service_id

app_name="auth"


jwt_placeholder_value_hash=(
  "{{app_name}}:$app_name"
  "{{project_name}}:$project_name"
)

jwt_entity_path_hash=(
  "controllers:${ROOT}/apps/${app_name}/controllers"
  "dtos:${ROOT}/apps/${app_name}"
  "init:${ROOT}/apps/${app_name}/init"
  "routes:${ROOT}/apps/${app_name}/routes"
  "services:${ROOT}/apps/${app_name}/services"
  "jwt_auth:${ROOT}/middlewares"
)

config_path="${ROOT}/config/conf.go"
middleware_path="${ROOT}/middlewares/init.go"
router_path="${ROOT}/config/router.go"
import_name="${project_name}/apps/${app_name}/init"
import_name_router="${project_name}/apps/${app_name}/routes"
jwt_fx_installed_app_string="var InstalledApps = fx.Options("
jwt_middleware_fx_module_string="var Module = fx.Options("


jwt_inject_dependency() {
    if [[ $os_name == "Darwin" ]]; then
        sed -i '' -e "/^import (/a\\
        ${app_name} \"$import_name\"
        " $config_path

        # installed app
        sed -i "" "s/${jwt_fx_installed_app_string}/${jwt_fx_installed_app_string}\n\t  ${app_name}.Module,/g" $config_path

        # middleware fx
        sed -i "" "s/${jwt_middleware_fx_module_string}/${jwt_middleware_fx_module_string}\n\t  fx.Provide(NewJWTAuthMiddleWare),/g" $middleware_path

        # router
        sed -i '' -e "/^import (/a\\
        ${app_name} \"$import_name_router\"
        " $router_path
        sed -i "" "s/func RoutersConstructor(/func RoutersConstructor(\n\t AuthRoutes ${app_name}.AuthRoute,/g" $router_path
        sed -i "" "s/return Routes{/return Routes{\n\t AuthRoutes,/g" $router_path
    else
        # installed app
        sed -i "s/${jwt_fx_installed_app_string}/${jwt_fx_installed_app_string}\n\t  ${app_name}.Module,/g" $config_path
        # middleware fx
        sed -i "s/${jwt_middleware_fx_module_string}/${jwt_middleware_fx_module_string}\n\t  fx.Provide(NewJWTAuthMiddleWare),/g" $middleware_path
        # router
        sed -i "s/func RoutersConstructor(/func RoutersConstructor(\n\t AuthRoutes ${app_name}.AuthRoute,/g" $router_path
        sed -i "s/return Routes{/return Routes{\n\t AuthRoutes,/g" $router_path

    fi
    echo "✅ Jwt authentication system 🔐 injected "
    echo "You can now use middlewares.JWTAuthMiddleWare.Handle() method to enable jwt auth middleware"

}
#create file
create_file() {
    #pass first parameter as entity name
    #pass second parameter as file to write
    #pass third parameter as service name -> jwt / firebase

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
    # Check if the app directory exists
    if [ ! -d "$app_directory/${app_name}" ]; then
        mkdir -p "$app_directory/${app_name}"
        cd ${app_directory}/${app_name}
        mkdir controllers services routes init
    fi
    for entity in "${jwt_entity_path_hash[@]}"; do
        entity_name="${entity%%:*}"
        entity_path="${entity##*:}"
        file_to_write="$entity_path/${entity_name}.go"
            create_file $entity_name $file_to_write "jwt"
    done
    jwt_inject_dependency
}

inject_firebase() {
    printf "Injecting Firebase authentication system"
    printf "Yet to implement "
}

# check if auth type is empty or not
if [ -z "$service_id" ]; then
    echo "🛑 Service Name id is empty. Please enter a non-empty name."
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
go mod tidy



