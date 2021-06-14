deploy:
	kubectl create cm envoy --from-file=configmap
	kubectl apply -f deployment.yaml

clean:
	kubectl delete deployment/envoy configmap/envoy
