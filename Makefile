.PHONY: help build-docker build-local clean docs run-docker

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Display this help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build-docker: clean
	wsc compile
	cp -r src/frontend/ui/node-modules/monaco-editor src/frontend/ui-dist/static/js/monaco-editor
	cp -n data/service-manager/config.json.bak data/service-manager/config.json || true
	cp -n data/service-manager/secrets.json.bak data/service-manager/secrets.json || true
	docker build -t api-server -f src/api-server/Dockerfile .
	docker build -t asset-manager -f src/asset-manager/Dockerfile .
	docker build -t auth-manager -f src/auth-manager/Dockerfile .
	docker build -t frontend -f src/frontend/Dockerfile .
	docker build -t model-manager -f src/model-manager/Dockerfile .
	docker build -t service-manager -f src/service-manager/Dockerfile .

build-local: clean  ## Build local binaries of all VelociModel services
	wsc compile
	mkdir dist
	for service in api-server asset-manager auth-manager frontend model-manager service-manager; do \
		cd src/$$service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o $$service ; \
		mkdir ../../dist/$$service ; \
		mv $$service ../../dist/$$service/$$service ; \
		cp launch.sh ../../dist/$$service/launch-$$service.sh ; \
		chmod +x ../../dist/$$service/launch-$$service.sh ; \
		cd ../.. ; \
	done
	cp -n data/service-manager/config.json.bak data/service-manager/config.json || true
	cp -n data/service-manager/secrets.json.bak data/service-manager/secrets.json || true
	cp -r src/frontend/ui-dist/templates dist/templates
	cp -r src/frontend/ui-dist/static dist/static
	cp -r src/frontend/ui/node-modules/monaco-editor dist/static/js/monaco-editor
	cp -r data/service-manager dist/service-manager/data
	cp -r data/auth-manager dist/auth-manager/data

bundle: build-local
	for service in api-server asset-manager auth-manager frontend model-manager service-manager; do \
		version=$$(cat VERSION) ; \
		mv dist/$$service dist/$$service-$$(cat VERSION) ; \
	done
	mkdir -p release
	for service in api-server asset-manager auth-manager frontend model-manager service-manager; do \
		version=$$(cat VERSION) ; \
		cd dist/$$service-$$version ; \
		tar -czvf $$service-$$version.tar.gz * ; \
		cp $$service-$$version.tar.gz ../../release ; \
		cd ../.. ; \
	done
	
publish:
	wsc compile
	cp -r src/frontend/ui/node-modules/monaco-editor src/frontend/ui-dist/static/js/monaco-editor
	cp -n data/service-manager/config.json.bak data/service-manager/config.json || true
	cp -n data/service-manager/secrets.json.bak data/service-manager/secrets.json || true
	
	for service in api-server asset-manager auth-manager frontend model-manager service-manager; do \
		version=$$(cat VERSION) ; \
		docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t jfcarter2358/velocimodel-$$service:$$version -f src/$$service/Dockerfile --push . ; \
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
