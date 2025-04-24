BINDATA_FILE := data/generated_bindata.go

# upstream data
EC2_INSTANCES_URL := "https://instances.vantage.sh/instances.json"
RDS_INSTANCES_URL := "https://instances.vantage.sh/rds/instances.json"
ELASTICACHE_INSTANCES_URL := "https://instances.vantage.sh/cache/instances.json"
OPENSEARCH_INSTANCES_URL := "https://instances.vantage.sh/opensearch/instances.json"
AZURE_INSTANCES_URL := "https://instances.vantage.sh/azure/instances.json"

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
	@curl -s --compressed $(EC2_INSTANCES_URL) -o data/instances.json
	@curl -s --compressed $(RDS_INSTANCES_URL) -o data/rds-instances.json
	@curl -s --compressed $(ELASTICACHE_INSTANCES_URL) -o data/elasticache-instances.json
	@curl -s --compressed $(OPENSEARCH_INSTANCES_URL) -o data/opensearch-instances.json
	@curl -s --compressed $(AZURE_INSTANCES_URL) -o data/azure-instances.json

run-example:
	@go get ./...
	@go run examples/instances/instances.go | sort | tee generated_instances_data.txt | less -S

run-azure-example:
	@go get ./...
	@go run examples/azure/azure-vm.go | tee generated_azure_instances_data.txt | less -S

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

azure: data/instances.json run-azure-example
.PHONY: azure