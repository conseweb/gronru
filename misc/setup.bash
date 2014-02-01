#!/bin/bash -e

# copy gronru command to a place in path
go build -o gronru bin/gronru.go
sudo mv gronru /usr/local/bin/

# copy default config file
sudo cp etc/gronru.conf /etc/

# starts gronru api web server
go build -o gronru-webserver webserver/main.go
./gronru-webserver > $HOME/gronru-webserver.out 2>&1 &
git daemon --base-path=/var/repositories --detach --export-all
