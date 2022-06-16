# Custom Service Definition

## Introduction
**CustomServiceDefinition**, shorted by ***CSD* and expanding the original abilities of [CustomResourceDefinition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#customresourcedefinition-v1-apiextensions-k8s-io),
is a unique and core concept of Kappital. It provides more abilities which are required and common used in industries such
as service dependency, service access, service deployment, monitoring and logging.

## Key Fields

### Top-Level Attributes

The top-level attributes of a cloud-native service define its apiVersion, kind, metadata and spec.

| Attribute    | Type                    | Required | Default Value               | Description                                                                                                                                                  |
|--------------|-------------------------|----------|-----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `apiVersion` | `string`                | Y        | `core.kappital.io/v1alpha1` | A string that identifies the version of the schema the object should have. The core types uses `core.kappital.io/v1alpha1` in this version of documentation. |
| `kind`       | `string`                | Y        | `CustonServiceDefinition`   | Must be `CustonServiceDefinition`                                                                                                                            |
| `metadata`   | [`Metadata`](#Metadata) | Y        |                             | Information about the CSD resource.                                                                                                                          |
| `spec`       | [`Spec`](#spec)         | Y        |                             | A specification for the CSD resource attributes.                                                                                                             |

### Metadata

Metadata provides basic information about the CSD.

| Attribute     | Type                | Required | Default Value | Description                                                                                                                                                                                                                                                   |
|---------------|---------------------|----------|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`        | `string`            | Y        |               | A name for the schematic. `name` is subject to the restrictions listed beneath this table.                                                                                                                                                                    |
| `labels`      | `map[string]string` | N        |               | A set of string key/value pairs used as arbitrary labels on this component. Labels follow the [Kubernetes specification](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).                                                          |
| `annotations` | `map[string]string` | N        |               | A set of string key/value pairs used as arbitrary descriptive text associated with this object.  Annotations follows the [Kubernetes specification](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#syntax-and-character-set). |

### Spec

The specification of cloud-native services defines service metadata, version list, service capabilities and plugins.

| Attribute      | Type                                 | Required | Default Value | Description                                                                             |
|----------------|--------------------------------------|----------|---------------|-----------------------------------------------------------------------------------------|
| `crd`          | [`CRDDescriptor`](####CRDDescriptor) | Y        |               |                                                                                         |
| `description`  | `string`                             | N        |               | Extend the CRD's description                                                            |
| `role`         | `string`                             | N        |               | Define a role to extended CRD including ServiceEntity, Operation. Default ServiceEntity |

#### CRDDescriptor
| Attribute  | Type           | Required | Default Value | Description                                          |
|------------|----------------|----------|---------------|------------------------------------------------------|
| `name`     | `string`       | Y        |               | Specified CRD name                                   |
| `version`  | `string`       | N        |               | Specified CRD version                                |
| `versions` | `string array` | N        |               | All versions and default values of the specified CRD |

## Example
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