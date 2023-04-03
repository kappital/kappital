# Installing Kappital

## Suggestion

Suggest deploying all Kappital components as containers. Also, the Kappital components can deploy as a process in virtual machine.

## Prerequisites

1. The Kapptial install should have **Kubernetes cluster with v1.17+** and **Helm with v3+**. In addition, the minikube's and kind's Kubernetes may not provide some functions. 
2. The kappital-manager depends on database. It uses **SQLite as default database**.
3. Both kappital-manager offers HTTPS as default and TLS 1.2 two-way authentication option. Thus, kappital-manager needs root private key (file name is `ca.crt` as pass in), and service certificate & private key (file names are `server.crt` and `server.key` as pass in). The binary tool `kappctl` need root private key, and client certificate & private key (file names are `client.crt` and `client.key` as config). **User need to create and offer certificates and private keys for each model. OR, the easy approach is re-use the Kubernetes certificates or Kubernetes creating certificate method ([link](https://kubernetes.io/docs/tasks/administer-cluster/certificates/)).**
   1. The deployed directory structure and limitation as following:
    > ```shell
    > kappital/chart/kappital-manager
    > ├── certs
    > │   └── conf
    > │       ├── ca.crt      # root certificate; if want to use TLS 1.2 two-way authentication, must have this file
    > │       ├── server.crt  # server end certificate file
    > │       └── server.key  # server end private key file
    > ├── Chart.yaml
    > ├── README.md
    > ├── templates
    > │   ├── manager.yaml
    > │   └── configs.yaml
    > └── values.yaml
    > ```
   2. The binary tool `kappctl` config file structure like:
    > ```shell
    > /root/.kappital
    > └── config   # the config file which contain the client certificate and private key, also, include the root private key.
    > ```
    > The detail information please read **3.4**
   3. If setting kappital-manager's `.Values.checkIdentity` attribute to the true, kappital-manager will double-check the CN (CommonName) is the `Kappital - Client`. If `kappctl` bring the certificate without `Kappital - Client` CN, kappital-manager will reject the request.
   

## 1. Build All Binary Files

Assume you have downloaded the kappital directory and in the kappital project root directory. All binaries are generated via Makefile. 

```shell
[root@kappital kappital]# make all
go fmt ./cmd/...
go fmt ./pkg/...
go vet ./cmd/...
go vet ./pkg/...
CGO_ENABLE=0 CGO_CFLAGS="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -ftrapv" go build -buildmode=pie -ldflags="-linkmode=external -extldflags '-Wl,-z,now' -X github.com/kappital/kappital/pkg/utils/version.gitVersion=0.1.0-dirty -X github.com/kappital/kappital/pkg/utils/version.gitCommit=9ee0c6a2d5489e183ca9978d8e08eabc68ecfad6 -X github.com/kappital/kappital/pkg/utils/version.gitTreeState="dirty" -X github.com/kappital/kappital/pkg/utils/version.buildDate=2022-10-13T07:03:36Z" -o bin/kappital-manager cmd/manager/main.go
CGO_ENABLE=0 CGO_CFLAGS="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -ftrapv" go build -buildmode=pie -ldflags="-linkmode=external -extldflags '-Wl,-z,now' -X github.com/kappital/kappital/pkg/utils/version.gitVersion=0.1.0-dirty -X github.com/kappital/kappital/pkg/utils/version.gitCommit=9ee0c6a2d5489e183ca9978d8e08eabc68ecfad6 -X github.com/kappital/kappital/pkg/utils/version.gitTreeState="dirty" -X github.com/kappital/kappital/pkg/utils/version.buildDate=2022-10-13T07:03:36Z" -o bin/kappital-engine cmd/engine/main.go
CGO_ENABLE=0 CGO_CFLAGS="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -ftrapv" go build -buildmode=pie -ldflags="-linkmode=external -extldflags '-Wl,-z,now' -X github.com/kappital/kappital/pkg/utils/version.gitVersion=0.1.0-dirty -X github.com/kappital/kappital/pkg/utils/version.gitCommit=9ee0c6a2d5489e183ca9978d8e08eabc68ecfad6 -X github.com/kappital/kappital/pkg/utils/version.gitTreeState="dirty" -X github.com/kappital/kappital/pkg/utils/version.buildDate=2022-10-13T07:03:36Z" -o bin/kappctl cmd/kappctl/main.go

```

After running the `make all` command, the project will have a `bin` directory which contains 3 binary files (`kappital-manager`, `kappital-engine` and `kappctl`). The `kappctl` is the binary file to connect to the `kappital-manager`. In addition, after deploy `kappital-manager`, please use the `kappctl config` command to set up the binary file.

## 2. Build All Images

Suggestion: **Deploy kappital project using containers.** 

```shell
[root@kappital kappital-fix]# make all-image
go fmt ./cmd/...
go fmt ./pkg/...
go vet ./cmd/...
go vet ./pkg/...
CGO_ENABLE=0 CGO_CFLAGS="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -ftrapv" go build -buildmode=pie -ldflags="-linkmode=external -extldflags '-Wl,-z,now' -X github.com/kappital/kappital/pkg/utils/version.gitVersion=0.1.0-dirty -X github.com/kappital/kappital/pkg/utils/version.gitCommit=9ee0c6a2d5489e183ca9978d8e08eabc68ecfad6 -X github.com/kappital/kappital/pkg/utils/version.gitTreeState="dirty" -X github.com/kappital/kappital/pkg/utils/version.buildDate=2022-10-13T07:09:52Z" -o bin/kappital-engine cmd/engine/main.go
docker build -t kappital/kappital-engine:latest -f build/kappital-engine/Dockerfile .
Sending build context to Docker daemon  174.8MB
Step 1/5 : FROM euleros:latest
latest: Pulling from library/euleros
8768d3961db9: Pull complete
Digest: sha256:49d8884e21486bb58c586976f6169e806dd087c0e60aee0a924ba73bb144565c
Status: Downloaded newer image for euleros:latest
 ---> edeb051fcc73
Step 2/5 : WORKDIR /opt/kappital/kappital-engine
 ---> Running in de66512bb5e7
Removing intermediate container de66512bb5e7
 ---> 9b1e95221c93
Step 3/5 : RUN groupadd -g 10000 kappital &&     useradd -u 10000 -g 10000 kappital &&     echo "Defaults targetpw" >> /etc/sudoers &&     mkdir -p /opt/kappital/kappital-engine &&     mkdir -p /opt/kappital/log &&     touch /opt/kappital/log/kappital-engine.log
 ---> Running in ee06c1839967
Removing intermediate container ee06c1839967
 ---> 7cd437d363d5
Step 4/5 : COPY bin/kappital-engine /opt/kappital/kappital-engine/
 ---> daa61cc33a94
Step 5/5 : RUN chown -R kappital:kappital /opt/kappital/kappital-engine &&     chown -R kappital:kappital /opt/kappital/log &&     chmod 750 /opt/kappital/kappital-engine &&     chmod 640 /opt/kappital/log/kappital-engine.log &&     chmod 550 /opt/kappital/kappital-engine/kappital-engine
 ---> Running in fce11c85fa4d
Removing intermediate container fce11c85fa4d
 ---> 72af693663da
Successfully built 72af693663da
Successfully tagged kappital/kappital-engine:latest
CGO_ENABLE=0 CGO_CFLAGS="-fstack-protector-all -D_FORTIFY_SOURCE=2 -O2 -ftrapv" go build -buildmode=pie -ldflags="-linkmode=external -extldflags '-Wl,-z,now' -X github.com/kappital/kappital/pkg/utils/version.gitVersion=0.1.0-dirty -X github.com/kappital/kappital/pkg/utils/version.gitCommit=9ee0c6a2d5489e183ca9978d8e08eabc68ecfad6 -X github.com/kappital/kappital/pkg/utils/version.gitTreeState="dirty" -X github.com/kappital/kappital/pkg/utils/version.buildDate=2022-10-13T07:09:52Z" -o bin/kappital-manager cmd/manager/main.go
docker build -t kappital/manager:latest -f build/kappital-manager/Dockerfile .
Sending build context to Docker daemon  174.8MB
Step 1/5 : FROM euleros:latest
 ---> edeb051fcc73
Step 2/5 : WORKDIR /opt/kappital/manager
 ---> Running in aea38ad84e21
Removing intermediate container aea38ad84e21
 ---> b0944ab4ab1d
Step 3/5 : RUN groupadd -g 10000 kappital &&     useradd -u 10000 -g 10000 kappital &&     echo "Defaults targetpw" >> /etc/sudoers &&     mkdir -p /opt/kappital/{manager,certs,log,audit,database} &&     mkdir -p /opt/kappital/certs/{conf}
 ---> Running in 214ee75db579
Removing intermediate container 214ee75db579
 ---> e168ff832e51
Step 4/5 : COPY bin/kappital-manager /opt/kappital/manager/
 ---> e8e43d0f2580
Step 5/5 : RUN chown -R kappital:kappital /opt/kappital/ &&     chmod 750 /opt/kappital/manager &&     chmod 550 /opt/kappital/manager/kappital-manager &&     chmod -R 700 /opt/kappital/certs
 ---> Running in fc7f6c658372
Removing intermediate container fc7f6c658372
 ---> 92bc3e850547
Successfully built 92bc3e850547
Successfully tagged kappital/manager:latest
```

Get the Kappital’s images will build the binary files first. If does not need the `kappctl` binary file, can just run the `make all-image` command to get all images.

## 3. Deploy Services

Kappital suggest deploying the services using the containers, and deploy at Kubernetes cluster. Here will introduce how to use Makefile and Helm Chart packages to deploy `kappital-manager`, and `kappital-engine` .

### 3.1 Deploy Kappital-Manager

```shell
[root@kappital kappital-fix]# make deploy-manager VERSION={image tag} MANAGER_SERVICE_IP={node port ip}
helm package charts/kappital-manager --version=0.0.1
Successfully packaged chart and saved it to: /root/kappital/kappital-manager-0.0.1.tgz
mkdir -p bin/package
mv kappital-manager-0.0.1.tgz bin/package/
# all kappital services will install in the namespace kappital-system
helm install kappital-manager \
        --create-namespace=false  \
        --namespace=kappital-system \
        --set manager.image.name=kappital/manager \
        --set manager.image.tag={image tag} \
        --set manager.image.registry="" \
        --set manager.image.pullPolicy=IfNotPresent \
        --set manager.logDir=/opt/kappital/manager/log \
        --set manager.hostNetwork=false \
        --set manager.service.clusterIP={node port ip} \
        --set manager.service.nodePort=30011  \
        bin/package/kappital-manager-0.0.1.tgz
NAME: kappital-manager
LAST DEPLOYED: Thu Oct 13 15:19:37 2022
NAMESPACE: kappital-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

Please refer to [installing by Helm](../../charts/kappital-manager/README.md).

### 3.3 Deploy Kappital-Engine

```shell
[root@kappital kappital-fix]# make deploy-kappital-engine VERSION={image tag}
helm package charts/kappital-engine --version=0.0.1
Successfully packaged chart and saved it to: /root/kappital/kappital-engine-0.0.1.tgz
mkdir -p bin/package
mv kappital-engine-0.0.1.tgz bin/package/
# all kappital services will install in the namespace kappital-system
helm install kappital-engine \
        --create-namespace=false  \
        --namespace=kappital-system \
        --set image.name=kappital/kappital-engine \
        --set image.tag={image tag} \
        --set image.registry="" \
        --set image.pullPolicy=IfNotPresent \
        --set livenessReadiness.port=8081 \
        --set livenessReadiness.scheme=HTTPS \
        bin/package/kappital-engine-0.0.1.tgz
NAME: kappital-engine
LAST DEPLOYED: Thu Oct 13 15:17:03 2022
NAMESPACE: kappital-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

Please refer to [installing by Helm](../../charts/kappital-engine/README.md)

### 3.4 Configure kappctl

```shell
[root@localhost kappital]# kappctl config \
  --manager-addr={Kappital-Manager IP address} \
  --manager-https-port={Kappital-Manager https service port if use http service} \
  --manager-client-cert={Kappital-Manager client certificat detail as base64} \
  --manager-client-key={Kappital-Manager client private key detail as base64} \
  --manager-ca={Kappital-Manager trust CA detail as base64} 
```

