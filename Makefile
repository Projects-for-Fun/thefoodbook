tests:
	go test -count=1 -v ./...

db-up:
	docker-compose  up -d  --remove-orphans neo4j
	docker-compose  up -d migrate

webservice: db-up
	docker-compose  up -d webservice

cleanup:
	docker-compose down --remove-orphans

logs:
	docker-compose logs

integration-tests:
	API_URL := http://localhost:3001
	# INTEGRATION_TEST_SUITE_PATH is used for run specific test in Golang, if it's not specified
	# go test -tags=integration_tests ./test/integration_tests -count=1 -v -run=$(INTEGRATION_TEST_SUITE_PATH)
	go test -tags=integration_tests ./test/integration_tests -count=1 -v


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