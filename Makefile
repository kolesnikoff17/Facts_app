COV_DIR := coverage
COV_FILE := cover.out
CONTAINER_NAME := fact_app_test


run:
	docker-compose up --force-recreate --build -d

down:
	docker-compose down

$(COV_DIR):
	mkdir -p $(COV_DIR)

test: $(COV_DIR)
	docker build -f test/Dockerfile -t $(CONTAINER_NAME) .
	docker run -v ${PWD}/$(COV_DIR):/testdir/$(COV_DIR) --rm $(CONTAINER_NAME)
	docker image rm $(CONTAINER_NAME)

test_ci:
	mkdir -p $(COV_DIR)
	touch $(COV_DIR)/$(COV_FILE)
	docker build -f test/Dockerfile -t $(CONTAINER_NAME) .
	docker run --rm $(CONTAINER_NAME)

$(COV_DIR)/$(COV_FILE):
	make test

report: $(COV_DIR)/$(COV_FILE)
	go tool cover -html=$(COV_DIR)/$(COV_FILE)

clean:
	rm -rf $(COV_DIR)

.PHONY: test