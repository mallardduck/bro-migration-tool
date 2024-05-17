# BRO Migration Dev Tool

**Disclaimer:** This tool is a dev tool meant for engineers and comes with no warranty or support.
This is not an official Rancher tool for end-users. Only use this tool under direct direction from SUSE/Rancher support.

## What does it do?
This tool essentially just extracts the `"clusters.management.cattle.io#v3/local.json"` file from the tar.
Then you can edit it as needed and then repackage the file into a backup.

1. Extract local cluster: `bro-migration-tool pull-local -f {tar.gz backup}`
2. Edit the `local.yaml` file as needed.
3. Repack the `local.yaml` file:
   `bro-migration-tool push-local -f {tar.gz backup} -o k3s-test1`
4. Restore the test backup you just modified.

The above example will output a new archive at `k3s-test1.tar.gz`.

### Why?
To assist with a very specific workflow being used while investigating potential improvements to BRO.

## "Automated Transformation"
There are new subcommands intended to be used as a stop gap solution for Rancher migrations between k3s<->rke2. 
One sub-command has been created for each migration direction you could use.
**Only use this method when directed by a Support Admin.**

Essentially the process to use this method will require you having local access to the backup.
We'll run our tool, then upload that backup somewhere the cluster can restore from it.  

Ultimately, this process is essentially the same as
[Migrating Rancher to a New Cluster](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/migrate-rancher-to-new-cluster) directions.
Just with this small middle step to prepare a backup for restore on the new cluster.

### K3s -> RKE2
1. Create a backup of your Rancher app running on k3s (ideally save on S3),
2. Download the backup from S3 to local system,
3. Run `bro-migration-tool k3s-rke2 -f {input.tar.gz} -o {output-file}` (note `-o` flag is just name string, w/o extension),
4. Reupload the new file created in step 3 to S3,
5. Setup a k8s cluster using equivalent RKE2 version as original k3s cluster,
6. Install Rancher Backup helm chart,
7. Apply a `Restore` manifest using the file we just uploaded,
8. Give the restore time to complete, then continue following the "Migrating Rancher to a New Cluster" directions.

### RKE2 -> K3s
1. Create a backup of your Rancher app running on k3s (ideally save on S3),
2. Download the backup from S3 to local system,
3. Run `bro-migration-tool rke2-k3s -f {input.tar.gz} -o {output-file}` (note `-o` flag is just name string, w/o extension),
4. Reupload the new file created in step 3 to S3,
5. Setup a k8s cluster using equivalent k3s version as original rke2 cluster,
6. Install Rancher Backup helm chart,
7. Apply a `Restore` manifest using the file we just uploaded,
8. Give the restore time to complete, then continue following the "Migrating Rancher to a New Cluster" directions.
