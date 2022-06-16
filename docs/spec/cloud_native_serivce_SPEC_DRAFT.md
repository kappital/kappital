# Cloud Native Service Package Specification

Cloud Native Service Specification(spec) provide a standard definition of cloud-native services irrelevant to the cloud platform, describe the deployment and governance capabilities of cloud-native services.

## Specification 

```shell
kappital-package/ 
├─ metadata.yaml                 # declare service meta message
├─ manifests/                    # custom resource and extension resource
│  ├─ crd.yaml                   # kubernetes CustomerResourceDefinition file
│  └─ csd.yaml                   # CustomerServiceDefinition file，extension of CustomerResourceDefinition
├─ capability/                   # capability config,inclue bbservability configuration, servicization configuration
│  ├─ monitor.yaml               # monitor demo 
│  ├─ logs.yaml                  # log demo
│  └─ alarm.yaml                 # alarm demon
├─ operator/                     # operator deploy config
│  ├─ operator-deployment.yaml   # operator deployment config
│  ├─ operator-clusterrole.yaml  # operator rabc config
│  └─ operator-role.yaml         # operator rabc config
└─ raw/                          # third-party package eg: helmchart,operator-framework 
    └─ third-party-v1.0.tgz   
```



## Metadata File

 Metadata File describe cloud-native service's basic attributes, such as name, version, type,description. More information see table below:

 

| Field Name        | Description                                                                                                  | Value Example                               | Required |
|-------------------|--------------------------------------------------------------------------------------------------------------|---------------------------------------------|----------|
| name              | service name, format "^[a-z][a-z0-9\_\-]*[a-z0-9]$"                                                          | kappital-service                            | YES      |
| version           | service version                                                                                              | 0.1.0                                       | YES      |
| type              | service' type                                                                                                | [helm, operator]                            | YES      |
| alias             | service alias name                                                                                           | kappital-service-0.1.0                      | NO       |
| briefDescription  | short description for service                                                                                | NA                                          | NO       |
| detailDescription | full description for service                                                                                 | NA                                          | NO       |
| source            | service's source                                                                                             | bitnami                                     | NO       |
| repository        | store repository                                                                                             | default                                     | NO       |
| minKubeVersion    | support kubernetes minimun version                                                                           | 1.15                                        | NO       |
| architecture      | support architecture                                                                                         | [x86, Arm64]                                | NO       |
| capabilities      | define service capability. Scope: Basic Install,Seamless Upgrades, Full Lifecycle, Deep Insights, Auto Pilot | [Basic Install]                             | NO       |
| category          | define service category, such as Database, Monitoring, Networking, Security, Storage                         | [Database, Storage]                         | NO       |
| devices           | define service use devices,such as CPU, GPU, MPU                                                             | [GPU]                                       | NO       |
| industries        | define service belong to industry                                                                            | [education]                                 | No       |
| logo              | service logo, include attributes: base64data, mediaType                                                      | base64data: xxxx <br/> mediaType: image/png | NO       |
| maintainers       | service maintainer, include attributes: name, email                                                          | NA                                          | NO       |
| provider          | service provider, include attributes: name, url                                                              | NA                                          | NO       |
| scenes            | support deploy scenes,                                                                                       | [Cloud, Edge]                               | NO       |
| links             | service link message, include attributes: name, url                                                          | NA                                          | NO       |



**medata yaml example:**

``` yaml
name: kappital-service
version: 0.0.3 
alias: kappital-service-v0.0.3 
briefDescription: kappital operator with an example instance and action 
detailDescription: | 
  kapppital service detail description. 
source: bitnami 
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
keywords: 
  - Streaming & Messaging 
experience: 
  name: Example Demo 
  url: https://example-demo.com/
```

## Manifest  Directory

manifest directory container two kind file: **CustomResourceDefinition(CRD)** and **CustomServiceDefinition(CSD)**. 

- CRD file define by kubernetes, more information refer to [CustomResourceDefinition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#customresourcedefinition-v1-apiextensions-k8s-io).

- CSD file define by Kappital, which to enhance CRD, support more capability such as service dependency, service access, service deployment.



**1. CustomServiceDefinition(CSD)**:  CSD use GVK to define what in CSD file, A CSD file w

```
apiVersion: osc.io/v1beta1 
kind: CustomServiceDefinition 
metadata: 
  name: csd-name
spec:
  ... ## CSD parameter config
```

The following table describes CSD' parameter:

| parameter    | descripition                                                 | example                    | required |
| ------------ | ------------------------------------------------------------ | -------------------------- | -------- |
| crd          | define to  associate with CRD,                               | NA                         | YES      |
| crd.name     | specify CRD name                                             | wildflyservers.wildfly.org | YES      |
| crd.version  | specify CRD version                                          | v1alpha1                   | YES      |
| crd.versions | specify CRD all version and default values.                  | NA                         | NO       |
| description  | extend CRD' more description                                 |                            | NO       |
| role         | define role to extend CRD,  include: ServiceEntity, Operation. Default ServiceEntity | ServiceEntity              | NO       |
|              |                                                              |                            |          |



**csd yaml example:**

```yaml
apiVersion: core.kappital.io/v1alpha1
kind: CustomServiceDefinition
metadata:
  name: wildfly-csd
spec:
  crd:
    name: wildflyservers.wildfly.org
    version: v1alpha1
    versions:
    - name: v1alpha1
      defaultValues: |-
        {
          "applicationImage": "quay.io/wildfly-quickstarts/wildfly-operator-quickstart:18.0",
          "replicas": 2
        }
    - name: v1
      defaultValues: |-
        {
          "applicationImage": "quay.io/wildfly-quickstarts/wildfly-operator-quickstart:18.0",
          "replicas": 2
        }
  role: ServiceEntity
  description: wildfly crd description
```



## Capability  Directory

**Capability directory** store all capability configuration files , include operation file link logging, monitoring, alarm, and more observability files.



## Operator Directory



**Operator directory** store all operator deploy files when this service's type is **Operator**. When service's type is **Helm**, this directory is empty.

Generally, operator diectory containers operator workload file, role and rolebind file , clusterrole and clusterrolebind file. 



## Rw Directory

**raw directory** store third-party service packages, support helmchart, operator-framework bundle.







