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
    --roles=*)
    ROLESPATH="${i#*=}"
    shift # past argument=value
    ;;
    --logs=*)
    LOGSPATH="${i#*=}"
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
echo "ROLESPATH=${ROLESPATH}"
echo "LOGSPATH=${LOGSPATH}"

docker stop gec2;
docker rm gec2;
docker run \
  --name gec2 \
  --mount type=bind,source="${CREDENTIALS}",target=/credentials \
  --mount type=bind,source="${DEPLOYCONTEXT}",target=/context \
  --mount type=bind,source="${SSHKEY}",target=/sshKey \
  --mount type=bind,source="${ROLESPATH}",target=/roles\
  --mount type=bind,source="${LOGSPATH}",target=/logs\
  -e EC2_REGION=${REGION} \
  gec2:1.3;

