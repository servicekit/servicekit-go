_ci: lint vet test

lint:
	@echo "lint"
	@echo "-------------------"
	@golint $$(go list ./...|grep -v vendor)

vet:
	@echo "vet"
	@echo "-------------------"
	@go vet $$(go list ./...|grep -v vendor)

test:
	@echo "test"
	@echo "-------------------"
	@go test $$(go list ./...|grep -v vendor)

ci:
	@docker build ${docker_build_args} -t servicekit-go-make .
	@docker run servicekit-go-make _ci

travis: _ci

fmt:
	@go fmt $$(go list ./...|grep -v vendor)

govendor:
	@echo "install dependencies..."
	@govendor init
	@govendor fetch +o -v
