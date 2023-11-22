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


# check if auth type is empty or not
if [ -z "$service_id" ]; then
    echo "🛑 Service Name id is empty. Please enter a non-empty name."
    exit
fi


app_name="auth"

placeholder_value_hash=(
  "{{app_name}}:$app_name"
  "{{project_name}}:$project_name"
)

entity_path_hash=(
  "controllers:${ROOT}/apps/${app_name}/controllers"
  "dtos:${ROOT}/apps/${app_name}"
  "init:${ROOT}/apps/${app_name}/init"
  "routes:${ROOT}/apps/${app_name}/routes"
  "services:${ROOT}/apps/${app_name}/services"
  "jwt_auth:${ROOT}/middlewares"
)

if [ "$service_id" == "2" ]; then
    entity_path_hash=(
    "init:${ROOT}/apps/${app_name}/init"
    "services:${ROOT}/apps/${app_name}/services"
    "firebase_auth:${ROOT}/middlewares"
    )

fi



# common 
config_path="${ROOT}/config/conf.go"
middleware_path="${ROOT}/middlewares/init.go"
router_path="${ROOT}/config/router.go"
import_name="${project_name}/apps/${app_name}/init"
import_name_router="${project_name}/apps/${app_name}/routes"
fx_installed_app_string="var InstalledApps = fx.Options("
middleware_fx_module_string="var Module = fx.Options("


jwt_inject_dependency() {
    if [[ $os_name == "Darwin" ]]; then
        sed -i '' -e "/^import (/a\\
        ${app_name} \"$import_name\"
        " $config_path

        # installed app
        sed -i "" "s/${fx_installed_app_string}/${fx_installed_app_string}\n\t  ${app_name}.Module,/g" $config_path

        # middleware fx
        sed -i "" "s/${middleware_fx_module_string}/${middleware_fx_module_string}\n\t  fx.Provide(NewJWTAuthMiddleWare),/g" $middleware_path

        # router
        sed -i '' -e "/^import (/a\\
        ${app_name} \"$import_name_router\"
        " $router_path
        sed -i "" "s/func RoutersConstructor(/func RoutersConstructor(\n\t AuthRoutes ${app_name}.AuthRoute,/g" $router_path
        sed -i "" "s/return Routes{/return Routes{\n\t AuthRoutes,/g" $router_path
    else
        # installed app
        sed -i "s/${fx_installed_app_string}/${fx_installed_app_string}\n\t  ${app_name}.Module,/g" $config_path
        # middleware fx
        sed -i "s/${middleware_fx_module_string}/${middleware_fx_module_string}\n\t  fx.Provide(NewJWTAuthMiddleWare),/g" $middleware_path
        # router
        sed -i "s/func RoutersConstructor(/func RoutersConstructor(\n\t AuthRoutes ${app_name}.AuthRoute,/g" $router_path
        sed -i "s/return Routes{/return Routes{\n\t AuthRoutes,/g" $router_path

    fi
    echo "✅ Jwt authentication system 🔐 injected "
    echo "You can now use middlewares.JWTAuthMiddleWare.Handle() method to enable jwt auth"

}

fb_inject_dependency() {
    if [[ $os_name == "Darwin" ]]; then
        sed -i '' -e "/^import (/a\\
        ${app_name} \"$import_name\"
        " $config_path

        # installed app
        sed -i "" "s/${fx_installed_app_string}/${fx_installed_app_string}\n\t  ${app_name}.Module,/g" $config_path

        # middleware fx
        sed -i "" "s/${middleware_fx_module_string}/${middleware_fx_module_string}\n\t  fx.Provide(NewFirebaseAuthMiddleware),/g" $middleware_path

    else
        # installed app
        sed -i "s/${fx_installed_app_string}/${fx_installed_app_string}\n\t  ${app_name}.Module,/g" $config_path
        # middleware fx
        sed -i "s/${middleware_fx_module_string}/${middleware_fx_module_string}\n\t  fx.Provide(NewFirebaseAuthMiddleware),/g" $middleware_path

    fi
    echo "✅ Firebase authentication system 🔐 injected "
    echo "You can now use middlewares.FirebaseAuthMiddleware.Handle() method to enable firebase auth"

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

    for item in "${placeholder_value_hash[@]}"; do
        placeholder="${item%%:*}"
        value="${item##*:}"
        if [[ $os_name == "Darwin" ]]; then
            sed -i "" "s/$placeholder/$value/g" $file_to_write
            continue
        fi
        sed -i "s/$placeholder/$value/g" $file_to_write

    done
}


inject_auth() {
    if [[ "$service_id" != "1" && "$service_id" != "2" ]]; then
        echo "Invalid service id"
        exit
    fi
    service_name="jwt"
    # Check if the app directory exists
    if [ ! -d "$app_directory/${app_name}" ]; then
        mkdir -p "$app_directory/${app_name}"
        cd ${app_directory}/${app_name}
        if [ "$service_id" == "1" ]; then
            echo "Injecting JWT authentication system"
            mkdir controllers services routes init
        else
            service_name="firebase"
            echo "Injecting Firebase authentication system"
            mkdir services init
        fi
    fi

    for entity in "${entity_path_hash[@]}"; do
        entity_name="${entity%%:*}"
        entity_path="${entity##*:}"
        file_to_write="$entity_path/${entity_name}.go"
            create_file $entity_name $file_to_write $service_name
    done

    if [ "$service_id" == "1" ]; then
        jwt_inject_dependency
    else
        fb_inject_dependency
    fi
}


inject_auth

go mod tidy



