# Command-Line Tool : kappctl

This section define kappital CLI , introduce core commands to manage cloud-native services' full lifecycle.

| Concept                    | Short Name |
|----------------------------|------------|
| CloudNativeService         | service    |
| CloudNativeServiceInstance | instance   |

## Definition

```shell
kappctl manages the full lifecycle of Kappital Resources

Usage:
  kappctl [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  config      Config kappctl with Kappital-Manager
  create      Create a Kappital resource
  delete      Delete a Kappital resource
  get         Display one or many Kappital resources
  help        Help about any command
  init        Create a Kappital package scaffold from scratch.

Flags:
  -h, --help   help for kappctl
```

## Core Command

### 1. Init the Service Package

```shell
Create a Kappital package scaffold from scratch.

Usage:
  kappctl init [flags]

Flags:
  -h, --help             help for init
      --name string      the kappital package name (default "kappital-demo")
  -v, --version string   the kappital package version (default "0.1.0")
```

### 2. Config the `kappctl`

```shell
Config kappctl with Kappital-Manager

Usage:
  kappctl config [flags]

Flags:
  -h, --help                         help for config
      --manager-addr string          the ip Addr of kappital-manager
      --manager-ca string            the HTTPS ca file for kappital-manager
      --manager-client-cert string   the HTTPS client certificate file of kappital-manager
      --manager-client-key string    the HTTPS client key file of kappital-manager
      --manager-https-port string    the HTTPS Port of kappital-manager
      --manager-skip-verify          connect to kappital-manager need to skip verify
```

### 3. Deploy the Service into Cluster

```shell
Use a Kappital package to create a Cloud Native Service in a cluster

Usage:
  kappctl create service [Cloud Native Package Directory Path] [flags]

Flags:
  -h, --help   help for service
```

### 4. Deploy the Service Instance into Cluster

```shell
Install an instance of a specific Cloud Native Service

Usage:
  kappctl create instance [service-name] [flags]

Flags:
  -d, --dir string    the Cloud Native Package Path
  -f, --file string   the custom resource file path
  -h, --help          help for instance
```

- If user does not deploy Service into cluster before deploying the Service Instance, it needs to use `-d, --dir` flag. The `Kappital-Manager` will deploy the Service, and then deploy the Service Instance.
- If user does not use `-f, --file` flag, the `-d, --dir` flag is required. `Kappital-Manager` will use the example Service Instance config content to deploy into cluster.
- If user has already deployed the Service into cluster, user can only use the `-f, --file` flag to deploy Service Instance config content into cluster.

### 5. Search the Service

```shell
Query one or many Cloud Native Services

Usage:
  kappctl get service [service-name] [flags]

Aliases:
  service, services, svc

Flags:
  -A, --all              query resources across all repos/clusters. When this flag is set, --repo or --cluster flag will HAS NO EFFECT
  -c, --cluster string   the cluster scope of the Cloud Native Service (default "default")
  -h, --help             help for service
  -o, --output string    the output format of the queried resource, can be yaml or json
```

- Now, the `Kappital-Manager` is only single-cluster function, the `-c, --cluster` is useless. Planning develop multi-cluster version in the future.

### 6. Search the Service Instance

```shell
Query instances of a Cloud Native Instance.

Usage:
  kappctl get instance instance-name [flags]

Aliases:
  instance, instances

Flags:
  -A, --all                query resources across all repos/clusters. When this flag is set, --repo or --cluster flag will HAS NO EFFECT
  -c, --cluster string     the cluster scope of the Cloud Native Service (default "default")
  -h, --help               help for instance
  -n, --namespace string   the namespace of the specified instance (default "default")
  -o, --output string      the output format of the queried resource, can be yaml or json
  -s, --service string     the cloud native service name
```

### 7. Uninstall the Service

```shell
Delete a Cloud Native Service in a cluster

Usage:
  kappctl delete service service-name [flags]

Flags:
  -c, --cluster string   the cluster scope of the Cloud Native Service (default "default")
  -h, --help             help for service
```

- If there has Service Instance belong the Service, the `Kappital-Manager` will uninstall all Service Instance which belong this Service, and then uninstall Service.

### 8. Uninstall the Service Instance

```shell
Delete a Cloud Native Service Instance in a cluster

Usage:
  kappctl delete instance instance-name [flags]

Flags:
  -c, --cluster string   the cluster scope of the Cloud Native Service (default "default")
  -h, --help             help for instance
  -s, --service string   the cloud native service name
```


