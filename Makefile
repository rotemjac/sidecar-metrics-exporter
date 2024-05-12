# import deploy config
# You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

build2:
	go mod tidy
	go mod vendor
	cd cmd && CGO_ENABLED=0 go build -o ../artifacts/${BINARY_NAME}

image: build2
	docker build -f build/Dockerfile -t jmx-sidecar-exporter:${TAG} .

push3: image
	kind load docker-image jmx-sidecar-exporter:${TAG} --name main
