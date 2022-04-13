.PHONY: help build-docker build-local clean docs run-docker run-local test-regression test-stress test-unit

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Display this help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-local: clean  ## Build local binaries of all VelociModel services
	rm -rf dist || true
	mkdir dist
	for service in service-manager asset-manager; do \
		cd src/$$service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o $$service ; \
		mv $$service ../../dist/$$service ; \
		cp launch.sh ../../dist/launch-$$service.sh ; \
		chmod +x ../../dist/launch-$$service.sh ; \
		cd ../.. ; \
	done
	cp -r templates dist/templates
	cp -r data dist/data

clean:  ## Remove build and test artifacts
	rm -r dist || true

docs:
	cd docs && make html
