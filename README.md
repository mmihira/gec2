# Gec2

GEC2 is a lightweight tool to provision EC2 nodes on an AWS compatible Elatic Compute service provider.
The default way to run Gec2 is to build the docker image and run in a docker container.

## Configuration

Currently AmazonAWS and NeCTAR are supported. EC2 nodes are specified in a yaml file like:

```yaml
provider: "AWS"
nodes:
   - node1:
       ami: "ami-e686e12c"
       type: "m1.medium"
       placement: "ap-southeast-2a"
       keyname: "sshkey"
       security_groups:
         - "ssh"
       roles:
         - "init"

roles: ["init"]
```

Furthermore roles can be assigned to nodes. Roles are post deployment instructions which are run against
nodes.

Roles are defined in a roles.yaml file and look like this:

```yaml
init:
  steps:
    - stepType: "script"
      scripts:
        - "roles/init/run.sh"
```

There are three types of steps available:<br>
- `script`
  Run a script on the server
- `copy`
  Copy a file from the docker context to the node
- `template`
  Template a file and then copy it to the node

All filenames are resolved relative to the context folder specified on startup.
When the nodes are provisioned a file called deployed_schema.json is created in the context which looks
like :

```json
{
  "node1": {
    "name": "node1",
    "keyname": "somesshkey",
    "roles": [
      "init",
    ],
    "ip": "34.xxx.xx.xxx"
  }
}
```

This data structure can be used in templating using golang template syntax:
`https://golang.org/pkg/text/template/`

See the examples folders for an example.

## Installating and running

First build this image

```bash
  git clone https://github.com/mmihra/gec2
  cd gec2
  ./docker_build.sh
```

Or pull the latest version from dockerhub

```bash
  docker pull mmihira/gec2:1.1
```

Then run using the provided script
```bash
./scripts/run.sh -c=<AWS_CREDENTIALS_PATH> -r=<REGION> -l=<CONTEXT_FOLDER> -s=<SSH_KEY_PATH> --roles=<PATH_TO_ROLES> --logs=<PATH_TO_LOGS>
```

`AWS_CREDENTIALS_PATH` must be the path to you AWS (or NeCTAR) credentails file which looks like:
```
[default]
aws_access_key_id = <>
aws_secret_access_key = <>
region = <>
```

`SSH_KEY_PATH` the path to the ssh key for access to the nodes. <br/>
The use of only one key is supported at this time

- The context foler should be a folder with at a minimum these files:
  - `config.yaml`
- The roles folder should be a folder with the roles
- The log folder is where deploy logs are written to

## Tests

Run `go test -v ./...`
