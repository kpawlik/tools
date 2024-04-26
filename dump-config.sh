#!/bin/bash


items=(features enums layers layer_groups networks datasources applications roles users groups table_sets settings)
ds=`date "+%d%m%Y_%H%M%S"`
BASE=/shared-data/config-backup/$ds

for item in "${items[@]}"
do
    dump_dir="${BASE}/${item}"
    mkdir -p "$dump_dir"
    /opt/iqgeo/platform/Tools/myw_db $PGDATABASE dump $dump_dir $item
done
