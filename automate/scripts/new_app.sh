#!/bin/bash

set -e

printf "======================\n"
printf "Creating new app\n"
printf "======================\n"

os_name=$(uname)
ROOT=$(pwd)
read -r _ project_name _ <go.mod
project_name=$(echo $project_name | tr -d '\r')
app_directory="${ROOT}/apps"

printf "*App name should be in snake case eg: my_app *\n"
echo "Enter app name:"
read app_name

# Check if the app exists
if [ -d "$app_directory/${app_name}" ]; then
  echo "$app_name already exists try with different app name."
  exit
else
  mkdir ${app_directory}/${app_name}
fi

make_pascal_case_str() {
  pascal_case=""
  # Split the string into an array using '_' as the delimiter
  IFS="_" read -ra words <<< "$1" #spliting string by _
  # Loop through the words and capitalize the first letter of each word
  for word in "${words[@]}"; do
    # Capitalize the first letter of the word
    capitalized_word=$(echo "$word" | awk '{print toupper(substr($0, 1, 1)) tolower(substr($0, 2))}')
    pascal_case="${pascal_case}${capitalized_word}"
  done

  echo "$pascal_case"
}

create_package_name(){
    output=$(echo "$1" | tr -d '_')
    echo "$output"
}

create_file() {
    #pass first parameter as entity name
    #pass second parameter as file to write
    entity_name=$1
    file_to_write=$2

    cat "${ROOT}/automate/templates/new_app/${entity_name}.txt" >>$file_to_write

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


package_name=$(create_package_name $app_name)
method_name=$(make_pascal_case_str $app_name)


placeholder_value_hash=(
  "{{package_name}}:$package_name"
  "{{app_name}}:$app_name"
  "{{project_name}}:$project_name"
  "{{app_uppercase}}:$method_name"
)

entity_path_hash=(
  "controllers:${ROOT}/apps/${app_name}"
  "dtos:${ROOT}/apps/${app_name}"
  "fx:${ROOT}/apps/${app_name}"
  "helpers:${ROOT}/apps/${app_name}"
  "models:${ROOT}/apps/${app_name}"
  "repository:${ROOT}/apps/${app_name}"
  "routers:${ROOT}/apps/${app_name}"
  "services:${ROOT}/apps/${app_name}"
)

for entity in "${entity_path_hash[@]}"; do
    entity_name="${entity%%:*}"
    entity_path="${entity##*:}"
    file_to_write="$entity_path/${entity_name}.go"
    create_file $entity_name $file_to_write

done



