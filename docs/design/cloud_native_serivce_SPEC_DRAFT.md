# Cloud Native Service Package Specification

Cloud Native Service Specification(spec) provide a standard definition of cloud-native services irrelevant to the cloud platform, describe the deployment and governance capabilities of cloud-native services.

## Specification 

```shell
kappital-package/ 
├─ metadata.yaml    	# declare service meta message
├─ manifests/           # custom resource and extension resource
│  ├─ crd.yaml          # kubernetes CustomerResourceDefinition file
│  └─ csd.yaml          # CustomerServiceDefinition file，extension of CustomerResourceDefinition
├─ capability/          # capability config, include serviceability configuration, serialization configuration (Planning)
│  ├─ monitor.yaml      # monitor demo 
│  ├─ logs.yaml 	    # log demo
│  └─ alarm.yaml        # alarm demon
├─ operator/            # operator deploy config
│  ├─ operator-deployment.yaml 
│  ├─ operator-clusterrole.yaml
│  └─ operator-role.yaml
└─ raw/ 	            # third-party package (Planning)   
    └─ third-party-v1.0.tgz   
```



## Package Example

 [Base Example](../../examples/quickstart-package)



## File Example

### 1. metadata.yaml

``` yaml
name: example-operator 
version: 0.0.3 
alias: example-operator.v0.0.3 
displayName: Example operator 
briefDescription: example operator with an example instance and action 
detail: | 
  example operator detail description. 
source: bitnami 
type: operator|helm 
containerImage: xxxx:xxx 
repository: xxx 
minKubeVersion: 1.15.0 
architecture: 
  - x86 
  - arm64 
capabilities: 
  - Basic Install 
categories: 
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
  url: "the url of website, format as http or https" 
scenes: 
  - edge 
links: 
  - name: example link 
    url: "the url of website, format as http or https" 
keywords: 
  - Apache Kafka 
  - Apache ZooKeeper 
  - Streaming & Messaging 
experience: 
  name: Example Demo 
  url: "the url of website, format as http or https" 
```



### 2. csd.yaml

```yaml
# nonk8s
apiVersion: core.kappital.io/v1alpha1
kind: CustomServiceDefinition
metadata:
  name: kafka-csd
spec:
  # refer Custome Resource Definition
  definitionRef:
    apiVersion: apiextensions.k8s.io/v1beta1
    kind: CustomResourceDefinition
    name: kafka.example.com

  # define role to extension crd, enum:{ServiceEntity,Operation}
  role: ServiceEntity

  # public capability, not support in v1alpha1
  capabilityRequirements:
    - apiVersion: core.kappital.io/v1
      kind: MonitorConfig
      defaultPath: capability/monitor_config.yaml
    - apiVersion: core.kappital.io/v1
      kind: LogConfig
      defaultPath: capability/log_config.yaml
    - apiVersion: core.kappital.io/v1
      kind: AlarmConfig
      defaultPath: capability/alarm_config.yaml
```









