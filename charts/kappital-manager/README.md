# Kappital-Manager

## Prerequisites
- Kubernetes 1.15+
- Helm v3+

## Installing the Manager

Switch to the upper-level directory (`kappital/charts`).
```shell
[root@localhost charts]$ helm install kappital-manager -n kappital-system ./kappital-manager
```

> Tip: list all release using `helm list -n kappital-system` or `helm list -A`

## Uninstall the Chart

```shell
[root@localhost kappital]$ helm uninstall kappital-manager -n kappital-system
```

## Configuration

| Name                        | Description                                                                                                                                                                                                                     | Default Value               |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------|
| `namespace`                 | Manager running namespace. For Kappital project, we suggest the Manager is running in the default namespace.                                                                                                                    | `kappital-system`           |
| `manager.image.name`        | The image name of Manager. If you edit the image name with others, please change this value.                                                                                                                                    | `kappital/manager`          |
| `manager.image.tag`         | The version of this image.                                                                                                                                                                                                      | `latest`                    |
| `manager.image.registry`    | The registry store the image.                                                                                                                                                                                                   | `""`                          |
| `manager.image.pullPolicy`  | The pull policy when deploy Manager.                                                                                                                                                                                            | `IfNotPresent`              |
| `manager.logDir`            | This configuration is using for volume the log file to the host machine. You can change the log path as you prefer. In addition, if you change this value to empty, the Manager will use empty directory method to volume logs. | `"/opt/kappital/manager/log/"` |
| `manager.hostNetwork`       | Manager does not use the host machine port as default.                                                                                                                                                                          | `false`                     |
| `manager.checkIdentity`      | Does the Manager whill check the client's CN attribute. Manager will check this attribute as default.                                                                                                                           | `true`                      | 
| `manager.enableMutualHttps` | Does the Manager will check the client's certificate. Manager will check the client certificate as default.                                                                                                                     | `true`                      |
| `manager.tlsConfig` | Manager will use which method to check the client's certificate.                                                                                                                                                                | `REQUIRE_AND_VERIFY_CLIENT_CERT` |
| `manager.service.clusterIP` | If Manager does not use host network, it will provided service through the Service with ClusterIP method.                                                                                                                       | `x.x.x.x`                   |
| `manager.service.httpsPort` | Which port does the Manager will use in the container                                                                                                                                                                           | `30330`                     |
| `manager.service.httpsCertFile` | The server certificate path for Manager.                                                                                                                                                                                        | `/opt/kappital/certs/conf/server.crt` |
| `manager.service.httpsKeyFile` | The server private key path for Manager.                                                                                                                                                                                        | `/opt/kappital/certs/conf/server.key` |
| `manager.service.trustCaFile` | The client root ca file path for Manager. (This will use for check the client certificate if given)                                                                                                                             | `/opt/kappital/certs/conf/ca.crt` |

## Service Exposure Mode

Kappital-Manager position is deploying in the Kubernetes cluster, and maintain the CloudNativeService and CloudNativeInstance resources. Sometime, it will maintain some custom resources (CR) with the CloudNativeInstance. Thus, the Kappital-Manager is a local tool for the Kubernetes Cluster. It will use the `ClusterIP` as the default service exposure mode.

If `ClusterIP` cannot be used due to the using case or situation, Please change the Service (Kubernetes resource) to NodePort or Other method. The Service resource is in the [manager.yaml](./templates/manager.yaml), Please change it by yourself.

**ATTN**: If using the `NodePort` or other similar method to exposure service, Please reinforce the firewall to prevent security problems.

## Cluster Permission Clarification

Kappital-Manager will get, create, delete, list and update different Custom Resources (from Kubernetes `CustomResourceDefinition` resource). The different `CustomResourceDefinition` will have different `resources` and `apiGroups` during the `ClusterRole`'s "rules" attribute.

> For example:
> 1. Kappital-Manager deploy a CloudNativeInstance for a Package A. Package A needs Kappital-Manager to deploy a CR with apiGroups named `api.groups.a` and resource named `resource-a`, and this CR need to deploy into `default` namespace. Here, Kappital-Manager needs the permission for apiGroups `api.groups.a` and resource `resource-a`, also, this CR is deployed into `default` namespace which is different with the Kappital-Manager running namespace `kappital-system` (does not change the Kappital-Manager running namespace). Thus, the Kappital-Manager at least needs the ClusterRole Permission.
> 2. Now, Kappital-Manager deploy a CloudNativeInstance for a Package B. Package B needs Kappital-Manager to deploy a CR with apiGroups named `api.groups.b` and resource named `resource-b`, and this CR need to deploy into `ns-b` namespace. Here, Kappital-Manager needs the permission for apiGroups `api.groups.b` and resource `resource-b`, also, this CR is deployed into `ns-b` namespace which is different with the Kappital-Manager running namespace `kappital-system` (does not change the Kappital-Manager running namespace). Thus, with the Package A situation, Kappital-Manager at lease need the Permission for the Package A and B.
> 3. Therefore, the Packages' CRD kinds are unpredictable. Kappital-Manager will have the all apiGroups and resource permission to solve this situation.

However, if the service packages and CRDs to be deployed are predictable. You can change the [configs.yaml](./templates/configs.yaml) file. Changing the `system:controller:manager` ClusterRole to these packages.

**ATTN**: If changing the ClusterRole permission, Kappital-Manager is not able to deploy other package without these CRDs.

