# Use a trick to preinstall all imported packages. 'go build' doesn't permit
# installing packages, only 'go install' or 'go test -i' can do. But 'go
# install' would install an ARM binary, which is not what we want.
#
# Luckily, 'go test -i' is super fast on second execution.
dotstar: *.go cmd/dotstar/*.go
	GOOS=linux GOARCH=arm go test -i ./cmd/dotstar
	GOOS=linux GOARCH=arm go build ./cmd/dotstar

# When an executable is running, it must be scp'ed aside then moved over.
push: dotstar
	scp -q dotstar raspberrypi1:dotstar2
	ssh raspberrypi1 "mv dotstar2 dotstar"

all: dotstar
