# 1. Cloud-Native Service Object

This section define cloud-native service object, and how __CloudNativeService__ are designed and published.



## Definition



// TODO



### Top-Level Attributes

The top-level attributes of a cloud-native service define its apiVersion, kind, metadata and spec.

| Attribute    | Type                    | Required | Default Value                                       | Description                                                                                                                                                  |
|--------------|-------------------------|----------|-----------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `apiVersion` | `string`                | Y        | `core.kappital.io/v1alpha1`                         | A string that identifies the version of the schema the object should have. The core types uses `core.kappital.io/v1alpha1` in this version of documentation. |
| `kind`       | `string`                | Y        | `CloudNativeService`                                | Must be `CloudNativeService`                                                                                                                                 |
| `metadata`   | [`Metadata`](#Metadata) | Y        |                                                     | Information about the cloud-native service.                                                                                                                  |
| `spec`       | [`Spec`](#spec)         | Y        || A specification for cloud-native serviceattributes. |



### Metadata

Metadata provides information about the contents of the object.

| Attribute     | Type                | Required | Default Value | Description                                                                                                                                                                                                                                                                                                              |
|---------------|---------------------|----------|---------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`        | `string`            | Y        |               | A name for the schematic. `name` is subject to the restrictions listed beneath this table.                                                                                                                                                                                                                               |
| `labels`      | `map[string]string` | N        |               | A set of string key/value pairs used as arbitrary labels on this component. Labels follow the [Kubernetes specification](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).<br/>cloud-native service has unique labels:<br/>- xxx                                                               |
| `annotations` | `map[string]string` | N        |               | A set of string key/value pairs used as arbitrary descriptive text associated with this object.  Annotations follows the [Kubernetes specification](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#syntax-and-character-set).<br/>cloud-native service has unique annotations:<br/>- xxx |



### Spec

The specification of cloud-native services defines service metadata, version list, service capabilities and plugins.

| Attribute      | Type         | Required | Default Value | Description                                                        |
|----------------|--------------|----------|---------------|--------------------------------------------------------------------|
| `descriptor`   | `Descriptor` | Y        |               | A string that identifies cloud-native service detail informations. |
| `versionKinds` | `[]Version`  | Y        |               | A set of string list describe as cloud-native service's  version.  |



### Descriptor

The descriptor of cloud-native service describe service detail information.

| Attribute          | Type     | Required | Default Value | Description                                           |
|--------------------|----------|----------|---------------|-------------------------------------------------------|
| `Name`             | `string` | N        |               |                                                       |
| `DisplayName`      | `string` | N        |               | Display service name to end users                     |
| `Type`             | `string` | N        |               |                                                       |
| `version`          | `string` | Y        |               | A string that identifies cloud-native service version |
| `briefDescription` | `string` | N        |               | Brief description of cloud-native service             |





### Version

The version collection of all cloud-native service versions.

| Attribute  | Type       | Required | Default Value | Description                                                        |
|------------|------------|----------|---------------|--------------------------------------------------------------------|
| `version`  | `string`   | Y        |               | cloud-native service version                                       |
| `replaces` | `[]string` | N        |               | A set of string list, that can be replaced by the current version. |

## Example

```yaml
# nonk8s
apiVersion: core.kappital.io/v1alpha1
kind: CloudNativeService
metadata:
  name: cloud-native-service-456465
  uid: c15ec843-96be-4a5f-8ab5-da14aed7ae93
spec:
  descriptor:
    version: 1.0
    type: operator|helm
    briefDescription: example operator with an example instance and action
    detail: |
      Kafka allows you to quickly deploy [Apache Kafka®](https://kafka.apache.org) in a [Kubernetes](https://kubernetes.io/) cluster.
      Features Supported by 21.9. 18
      - InstallationDeployment: Deploy and manage Apache Kafka clusters and components such as Apache ZooKeeper® on which they depend.
      - Elastic capacity expansion: Quickly scale out ZooKeeper nodes or Kafka Broker nodes.
      - Topic Management: Use Kafka Manager to manage topics in the Kafka cluster.
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
    scenes:
      - edge
    logo:
      base64data: xxxx
      mediaType: image/png
    maintainers:
      - email: test@test.com
        name: test
    provider:
      name: Example provider
      url: https://example.com/
    links:
      - name: example link
        url: https://example.com/
    keywords:
      - Apache Kafka
      - Apache ZooKeeper
      - Streaming & Messaging
    experience:
      name: Example Demo
      url: https://example-deom.com/
    industries:
      - education
      - media

  versionKinds:
    - version: 1.0
      replaces: [ ]
    - version: 2.0
      replaces: [ 1.0,1.2,1.3 ]
    - version: 3.0
      replaces: [ 2.0,2.5 ]

  manifests:
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        annotations:
          controller-gen.kubebuilder.io/version: v0.3.0
          kappital.io/resource-role: ServiceEntity
        name: kafkas.example.com
      spec:
        group: example.com
        names:
          kind: Kafka
          listKind: KafkaList
          plural: kafkas
          shortNames:
            - kfk
          singular: kafka
        scope: Namespaced
        subresources:
          status: { }
  capabilities:
    - apiVersion: core.kappital.io/v1alpha1
      kind: AlarmConfig
      metadata:
        name: kafka-alarm
      spec:

    - apiVersion: core.kappital.io/v1alpha1
      kind: LogConfig
      metadata:
        name: kafka-logconfig
      spec:


  plugins:
    - apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: redis-operator-deployment
        labels:
          app: nginx
      spec:
        replicas: 1
        selector:
          matchLabels:
            name: redis-operator-deployment
        strategy: { }
        template:
          metadata:
            labels:
              name: redis-operator-deployment
          spec:
            containers:
              - args:
                  - --enable-leader-election
                command:
                  - /bin/bash
                  - -ec
                  - test
                image: xxxx
                imagePullPolicy: IfNotPresent
                name: redis-operator
                volume:
                  - name: test
                    path: /var/paas/sys/log/operator/xxx-operator.log
            volumeMount:
              - name: test
                path: /var/paas/sys/log/operator/xxx-operator.log
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: RoleBinding
      metadata:
        name: read-pods
        namespace: default
      subjects:
        - kind: User
          name: dave
          apiGroup: rbac.authorization.k8s.io
      roleRef:
        kind: ClusterRole
        name: redis-operator-clusterrole
        apiGroup: rbac.authorization.k8s.io
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRole
      metadata:
        name: redis-operator-clusterrole
      rules:
        - apiGroups:
            - ""
          resources:
            - ""
          verbs:
            - ""
    - apiVersion: rbac.authorization.k8s.io/v1
      kind: RoleBinding
      metadata:
        name: read-pods
        namespace: default
      subjects:
        - kind: User
          name: dave
          apiGroup: rbac.authorization.k8s.io
      roleRef:
        kind: ClusterRole
        name: redis-operator-clusterrole
        apiGroup: rbac.authorization.k8s.io

  raw:
    type: helm
    spec:
      apiVersion: helm.toolkit.kappital.io/v1alpha1
      kind: HelmRelease
      spec:
        release: HelmReleaseSpec
        repository: HelmRepositorySpec
```







