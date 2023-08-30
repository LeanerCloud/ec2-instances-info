BINDATA_FILE := data/generated_bindata.go

# upstream data
INSTANCES_URL := "https://instances.vantage.sh/instances.json"
RDS_INSTANCES_URL := "https://instances.vantage.sh/rds/instances.json"

DEPS := "curl git jq"

all: update-data run-example
.PHONY: all

check_deps:                                 ## Verify the system has all dependencies installed
	@for DEP in $(shell echo "$(DEPS)"); do \
		command -v "$$DEP" > /dev/null 2>&1 \
		|| (echo "Error: dependency '$$DEP' is absent" ; exit 1); \
	done
	@echo "all dependencies satisifed: $(DEPS)"
.PHONY: check_deps

data/instances.json:
	@mkdir -p data
	@curl $(INSTANCES_URL) -o data/instances.json
	@curl $(RDS_INSTANCES_URL) -o data/rds-instances.json


run-example:
	@go get ./...
	@go run examples/instances/instances.go | sort | tee generated_instances_data.txt | less -S

clean:
	@rm -rf data
.PHONY: clean

update-data: clean data/instances.json
.PHONY: update-data

update-data-from-local-file: all
.PHONY: update-data-from-local-file


test:
	@go test
.PHONY: test
