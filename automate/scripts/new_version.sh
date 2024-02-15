set -e

os_name=$(uname)
ROOT=$(pwd)
read -r _ project_name _ <go.mod
project_name=$(echo $project_name | tr -d '\r')
config_directory="${ROOT}/apps"

default_version="v0"

get_current_version() {

# Check if the directory exists
if [ -d "$config_directory" ]; then
    # Get list of folders into a string
    folders_string=$(ls -d "$config_directory"/v*/ 2>/dev/null)

    # Extract version numbers
    versions=()
    for folder in $folders_string; do
        version=$(basename "$folder" | sed 's/v//')
        versions+=("$version")
    done

    # Find the greatest version number
    current_version=$(printf "%s\n" "${versions[@]}" | sort -n | tail -n 1)
    # Print the greatest version
    echo "v${current_version}"
    
else
    echo $default_version 
fi
}

current_version=$(get_current_version)

get_new_version() {
     _current_version=$(echo "$1" | sed 's/v//')
    _new_version=$(($_current_version+1))
    echo "v${_new_version}"
}

new_version=$(get_new_version $current_version)

echo  $current_version  "current version"
echo  $new_version  "new version"



# Check if the app exists
if [ -d "$config_directory/${new_version}" ]; then
  echo "$new_version already exists."
  exit
else
  mkdir ${config_directory}/${new_version}

  mkdir ${ROOT}/config/${new_version}
fi

# create config.go and router.go in Root/config


#create file
create_file() {
    #pass first parameter as entity name
    #pass second parameter as file to write
    entity_name=$1
    file_to_write=$2

    cat "${ROOT}/automate/templates/new_version/${entity_name}.txt" >>$file_to_write

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

placeholder_value_hash=(
  "{{version}}:$new_version"
)

entity_path_hash=(
  "conf:${ROOT}/config/${new_version}"
  "router:${ROOT}/config/${new_version}"
)
# create these two file
for entity in "${entity_path_hash[@]}"; do
    entity_name="${entity%%:*}"
    entity_path="${entity##*:}"
    file_to_write="$entity_path/${entity_name}.go"
    create_file $entity_name $file_to_write
done

#setup config and routers in bootstrap

bootstrap_path="${ROOT}/config/bootstrap.go"
fx_module_string="var Module = fx.Options("
func_bootstrap_string="func bootstrap("
route_setup_point="middlewares.Setup()"
route_info="${new_version}_routes.Setup())"

import_name="${project_name}/config/${new_version}"
if [[ $os_name == "Darwin" ]]; then
#import
  sed -i '' -e "/^import (/a\\
  ${new_version} \"$import_name\"
  " $bootstrap_path
  # in Module
  sed -i "" "s/${fx_module_string}/${fx_module_string}\n\t  ${new_version}.InstalledApps,/g" $bootstrap_path
  # in func bootstrap()
  sed -i "" "s/${func_bootstrap_string}/${func_bootstrap_string}\n\t  ${new_version}_routes ${new_version}.Routes,/g" $bootstrap_path
  #// setup route
  sed -i "" "s/${route_setup_point}/${route_setup_point}\n\t  ${new_version}_routes.Setup()/g" $bootstrap_path

else
  sed -i "s/${fx_module_string}/${fx_module_string}\n\t  ${new_version}.Module,/g" $bootstrap_path
  sed -i "s/${func_bootstrap_string}/${func_bootstrap_string}\n\t ${new_version}_routes ${new_version}.Routes,/g" $bootstrap_path
#// setup route
  sed -i "s/${route_setup_point}/${route_setup_point}\n\t  ${new_version}_routes.Setup()/g" $bootstrap_path

fi

