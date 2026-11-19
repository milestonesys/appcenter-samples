include ../../build/common.mak

# The microfrontend must be built before the container image, which copies client/dist.
CLIENT_DIR := containers/restful-config-api/client

build-client:
	cd $(CLIENT_DIR) && npm ci && npm run build:federation

$(addprefix containers/,$(IMAGES)): build-client
