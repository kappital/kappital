
```
apiVersion: core.kappital.io/v1alpha1 
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