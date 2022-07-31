appname := smarthome-hw
version := 0.1.0
sources := $(wildcard *.go)

build = mkdir -p smarthome-hw-bin && cp -r dist/* smarthome-hw-bin && GOOS=$(1) GOARCH=$(2) go build -o ./smarthome-hw-bin/$(appname)$(3) $(4)
tar = mkdir -p build && tar -cvzf ./$(appname)_v$(version)_$(1)_$(2).tar.gz smarthome-hw-bin && mv $(appname)_v$(version)_$(1)_$(2).tar.gz build

.PHONY: all linux


all:	linux

# Update the current version in all locations
version:
	python3 update_version.py
	cd web && npm i


# Prepares everything for a version-release
# In order to publish the release to Github, run `make gh-release
release: cleanall build

# Publishes the local release to Github releases
gh-release:
	gh release create v$(version) ./build/*.tar.gz -F ./CHANGELOG.md -t 'Node v$(version)'

run:
	go run .

clean:
	rm -rf smarthome-hw-bin
	rm -rf *.log

cleanall: clean
	rm -rf build
	rm -rf config.json

# Builds
build: all linux clean

# Build architectures, in this case only amd64 for local testing and arm for the Raspberry Pi
linux: build/linux_arm.tar.gz build/linux_amd64.tar.gz

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64, -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build')
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)
