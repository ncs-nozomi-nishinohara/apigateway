SHELL := /bin/bash
BUILDX_ARCHS := linux/amd64,linux/ppc64le,linux/arm64
run:
	docker \
		buildx \
		build \
		--platform ${BUILDX_ARCHS} \
		-o type=registry \
		-t ncsnozominishinohara/apigateway:1.0.0 .