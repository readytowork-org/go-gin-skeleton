#!/bin/bash

ROOT=$(pwd)

fileName=$(echo $1 | tr '[:upper:]' '[:lower:]')
titleFileName=$(echo ${fileName:0:1} | tr '[:lower:]' '[:upper:]')${fileName:1}


file_create=("controllers" "services" "repository" "routes" "models" )

# if files already exists then terminate the process 
for str in ${file_create[@]}; do
    FILE="${ROOT}/api/${str}/${fileName}.go"
    if test -f "$FILE"; then
        echo "${str}/${fileName} exists."
        exit
    fi
done

#  create with a template 
for str in ${file_create[@]}; do
    if [ "$str" == "models" ]; then
            sed 's/user/$fileName/g; s/User/$titleFileName/g' "${ROOT}/automate/automate-templates/${str}.txt" > temp.txt 
            mv temp.txt "${ROOT}/${str}/${fileName}.go" 
    else
        sed 's/user/$fileName/g; s/User/$titleFileName/g' "${ROOT}/automate/automate-templates/${str}.txt" > temp.txt 
        mv temp.txt "${ROOT}/api/${str}/${fileName}.go" 
    fi
                
done

echo "file created successfully"