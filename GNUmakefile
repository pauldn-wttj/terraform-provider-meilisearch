default: build

build:
	go install .

testacc:
	echo "Starting Meilisearch and waiting until it's ready"
	docker-compose -f docker_compose/docker-compose.yml up -d
	sleep 5

	echo "Creating test resources"
	./docker_compose/test_init.sh

	echo "Running Terraform tests"
	TF_ACC=1 go test -count=1 -v ./internal/provider -timeout 120m

	echo "Stopping Meilisearch and cleaning"
	docker-compose -f docker_compose/docker-compose.yml down
	docker volume rm -f docker_compose_meili_data

dev:
	docker-compose -f docker_compose/docker-compose.yml up

clean:
	docker-compose -f docker_compose/docker-compose.yml rm
	docker volume rm -f docker_compose_meili_data
	rm -f examples/resources/terraform.tfstate
	rm -f examples/resources/terraform.tfstate.backup
	rm -f examples/data-sources/terraform.tfstate
	rm -f examples/data-sources/terraform.tfstate.backup
	rm -f examples/provider/terraform.tfstate
	rm -f examples/provider/terraform.tfstate.backup
