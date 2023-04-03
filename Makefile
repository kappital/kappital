# Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Image Version
VERSION ?= latest
# Helm Chart Version
HELM_VERSION ?= 0.0.1

# Directories
BIN=bin
PACKAGE=bin/package

# Image Names
SERVICE_ENGINE ?= kappital/kappital-engine
MANAGER ?= kappital/manager

# Image-Related Variables
CREATE_NAMESPACE ?= false # if namespace is exist, set it to the false
REGISTRY ?= ""
PULL_POLICY ?= IfNotPresent

# Liveness and Readiness deploy variables for service engine
LIVENESS_PORT ?= 8081
LIVENESS_SCHEME ?= HTTPS

# Manager Deploy Configs
MANAGER_LOG ?= /opt/kappital/manager/log
MANAGER_SERVICE_IP ?= x.x.x.x # service.clusterIP

# Common Configs
USE_HOST_NETWORK ?= false

# Git Information
GIT_VERSION ?= $(shell git describe --tags --dirty)
GIT_COMMIT_HASH ?= $(shell git rev-parse HEAD)
GIT_TREESTATE = "clean"
GIT_DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(GIT_DIFF), 1)
    GIT_TREESTATE = "dirty"
endif
BUILDDATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# go build parameters
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
SOURCES := $(shell find . -type f  -name '*.go')
LDFLAGS := "-linkmode=external -extldflags '-Wl,-z,now,-s' \
				-X github.com/kappital/kappital/pkg/utils/version.gitVersion=$(GIT_VERSION) \
				-X github.com/kappital/kappital/pkg/utils/version.gitCommit=$(GIT_COMMIT_HASH) \
				-X github.com/kappital/kappital/pkg/utils/version.gitTreeState=$(GIT_TREESTATE) \
				-X github.com/kappital/kappital/pkg/utils/version.buildDate=$(BUILDDATE)"
CGO_FLAG="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -ftrapv"
COVERAGE=$(shell pwd)/tests/coverage

all: manager kappital-engine kappctl
all-image: kappital-engine-image manager-image
deploy-all: all-image deploy-kappital-engine deploy-manager

# generate and check the go files
fmt:
	go fmt ./cmd/...
	go fmt ./pkg/...

vet:
	go vet ./cmd/...
	go vet ./pkg/...

test: fmt vet
	mkdir -p ${COVERAGE}
	go test --coverprofile=${COVERAGE}/cmd.out ./cmd/...
	go test --coverprofile=${COVERAGE}/pkg.out ./pkg/...
	go tool cover -func=${COVERAGE}/cmd.out -o ${COVERAGE}/cmd.txt
	go tool cover -func=${COVERAGE}/pkg.out -o ${COVERAGE}/pkg.txt
	cat ${COVERAGE}/cmd.txt
	cat ${COVERAGE}/pkg.txt

# construct the binary files to the 'bin' directory
kappital-engine: fmt vet
	CGO_ENABLE=0 CGO_CFLAGS=${CGO_FLAG} go build -buildmode=pie -ldflags=${LDFLAGS} -o bin/kappital-engine cmd/engine/main.go

manager: fmt vet
	CGO_ENABLE=0 CGO_CFLAGS=${CGO_FLAG} go build -buildmode=pie -ldflags=${LDFLAGS} -o bin/kappital-manager cmd/manager/main.go


kappctl: fmt vet
	CGO_ENABLE=0 CGO_CFLAGS=${CGO_FLAG} go build -buildmode=pie -ldflags=${LDFLAGS} -o bin/kappctl cmd/kappctl/main.go

# construct the docker images from the docker file
kappital-engine-image: kappital-engine
	docker build -t ${SERVICE_ENGINE}:${VERSION} -f build/kappital-engine/Dockerfile .

manager-image: manager
	docker build -t ${MANAGER}:${VERSION} -f build/kappital-manager/Dockerfile .


# construct the helm chart packages to the 'bin' directory
kappital-engine-package:
	helm package charts/kappital-engine --version=${HELM_VERSION}
	mkdir -p bin/package
	mv kappital-engine-${HELM_VERSION}.tgz bin/package/

manager-package:
	helm package charts/kappital-manager --version=${HELM_VERSION}
	mkdir -p bin/package
	mv kappital-manager-${HELM_VERSION}.tgz bin/package/

# deploy the kappital module to the cluster
deploy-kappital-engine: kappital-engine-package
	# all kappital services will install in the namespace kappital-system
	helm install kappital-engine \
        --create-namespace=${CREATE_NAMESPACE} \
        --namespace=kappital-system \
		--set image.name=${SERVICE_ENGINE} \
		--set image.tag=${VERSION} \
		--set image.registry=${REGISTRY} \
		--set image.pullPolicy=${PULL_POLICY} \
		--set env.runningMode=${RUNNING_MODE} \
		--set livenessReadiness.port=${LIVENESS_PORT} \
		--set livenessReadiness.scheme=${LIVENESS_SCHEME} \
		bin/package/kappital-engine-${HELM_VERSION}.tgz

deploy-manager: manager-package
	# all kappital services will install in the namespace kappital-system
	helm install kappital-manager \
		--create-namespace=${CREATE_NAMESPACE} \
		--namespace=kappital-system \
		--set manager.image.name=${MANAGER} \
		--set manager.image.tag=${VERSION} \
		--set manager.image.registry=${REGISTRY} \
		--set manager.image.pullPolicy=${PULL_POLICY} \
		--set manager.logDir=${MANAGER_LOG} \
		--set manager.hostNetwork=${USE_HOST_NETWORK} \
		--set manager.service.clusterIP=${MANAGER_SERVICE_IP} \
		bin/package/kappital-manager-${HELM_VERSION}.tgz

clean:
	rm -rf bin/kappital-engine
	rm -rf bin/kappital-manager
	rm -rf bin/package
