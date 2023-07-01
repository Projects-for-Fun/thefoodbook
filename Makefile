.PHONY: deps-up deps-down webservice cleanup logs run-tests lint

deps-up:
	docker-compose -f docker-compose.dependencies.yml up -d neo4j
	docker-compose -f docker-compose.dependencies.yml up -d migrate

deps-down:
	docker-compose -f docker-compose.dependencies.yml down --remove-orphans

webservice: deps-up
	docker-compose -f docker-compose.service.yml up -d webservice --build -d --force-recreate

cleanup:
	docker-compose down --remove-orphans

logs:
	docker-compose -f docker-compose.dependencies.yml logs
	docker-compose -f docker-compose.service.yml logs
	exit 1

run-test:
	go test -count=1 -cover -v ./...

run-integration: deps-up
	docker-compose -f docker-compose.service.yml up -d integration_tests --build -d --force-recreate

lint:
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run \
		-it \
		--rm \
		-v ~/go:/root/go \
		-v golangci-lint-cache:/root/.cache \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:latest golangci-lint run --deadline=65s $(OUTPUT_OPTIONS) -v
		-c ./.golangci.yml