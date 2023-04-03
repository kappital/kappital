# Kappital

## Kappital: Open, Cloud Native Service Full Lifecycle Management Platform

**Kappital** is an open source project that enables developer to manage cloud native application across multiple  clouds and edge, with no changes to developer's application. Kappital define [Cloud Native Service Package specification](docs/design/cloud_native_serivce_SPEC_DRAFT.md). By complying with spec, kappital enhance service capabilities,observability to application.



## Why Kappital



- **Unified Cloud Native Service Management**
  
  - Web console and command-line client for centrally manages multiple clusters
  - Full lifecycle management,such as install,upgrade,update,unInstall,state management
  
- **Declarative Observability With Non-Intrusive**
  
  - Zero-Code for logging,monitoring and alarm
  - Support CNCF Observability project,such as prometheus, OpenTelemetry,Thanos
  
- **Multi-Deployment Engine**

  - Support  Kubernetes Native deployment 

  - Support  Helm,Operator-Framework, and other deploy framework

    

## Architecture



![Architecture](docs/images/architecture.png)

The Kappital consists of the following components:

- __Kappital-Manager__
  - Support full lifecycle of Cloud-Native Service，such install, upgrade,update,uninstall .
  - Deploy Multi-cloud/multi-cluster/cloud-edge without any 
  - Support Day2 Operation
- __Kappital-Service-Engine__
  - Support popular framework such as helm,operator-framework
  - Provide observability plugin,enable monitor,logging,alarm,etc



## Specification

Kappital define CloudNativeService Specification to 



More information, please refer to [Cloud Native Service SPEC](docs/design/cloud_native_serivce_SPEC_DRAFT.md).

## Core Concepts



- __CloudNativeService__

  Kappital define **CloudNativeService** to describe application that supports multi-cloud and cloud-edge environments, use cloud-native technology stack. Provides full lifecycle management capabilities, including sales, installation, management, and maintenance.
  
  
  
  More information, please refer to [CloudNativeService](docs/design/cloud_native_service.md).



- __CloudNativeServiceInstance__

  Kappital define **CloudNativeServiceInstance** to describe when deploy a __CloudNativeService__ in a cluster environment will return a __ServiceInstance__. It include application instance, service capabilities and observability capabilities added by platform.

  

   More information, please refer to [CloudNativeServiceInstance](docs/design/2 - cloud_native_service_instance.md).



## Getting Started



### Installation

Kappital v0.1.0(alpha) requires Kubernetes version >= 1.17, [installation reference](./docs/installation/installation.md).


## License

Kappital is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.