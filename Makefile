.PHONY: help build-docker build-local clean docs run-docker

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Display this help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-docker: clean
	wsc compile
	cp -r src/frontend/ui/node-modules/monaco-editor src/frontend/ui-dist/static/js/monaco-editor
	cp -n data/config.json.bak data/config.json || true
	cp -n data/secrets.json.bak data/secrets.json || true
	docker build -t api-server -f src/api-server/Dockerfile .
	docker build -t asset-manager -f src/asset-manager/Dockerfile .
	docker build -t frontend -f src/frontend/Dockerfile .
	docker build -t model-manager -f src/model-manager/Dockerfile .
	docker build -t service-manager -f src/service-manager/Dockerfile .

build-local: clean  ## Build local binaries of all VelociModel services
	wsc compile
	mkdir dist
	for service in api-server asset-manager frontend model-manager service-manager; do \
		cd src/$$service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o $$service ; \
		mv $$service ../../dist/$$service ; \
		cp launch.sh ../../dist/launch-$$service.sh ; \
		chmod +x ../../dist/launch-$$service.sh ; \
		cd ../.. ; \
	done
	cp -n data/config.json.bak data/config.json || true
	cp -n data/secrets.json.bak data/secrets.json || true
	cp -r src/frontend/ui-dist/templates dist/templates
	cp -r src/frontend/ui-dist/static dist/static
	cp -r src/frontend/ui/node-modules/monaco-editor dist/static/js/monaco-editor
	cp -r data dist/data

bundle: build-local
	for service in api-server asset-manager frontend model-manager service-manager; do \
		mkdir dist/$$service-$$(cat VERSION) ; \
		mv dist/$$service dist/$$service-$$(cat VERSION) ; \
		mv dist/launch-$$service.sh dist/$$service-$$(cat VERSION) ; \
	done
	mkdir dist/api-server-$$(cat VERSION)/templates
	mkdir dist/api-server-$$(cat VERSION)/static
	mkdir dist/service-manager-$$(cat VERSION)/data
	mv dist/templates dist/api-server-$$(cat VERSION)
	mv dist/static dist/api-server-$$(cat VERSION)
	mv dist/data dist/service-manager-$$(cat VERSION)
	mkdir -p release
	for service in api-server asset-manager frontend model-manager service-manager; do \
		version=$$(cat VERSION) ; \
		cd dist/$$service-$$version ; \
		tar -czvf $$service-$$version.tar.gz * ; \
		cp $$service-$$version.tar.gz ../../release ; \
		cd ../.. ; \
	done
	
	

clean:  ## Remove build and test artifacts
	rm -r dist || true
	rm -r release || true
	rm -r src/frontend/ui-dist || true

docs:
	cd docs && make html

run-docker:
	docker-compose rm -f
	docker-compose up
