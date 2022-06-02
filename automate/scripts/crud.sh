#!/bin/bash -e

if [[ uname -eq "Linux" ]]
then
  resed="sed -i"
elif [[ uname -eq "Darwin" ]]
then
  resed="sed -i ''"
else
  resed="sed -i"
fi


first_lower () {
  echo `echo $1 | awk '{$1=tolower(substr($1,0,1))substr($1,2)}1'`
}

printf "\n *** Go Gin GORM Scaffold Generator *** \n"
printf "This scaffolder assumes that you are using RTW clean-gin template.\n"
echo "Enter resource name(eg: ProductCategory):"; read uc_resource
echo "Enter resource table name(eg: product_category):"; read resource_table
echo "Enter plural resource table name(eg: product_categories):"; read plural_resource_table
echo "Enter plural resource name(eg: ProductCategories):"; read plural_resource

lc_resource=$(first_lower $uc_resource)
plc_resource=$(first_lower $plural_resource)
ROOT=$(pwd)

printf "\n* Generating Scaffold for ${uc_resource} *\n\n"

# getting project name from go.mod file
# this code will grab second word of first line from file go.mod and store value to the project name
read -r _ project_name _ < go.mod

placeholder_value_hash=(
  "{{ucresource}}:$uc_resource"
  "{{plcresource}}:$plc_resource"
  "{{lcresource}}:$lc_resource"
  "{{projectname}}:$project_name"
  "{{resourcetable}}:$resource_table"
  "{{pluralresourcetable}}:$plural_resource_table"
)
entity_path_hash=(
  "models:${ROOT}/models"
  "routes:${ROOT}/api/routes"
  "controllers:${ROOT}/api/controllers"
  "services:${ROOT}/api/services"
  "repository:${ROOT}/api/repository"
)

# if files already exists then terminate the process 
for str in ${entity_path_hash[@]}; do
    FILE="${entity##*:}/${resource_table}.go"
    if test -f "$FILE"; then
        echo "${str}/${fileName} exists."
        exit
    fi
done


for entity in "${entity_path_hash[@]}"; do
  entity_name="${entity%%:*}"
  entity_path="${entity##*:}"
  file_to_write="$entity_path/${resource_table}.go" 

  cat "${ROOT}/automate/automate-templates/${entity_name}.txt" >> $file_to_write
  for item in "${placeholder_value_hash[@]}"; do
    placeholder="${item%%:*}"
    value="${item##*:}"
    $resed "" "s/$placeholder/$value/g" $file_to_write
  done
  echo $file_to_write "created."
done

# inject fx deps
fx_path_hash=(
  "Controller:${ROOT}/api/controllers/controllers.go"
  "Service:${ROOT}/api/services/services.go"
  "Repository:${ROOT}/api/repository/repository.go"
)
fx_init_string="var Module = fx.Options("
for deps_value in "${fx_path_hash[@]}"; do
  deps_name="${deps_value%%:*}"
  deps_path="${deps_value##*:}"
  $resed "" "s/${fx_init_string}/${fx_init_string}\n\t  fx.Provide(New${uc_resource}${deps_name}),/g" $deps_path
  echo $deps_path "updated."
done

# fx routes
fx_route_path="${ROOT}/api/routes/routes.go"
$resed "" "s/func NewRoutes(/func NewRoutes(\n\t ${lc_resource}Routes ${uc_resource}Routes,/g" $fx_route_path
$resed "" "s/return Routes{/return Routes{\n\t ${lc_resource}Routes,/g" $fx_route_path
$resed "" "s/fx.Provide(NewRoutes),/fx.Provide(NewRoutes),\n  fx.Provide(New${uc_resource}Routes),/g" $fx_route_path
echo $fx_route_path "updated."

printf "\n\n*** Scaffolding Completely Successfully ***\n"