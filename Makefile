all: install

# dependencies installs all of the dependencies that are required for
# building the block explorer
dependencies:
	go install -race std
	go get -u github.com/agl/ed25519
	go get -u github.com/dchest/blake2b
	go get -u golang.org/x/crypto/twofish
	go get -u github.com/NebulousLabs/merkletree
	go get -u github.com/NebulousLabs/Sia/types
	go get -u github.com/NebulousLabs/Sia/modules

# Fmt formats all packages properly
fmt:
	go fmt ./...

install: fmt
	go install ./...
