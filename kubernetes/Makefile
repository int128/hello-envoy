deploy:
	kubectl create cm envoy --from-file=../envoy_dynamic
	kubectl apply -f deployment.yaml

clean:
	kubectl delete deployment/envoy configmap/envoy
