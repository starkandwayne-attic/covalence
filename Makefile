BUILD_TYPE?=build

covalence:
	go $(BUILD_TYPE) ./cmd/covalence-api
	go $(BUILD_TYPE) ./cmd/covalence-agent


ARTIFACTS := artifacts/covalence-{{.OS}}-{{.Arch}}
LDFLAGS := -X main.Version=$(VERSION)
release:
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"

	@echo "Checking that TARGETS was defined in the calling environment"
	@test -n "$(TARGETS)"
	@echo "OK.  TARGETS='$(TARGETS)'"
	rm -rf artifacts
	gox -osarch="$(TARGETS)" -ldflags="$(LDFLAGS)" --output="$(ARTIFACTS)/covalence-api"      ./cmd/covalence-api
	gox -osarch="$(TARGETS)" -ldflags="$(LDFLAGS)" --output="$(ARTIFACTS)/covalence-agent"    ./cmd/covalence-agent

	cd artifacts && for x in covalence-*; do cp -a ../ui/ $$x/ui; tar -czvf $$x.tar.gz $$x; rm -r $$x;  done
