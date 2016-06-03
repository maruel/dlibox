# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Set this variable to your host to enable "make push" or use the form:
#   make HOST=mypi push
HOST ?= raspberrypi1


.PHONY: push run setup log


gofiles := $(wildcard **/*.go)
imgfiles := $(wildcard cmd/dotstar/images/*)
webfiles := $(wildcard cmd/dotstar/web/*)


# Regenerate the embedded files as needed.
cmd/dotstar/static_files_gen.go: $(imgfiles) $(webfiles) cmd/package/main.go
	go generate ./...

# Use a trick to preinstall all imported packages. 'go build' doesn't permit
# installing packages, only 'go install' or 'go test -i' can do. But 'go
# install' would install an ARM binary, which is not what we want.
#
# Luckily, 'go test -i' is super fast on second execution.
dotstar: $(gofiles)
	GOOS=linux GOARCH=arm go test -i ./cmd/dotstar
	GOOS=linux GOARCH=arm go build ./cmd/dotstar

# When an executable is running, it must be scp'ed aside then moved over.
# dotstar will exit safely when it detects its binary changed.
push: dotstar
	scp -q dotstar $(HOST):bin/dotstar2
	ssh $(HOST) "mv bin/dotstar2 bin/dotstar"


# Runs it locally as a fake display with the web server running on port 8010.
run: $(gofiles) cmd/dotstar/static_files_gen.go
	go install ./cmd/dotstar
	dotstar -fake -n 80 -port 8010


# Sets up a new raspberry pi.
setup: push
	scp setup/dotstar.service $(HOST):.
	ssh $(HOST) 'sed -i -e "s/pi/$$USER/g" dotstar.service && sudo -S cp dotstar.service /etc/systemd/system/dotstar.service && sudo systemctl daemon-reload && sudo systemctl enable dotstar.service && sudo systemctl start dotstar.service && rm dotstar.service'


log:
	ssh $(HOST) 'sudo -S journalctl -u dotstar'


# Defaults to cross building to ARM.
all: dotstar
