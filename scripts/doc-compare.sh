#!/usr/bin/env bash

# Compare docs
# Run this when doc updating is finished, this will save patch files for next doc updating
SOURCEDIR="/home/zhangjishu/go/src/github.com/g42cloud-terraform/terraform-provider-g42cloud"
DESTINATIONDIR="/home/zhangjishu/go/src/github.com/huaweicloud/terraform-provider-huaweicloud"

echo "==> Comparing docs..."
data_source_filenames=$(ls ${SOURCEDIR}/docs/data-sources/*.md)
for file in ${data_source_filenames};do
    echo ${file}
    echo ${DESTINATIONDIR}/docs/data-sources/${file##*/}
    diff ${DESTINATIONDIR}/docs/data-sources/${file##*/} ${file} > `p=${SOURCEDIR}/patchfiles/data-sources;[[ ! -d "${p}" ]] && mkdir -p ${p};echo ${p}/${file##*/}.patch`
done

resource_filenames=$(ls ${SOURCEDIR}/docs/resources/*.md)
for file in ${resource_filenames};do
    echo ${file}
    echo ${DESTINATIONDIR}/docs/resources/${file##*/}
    diff ${DESTINATIONDIR}/docs/resources/${file##*/} ${file} > `p=${SOURCEDIR}/patchfiles/resources;[[ ! -d "${p}" ]] && mkdir -p ${p};echo ${p}/${file##*/}.patch`
done

exit 0
