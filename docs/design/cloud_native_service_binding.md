# 2. Cloud-Native Service Instance Object

This section define cloud-native service instance object, and how __CloudNativeServiceInstance__  are designed and managerment.

## Definition

// TODO

### Top-Level Attributes

The top-level attributes of a cloud-native service instance define its apiVersion, kind, metadata and spec.

| Attribute    | Type                    | Required | Default Value                | Description                                                                                                                                                  |
|--------------|-------------------------|----------|------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `apiVersion` | `string`                | Y        | `core.kappital.io/v1alpha1`  | A string that identifies the version of the schema the object should have. The core types uses `core.kappital.io/v1alpha1` in this version of documentation. |
| `kind`       | `string`                | Y        | `CloudNativeServiceInstance` | Must be `CloudNativeServiceInstance`                                                                                                                         |
| `metadata`   | [`Metadata`](#Metadata) | Y        |                              | Information about the cloud-native service instance.                                                                                                         |
| `spec`       | `Spec`                  | Y        |                              | A specification for cloud-native service instance attributes.                                                                                                |





### Metadata

Metadata provides information about the contents of the object.

| Attribute     | Type                | Required | Default Value | Description                                                                                                                                                                                                                                                                                                                       |
|---------------|---------------------|----------|---------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`        | `string`            | Y        |               | A name for the schematic. `name` is subject to the restrictions listed beneath this table.                                                                                                                                                                                                                                        |
| `labels`      | `map[string]string` | N        |               | A set of string key/value pairs used as arbitrary labels on this component. Labels follow the [Kubernetes specification](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).<br/>cloud-native service instance has unique labels:<br/>- xxx                                                               |
| `annotations` | `map[string]string` | N        |               | A set of string key/value pairs used as arbitrary descriptive text associated with this object.  Annotations follows the [Kubernetes specification](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#syntax-and-character-set).<br/>cloud-native service instance has unique annotations:<br/>- xxx |





## Example

```yaml
# nonk8s
apiVersion: core.kappital.io/v1beta1
kind: CloudNativeServiceBinding
metadata:
  name: cloud-native-service-binding-c15ec843
  namespace: default
spec:
  clusterId: c15dsafsd-96be-4a5f-8ab5-asdfads
  name: "xxx"
  version: v1.1.0
  serviceName: "xxx"
  serviceId: "xxx"
  serviceType: "operator"
  customResources:
    - apiVersion: redis.io/v1beta1
      kind: Redis
      name: redis-0aa4t1
    - apiVersion: redis.io/v1beta1
      kind: RedisClusterRestore
      name: redis-restore-0aa4t2
    - apiVersion: redis.io/v1beta1
      kind: RedisClusterBackup
      name: redis-backup-0aa4t3
  capabilitiesResources:
    - apiVersion: core.kappital.io/v1alpha1
      kind: AlarmConfig
      name: alarmconfig-xxx
status:
  installStatus:
    status: ""
    phase: ""
    message: ""
    servicePackName: ""
    servicePackVersion: ""
  dependencyStatus:
    - status: ""
      phase: ""
      message: ""
      servicePackName: ""
      servicePackVersion: ""
```

