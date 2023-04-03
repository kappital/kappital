# Kappital-ServiceEngine

## Prerequisites
- Kubernetes 1.15+
- Helm v3+

## Installing the Service Engine

Switch to the upper-level directory (`kappital/charts`).
```shell
[root@localhost charts]$ helm install kappital-engine -n kappital-system ./kappital-engine
```

> Tip: list all release using `helm list -n kappital-system` or `helm list -A`

## Uninstall the Chart

```shell
[root@localhost kappital]$ helm uninstall kappital-engine -n kappital-system
```

## Configuration

| Name                       | Description                                                  | Default Value            |
| -------------------------- | ------------------------------------------------------------ |--------------------------|
| `namespace`                | Service Engine running namespace. For Kappital project, we suggest the Service Engine is running in the default namespace. | `kappital-system`        |
| `image.name`               | The image name of Service Engine. If you edit the image name with others, please change this value. | `kappital/kappital-engine` |
| `image.tag`                | The version of this image.                                   | `latest`                 |
| `image.registry`           | The registry store the image.                                | `""`                       |
| `image.pullPolicy`         | The pull policy when deploy Service Engine.                  | `IfNotPresent`           |
| `log.hostPath`             | This configuration is using for volume the log file to the host machine. You can change the log path as you prefer. In addition, if you change this value to empty, the Service Engine will use empty directory method to volume logs. | `/opt/kappital/kappital-engine/log/` |
| `livenessReadiness.port`   | The liveness and rediness port for Kubernetes. Because the Service Engine is an Operator, here will use the default port number to do the health check. If you intrusively modify the code, please change this port as your code. | `8081`                   |
| `livenessReadiness.scheme` | The liveness and rediness port for Kubernetes for health checking. The default operator health check is using the HTTP method. If you think HTTP is insecure, please change it to the HTTPS and intrusively modify the code. | `HTTPS`                  |

## Cluster Permission Clarification

Kappital-Engine converts and manages cloud-native service concepts and Kubernetes built-in resources. During actual running, Kappital-Engine listens to the service package resource objects defined by Kappital. Convert to Kubernetes built-in resource objects (such as Deployments, CustomResourceDefinition, etc.) and create them in Kubernetes. Kappital-Engine needs to manage the entire lifecycle of services and service instances and maintain Kubernetes resources. Therefore, Kappital-Engine must have the permission to manage all Kubernetes resources.
