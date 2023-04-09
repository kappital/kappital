# Cloud Native Service Package Specification

Cloud Native Service Package Specification provides a standard definition of cloud-native service irrelevant to the cloud
platform and describes the deployment and governance of cloud-native service.

## Specification 

```shell
service-package/ 
├─ metadata.yaml                 # declare service metadata
├─ manifest/                     # custom resource and extension resource
│  ├─ crd.yaml                   # Kubernetes CustomerResourceDefinition file
│  └─ csd.yaml                   # CustomerServiceDefinition file, extension of CustomerResourceDefinition
├─ capability/                   # capability directory including observability and servitization configuration
│  ├─ monitor.yaml               # monitoring configuration 
│  ├─ logs.yaml                  # logging configuration
│  └─ alarm.yaml                 # alarm configuration
├─ operator/                     # operator deployment directory
│  ├─ operator-deployment.yaml   # operator deployment configuration
│  ├─ operator-clusterrole.yaml  # operator RBAC configuration, eg: ClusterRole.
│  └─ operator-role.yaml         # operator RBAC configuration, eg: Role.
└─ raw/                          # third-party package directory, eg: Helm chart, operator-framework. 
    └─ third-party-v1.0.tgz   
```

## Metadata File

 Metadata File describes cloud-native service's basic attributes such as name, version, type, description. More information see the table below:

| Field Name          | Description                                                                                                    | Required | Value Example                               |
|---------------------|----------------------------------------------------------------------------------------------------------------|----------|---------------------------------------------|
| `name`              | service name, format "^[a-z][a-z0-9\_\-]*[a-z0-9]$"                                                            | Y        | kappital-service                            |
| `version`           | service version                                                                                                | Y        | 0.1.0                                       |
| `type`              | service type                                                                                                   | Y        | [helm, operator]                            |
| `alias`             | service alias name                                                                                             | N        | kappital-service-0.1.0                      |
| `briefDescription`  | short description for service                                                                                  | N        | NA                                          |
| `detailDescription` | full description for service                                                                                   | N        | NA                                          |
| `source`            | service source such as artifacthub                                                                             | N        | bitnami                                     |
| `repository`        | store repository                                                                                               | N        | default                                     |
| `minKubeVersion`    | support Kubernetes minimum version                                                                             | N        | 1.15                                        |
| `architecture`      | support architecture                                                                                           | N        | [x86, Arm64]                                |
| `capabilities`      | define service capability. Scope: Basic Install, Seamless Upgrades, Full Lifecycle, Deep Insights, Auto Pilot. | N        | [Basic Install]                             |
| `category`          | define service category, such as Database, Monitoring, Networking, Security, Storage.                          | N        | [Database, Storage]                         |
| `devices`           | devices the service making use of such as CPU, GPU, MPU.                                                       | N        | [GPU]                                       |
| `industries`        | the industry the service applied to, for example education application.                                        | N        | [education]                                 |
| `logo`              | service logo including base64data, mediaType.                                                                  | N        | base64data: xxxx <br/> mediaType: image/png |
| `maintainers`       | service maintainer including attributes: name, email.                                                          | N        | NA                                          |
| `provider`          | service provider including attributes: name, url.                                                              | N        | NA                                          |
| `scenes`            | scenes the service applied to such as cloud, edge.                                                             | N        | [Cloud, Edge]                               |
| `links`             | link providing more information for the service including attributes: name, url.                               | N        | NA                                          |


**medata yaml example:**

``` yaml
name: kappital-service
version: 0.0.1 
alias: kappital-service-v0.0.1 
briefDescription: kappital operator with an example and action
detailDescription: | 
  kapppital service detail description. 
source: artifacthub 
type: operator|helm 
repository: defualt
minKubeVersion: 1.15.0 
architecture: 
  - x86 
  - arm64 
capabilities: 
  - Basic Install 
category: 
  - Database 
devices: 
  - GPU 
  - MPU 
industries: 
  - education 
  - media 
logo: 
  base64data: xxxx 
  mediaType: image/png 
maintainers: 
  - email: test@test.com 
    name: test 
provider: 
  name: Example provider 
  url: https://example.com/ 
scenes: 
  - edge 
links: 
  - name: example link 
    url: http://github.com/ 
```

## Manifest Directory

Manifest directory contains two kinds of files: **CustomResourceDefinition(CRD)** and **CustomServiceDefinition(CSD)**. 

- CRD file is defined following Kubernetes. More information refers to [CustomResourceDefinition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#customresourcedefinition-v1-apiextensions-k8s-io).
- CSD file is defined following Kappital, which enhances CRD, supporting more abilities such as service dependency, service access, service deployment. More information refers to [CustomServiceDefinition](./custom_service_definition.md).


## Capability  Directory

**Capability Directory** stores all capability configuration files including operation files such as logging, monitoring, alarm and observability.

## Operator Directory

**Operator Directory** stores all operator deployment files when the service's type is **Operator**, while the directory is empty when the service's type is **Helm**.

Generally, operator directory contains operator workload files, role and rolebinding files, clusterrole and clusterrolebinding files. 



## Raw Directory

**Raw Directory** stores third-party service packages which supports Helm chart, operator-framework bundle.