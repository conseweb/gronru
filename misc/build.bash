#!/bin/bash -e

# Copyright 2013 gronru authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This script is used to build components from gronru server (webserver and
# git wrapper). It's based on misc/build-server.bash from tsuru repository.

destination_dir="dist-server"

function build_and_package {
	echo "Building $1... "
 	go build -o ${destination_dir}/gronru-${1} github.com/globocom/gronru/${1}
	tar -C $destination_dir -czf $destination_dir/gronru-${REVISION}-${1}.tar.gz gronru-$1
	rm $destination_dir/gronru-$1
}

echo -n "Creating \"$destination_dir\" directory... "
mkdir -p $destination_dir
echo "ok"

build_and_package bin
build_and_package webserver
