.PHONY: test build image update log clean

test:
	go test ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o bin/admission-webhook-demo ./cmd

image: build
	docker build -t admission-webhook-demo .

update: image
	kubectl delete -f deployment/mutatingwebhook.yaml || true
	kubectl delete -f deployment/validatingwebhook.yaml || true
	kubectl get rs|grep admission-webhook-demo|awk '{print $$1}'|xargs kubectl delete rs
	kubectl apply -f deployment/mutatingwebhook.yaml
	kubectl apply -f deployment/validatingwebhook.yaml
	
log:
	kubectl get pod|grep admission-webhook-demo|awk 'NR==1{print $$1}'|xargs kubectl logs -f --tail 100

clean:
	rm -rf ./bin
