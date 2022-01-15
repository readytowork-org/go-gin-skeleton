#!/bin/bash -e

RED='\033[0;31m'
GREEN='\033[0;32m'
LGREEN='\033[1;32m'
BROWN='\033[0;33m'
LGRAY='\033[0;37m'
NC='\033[0m' # No Color

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



echo -e "\n${LGREEN} *** Go Gin GORM Scaffold Generator *** \n"
echo -e "This scaffolder assumes that you are using RTW clean-gin template.\n ${GREEN}"
echo -e "${NC}Enter project name (eg: ecommerce-api): ${GREEN}"; read project_name
echo -e "${NC}Enter resource name(eg: ProductCategory): ${GREEN}"; read uc_resource
echo -e "${NC}Enter resource table name(eg: product_category): ${GREEN}"; read resource_table
echo -e "${NC}Enter plural resource name(eg: ProductCategories): ${GREEN}"; read plural_resource

lc_resource=$(first_lower $uc_resource)
plc_resource=$(first_lower $plural_resource)

echo -e "\n${BROWN}Generating Scaffold for ${uc_resource}${NC}:\n"

placeholder_value_hash=(
  "{{ucresource}}:$uc_resource"
  "{{plcresource}}:$plc_resource"
  "{{lcresource}}:$lc_resource"
  "{{projectname}}:$project_name"
  "{{resourcetable}}:$resource_table"
)
entity_path_hash=(
  "models:../models"
  "routes:../api/routes"
  "controllers:../api/controllers"
  "services:../api/services"
  "repository:../api/repository"
)
for entity in "${entity_path_hash[@]}"; do
  entity_name="${entity%%:*}"
  entity_path="${entity##*:}"
  file_to_write="$entity_path/${resource_table}.go"

  cat "./templates/${entity_name}-template.go" > $file_to_write
  for item in "${placeholder_value_hash[@]}"; do
    placeholder="${item%%:*}"
    value="${item##*:}"
    $resed "s/$placeholder/$value/g" $file_to_write
  done
  echo -e ${GREEN} $file_to_write "created." ${NC}
done

# inject fx deps
fx_path_hash=(
  "Controller:../api/controllers/controllers.go"
  "Service:../api/services/services.go"
  "Repository:../api/repository/repository.go"
)
fx_init_string="var Module = fx.Options("
echo -e "\n${BROWN}Injecting Dependencies:\n"
for deps_value in "${fx_path_hash[@]}"; do
  deps_name="${deps_value%%:*}"
  deps_path="${deps_value##*:}"
  $resed "s/${fx_init_string}/${fx_init_string}\n\tfx.Provide(New${uc_resource}${deps_name}),/g" $deps_path
  echo -e ${BROWN} $deps_path "updated." ${NC}
done

# fx routes
fx_route_path="../api/routes/routes.go"
$resed "s/func NewRoutes(/func NewRoutes(\n\t${lc_resource}Routes ${uc_resource}Routes,/g" $fx_route_path
$resed "s/return Routes{/return Routes{\n\t\t${lc_resource}Routes,/g" $fx_route_path
$resed "s/fx.Provide(NewRoutes),/fx.Provide(NewRoutes),\n\tfx.Provide(New${uc_resource}Routes),/g" $fx_route_path
echo -e "${BROWN} $fx_route_path "updated." ${NC}"

echo -e "\n${GREEN}To generate migrations run: ./migration.sh${NC}"

echo -e "\n${LGREEN}*** Scaffolding Completely Successfully ***\n${NC}"