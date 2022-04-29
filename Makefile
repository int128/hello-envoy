IMAGE := envoyproxy/envoy:v1.22-latest

all:

validate: validate_static validate_dynamic

validate_static: envoy_static.yaml
	docker run --rm -v $(CURDIR)/envoy_static.yaml:/envoy_static.yaml $(IMAGE) -c /envoy_static.yaml --mode validate

validate_dynamic: envoy_dynamic
	docker run --rm -v $(CURDIR)/envoy_dynamic:/etc/envoy $(IMAGE) -c /etc/envoy/bootstrap.yaml --mode validate
