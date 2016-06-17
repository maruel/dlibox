# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Set this variable to your host to enable "make push" or use the form:
#   make HOST=raspberrypi push
HOST ?= dlibox


.PHONY: log push run setup setup_internal


gofiles := $(wildcard **/*.go)
imgfiles := $(wildcard cmd/dlibox/images/*)
webfiles := $(wildcard cmd/dlibox/web/*)


# Regenerate the embedded files as needed.
cmd/dlibox/static_files_gen.go: $(imgfiles) $(webfiles) cmd/package/main.go
	go generate ./...

# Use a trick to preinstall all imported packages. 'go build' doesn't permit
# installing packages, only 'go install' or 'go test -i' can do. But 'go
# install' would install an ARM binary, which is not what we want.
#
# Luckily, 'go test -i' is super fast on second execution.
dlibox: $(gofiles) cmd/dlibox/static_files_gen.go
	GOOS=linux GOARCH=arm go test -i ./cmd/dlibox
	GOOS=linux GOARCH=arm go build ./cmd/dlibox

# When an executable is running, it must be scp'ed aside then moved over.
# dlibox will exit safely when it detects its binary changed.
push: dlibox
	scp -q dlibox $(HOST):bin/dlibox2
	ssh $(HOST) "mv bin/dlibox2 bin/dlibox"


# Runs it locally as a fake display with the web server running on port 8010.
run: $(gofiles) cmd/dlibox/static_files_gen.go
	go install ./cmd/dlibox
	dlibox -fake


# Sets up a new raspberry pi.
setup: setup_internal push


setup_internal:
	ssh $(HOST) 'mkdir -p bin'
	scp setup/dlibox.service $(HOST):.
	# Replace 'pi' in dlibox.server by the actual remote user. To change the
	# default user, modify your local .ssh/config.
	ssh $(HOST) 'sed -i -e "s/pi/$$USER/g" dlibox.service && \
	    sudo -S cp dlibox.service /etc/systemd/system/dlibox.service && \
	    sudo systemctl daemon-reload && \
	    sudo systemctl enable dlibox.service && \
	    sudo systemctl start dlibox.service && \
	    rm dlibox.service'


log:
	ssh $(HOST) 'sudo -S journalctl -u dlibox'


# Defaults to cross building to ARM.
all: dlibox
