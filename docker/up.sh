#!/bin/bash
if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

# starting the server
cd ./../surgical-extraction
sudo docker build . -t surgical 
cd ./../docker

sudo docker-compose build 
sudo docker-compose run atlantgo bash