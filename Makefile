.PHONY: infra-up infra-down webservice cleanup logs run-tests lint

infra-up:
	docker-compose -f docker-compose.local.yml up -d neo4j
	docker-compose -f docker-compose.local.yml up -d migrate

infra-down:
	docker-compose -f docker-compose.local.yml down neo4j --remove-orphans
	docker-compose -f docker-compose.local.yml down migrate --remove-orphans

webservice: infra-up
	docker-compose -f docker-compose.local.yml up -d webservice --build -d --force-recreate

cleanup:
	docker-compose -f docker-compose.local.yml down --remove-orphans

logs:
	docker-compose -f docker-compose.local.yml logs
	exit 1

run-test:
	go test -count=1 -cover -v ./...

run-integration: create-integration-tests cleanup

create-integration-tests:
	docker-compose -f docker-compose.local.yml up --build -d neo4j --force-recreate
	docker-compose -f docker-compose.local.yml up --build -d migrate --force-recreate
	docker-compose -f docker-compose.local.yml up integration_tests --build --abort-on-container-exit --exit-code-from=integration_tests --force-recreate


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