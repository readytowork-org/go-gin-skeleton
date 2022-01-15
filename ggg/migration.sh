#!/bin/bash -e
cd ..
echo "Enter resource table name(eg: product_category):"; read resource_table
migration_name="create_${resource_table}_table"
make auto-create GGG_NAME=$migration_name

migration_up_filepath=`find . -name "*${migration_name}.up.sql"`
migration_down_filepath=`find . -name "*${migration_name}.down.sql"`

sed "s/{{resourcetable}}/$resource_table/g" ./ggg/templates/up-migration-template.go > "${migration_up_filepath}"
echo "DROP TABLE IF EXISTS ${resource_table};" > $migration_down_filepath
cd ggg/