build-dev:
	docker build . --target development --tag openslides-permission-dev

run-tests:
	docker build . --target testing --tag openslides-permission-test
	docker run openslides-permission-test

cover.out:
	go test ./...  -coverprofile=cover.out -coverpkg=./internal/collection

cover: cover.out
	go tool cover -html cover.out