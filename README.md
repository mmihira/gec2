# Gec2

GEC2 is a lightweight tool to provision EC2 nodes on an AWS compatible Elatic Compoute service provider.
The only dependecy is docker.

Currently AmazonAWS and NeCTAR are supported. EC2 nodes are specified in a yaml file like:

```yaml
provider: "AWS"
nodes:
   - node1:
       ami: "ami-e686e12c"
       type: "m1.medium"
       placement: "ap-southeast-2a"
       attach_volume: false
       volume: "vol-02b6e03aa9a40d5f3"
       volume_mount_point: "/dev/xvdb"
       volume_mount_dir: "/data"
       keyname: "sshkey"
       security_groups:
         - "ssh"
       roles:
         - "init"
```

Furthermore roles can be assigned to nodes. Roles specific post deployment scripts or actions
which are run remotely on the nodes.

## Usage

First build this image

```bash
  git clone https://github.com/mmihra/gec2
  cd gec2
  ./docker_build.sh

```

Then run using the provided script
```bash
./scripts/run.sh -c=<AWS_CREDENTIALS_PATH> -r=<REGION> -l=<CONTEXT_FOLDER> -s=<SSH_KEY_PATH>
```

`AWS_CREDENTIALS_PATH` must be the path to you AWS credentails file which looks like:
```
[default]
aws_access_key_id = <>
aws_secret_access_key = <>
region = <>
```

`SSH_KEY_PATH` the path to the ssh key for access to the nodes. <br/>
The use of only one key is supported at this time



