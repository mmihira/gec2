#!/bin/sh

for i in "$@"
do
case $i in
    -c=*|--credentials=*)
    CREDENTIALS="${i#*=}"
    shift # past argument=value
    ;;
    -r=*|--region=*)
    REGION="${i#*=}"
    shift # past argument=value
    ;;
    -l=*|--deploycontext=*)
    DEPLOYCONTEXT="${i#*=}"
    shift # past argument=value
    ;;
    -s=*|--sshkey=*)
    SSHKEY="${i#*=}"
    shift # past argument=value
    ;;
    *)
          # unknown option
    ;;
esac
done

echo "CREDENTIALS  = ${CREDENTIALS}"
echo "REGION = ${REGION}"
echo "DEPLOYCONTEXT= ${DEPLOYCONTEXT}"
echo "SSHKEY= ${SSHKEY}"

docker stop gec2;
docker rm gec2;
docker run \
  --name gec2 \
  --mount type=bind,source="${CREDENTIALS}",target=/credentials \
  --mount type=bind,source="${DEPLOYCONTEXT}",target=/context \
  --mount type=bind,source="${SSHKEY}",target=/sshKey \
  -e EC2_REGION=${REGION} \
  gec2:1.0;

