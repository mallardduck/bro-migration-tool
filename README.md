# BRO Migration Dev Tool

**CAUTION WILL EAT YOUR CAT**

**Disclaimer:** This is not an official Rancher tool, nor is it supported by SUSE/Rancher support.
This tool is a dev tool meant for engineers and comes with no warranty or support.

## What it do?
This tool essentially just extracts the `"clusters.management.cattle.io#v3/local.json"` file from the tar.
Then you can edit it as needed and then repackage the file into a backup.

1. Extract local cluster: `bro-migration-tool pull-local -f {tar.gz backup}`
2. Edit the `local.yaml` file as needed.
3. Repack the `local.yaml` file:
   `bro-migration-tool push-local -f {tar.gz backup} -o k3s-test1`
4. Restore the test backup you just modified.

The above example will output a new archive at `k3s-test1.tar.gz`.