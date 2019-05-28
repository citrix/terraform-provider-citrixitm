TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=citrixitm
DOCKER_SOURCE_DIR=/terraform-provider-$(PKG_NAME)
DOCKER_IMAGE_NAME=$(PKG_NAME)-terraform
DOCKER_CONTAINER_NAME=$(PKG_NAME)_tf_dev_container

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test ./...

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

docker-build:
	@docker build -t $(DOCKER_IMAGE_NAME) .

# docker-start is used to resume an existing exited container named
# $(DOCKER_CONTAINER_NAME). Use docker-run to create a new container.
docker-start:
	docker/start_container.sh $(DOCKER_CONTAINER_NAME)

docker-exec-bash:
	@docker exec -it $(DOCKER_CONTAINER_NAME) /bin/bash

docker-exec-install-plugin:
	@docker exec -it $(DOCKER_CONTAINER_NAME) docker/install_plugin_binary.sh

docker-show-env:
	@echo ITM variables...
	@echo ITM_BASE_URL: $(ITM_BASE_URL)
	@echo ITM_CLIENT_ID: $(ITM_CLIENT_ID)
	@echo ITM_CLIENT_SECRET: $(ITM_CLIENT_SECRET)
	@echo ITM_HOST_MODULE_DIR: $(ITM_HOST_MODULE_DIR)

check-itm-env:
ifndef ITM_BASE_URL
	$(error ITM_BASE_URL is undefined)
endif
ifndef ITM_CLIENT_ID
	$(error ITM_CLIENT_ID is undefined)
endif
ifndef ITM_CLIENT_SECRET
	$(error ITM_CLIENT_SECRET is undefined)
endif

# docker-run is used to create a container named $(DOCKER_CONTAINER_NAME).
# Use docker-start to restart an existing exited container by that name.
docker-run: docker-show-env check-itm-env
ifndef ITM_HOST_MODULE_DIR
	@docker run -it --name $(DOCKER_CONTAINER_NAME) --env ITM_BASE_URL --env ITM_CLIENT_ID --env ITM_CLIENT_SECRET --mount type=bind,src=$(PWD),dst=$(DOCKER_SOURCE_DIR) $(DOCKER_IMAGE_NAME) /bin/bash
else
	@echo "Attaching host directory $(ITM_HOST_MODULE_DIR) as a bind mount to /terraform-module inside the container."
	@docker run -it --name $(DOCKER_CONTAINER_NAME) --env ITM_BASE_URL --env ITM_CLIENT_ID --env ITM_CLIENT_SECRET --mount type=bind,src=$(PWD),dst=$(DOCKER_SOURCE_DIR) --mount type=bind,src=$(ITM_HOST_MODULE_DIR),dst=/terraform-module $(DOCKER_IMAGE_NAME) /bin/bash
endif

.PHONY: build test testacc vet fmt fmtcheck vendor-status test-compile website website-test
