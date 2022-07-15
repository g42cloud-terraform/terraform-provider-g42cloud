!/usr/bin/env bash

Update docs
Run this to updating the files, you must update you local huaweicloud repo to latest version before run this
SOURCEDIR="/home/zhangjishu/go/src/github.com/g42cloud-terraform/terraform-provider-g42cloud"
DESTINATIONDIR="/home/zhangjishu/go/src/github.com/huaweicloud/terraform-provider-huaweicloud"

echo "==> Updating docs..."
data_source_filenames=$(ls ${SOURCEDIR}/docs/data-sources/*.md)
for file in ${data_source_filenames};do
    echo ${file}
    echo ${DESTINATIONDIR}/docs/data-sources/${file##*/}
    cp ${DESTINATIONDIR}/docs/data-sources/${file##*/} ${file}
    patch ${file} ${SOURCEDIR}/patchfiles/data-sources/${file##*/}.patch
done

resource_filenames=$(ls ${SOURCEDIR}/docs/resources/*.md)
for file in ${resource_filenames};do
    echo ${file}
    echo ${DESTINATIONDIR}/docs/resources/${file##*/}
    cp ${DESTINATIONDIR}/docs/resources/${file##*/} ${file}
    patch ${file} ${SOURCEDIR}/patchfiles/resources/${file##*/}.patch
done

exit 0
