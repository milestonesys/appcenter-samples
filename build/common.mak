
APP_DEFINITION := app-definition.yaml
CUSTOM_DEFINITION := custom-definition
APP_BUILDER := $(dir $(lastword $(MAKEFILE_LIST)))app-builder.sh

# Check if custom definition exists
ifeq ($(shell test -d $(CUSTOM_DEFINITION) && echo yes), yes)
CUSTOM_FLAG := -c $(CUSTOM_DEFINITION)
else
CUSTOM_FLAG :=
endif

default: build

~/.kube/config:
	@echo "Logging into system with sandbox environment enabled"
	@$(APP_BUILDER) login

login: ~/.kube/config

logout:
	@echo "Logging out from system"
	@$(APP_BUILDER) logout

build: build-image
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) build

dashboard: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) dashboard

list: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) list

push: login push-image
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) push

remove: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) remove

install-from-file: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) install-from-file

install-from-repo: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) install-from-repo

uninstall: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) uninstall

restart: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) restart

status: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) status

events: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) events

volumes: login
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) volumes

bundle: login build-image push-image
	@$(APP_BUILDER) -f $(APP_DEFINITION) $(CUSTOM_FLAG) bundle


# --- Docker rules

IMAGES = $(shell yq -r '.services[].containers[].image | capture("^sandbox.io/(?<image>.*):(?<tag>.*)$$") | .image+"_v"+.tag' app-definition.yaml)
SYSTEM_IP = $(shell kubectl get configmap -n kube-public application-platform-metadata -o "jsonpath={.data['system_ip']}")
DOCKER_REGISTRY = $(SYSTEM_IP):5000/sandbox.io

build-image: login $(addprefix containers/,$(IMAGES))
push-image: login $(addprefix push-image/,$(IMAGES))
remove-image: login $(addprefix remove-image/,$(IMAGES))

DOCKER_DIR = $(word 1, $(subst _v, ,$1))
DOCKER_TAG = $(word 2, $(subst _v, ,$1))
DOCKER_IMAGE = $(DOCKER_REGISTRY)/$(call DOCKER_DIR, $1):$(call DOCKER_TAG, $1)

containers/%:
	@docker build $(call DOCKER_DIR, $@) -t $(call DOCKER_IMAGE, $*)

push-image/%:
	@docker push $(call DOCKER_IMAGE, $*)

remove-image/%:
	@docker rmi $(call DOCKER_IMAGE, $*)

clear-image-cache:
	@docker system prune -a --volumes --force