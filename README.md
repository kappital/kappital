# Kappital
[![Build Status]()]()
[![Go Report Card]()]()
[![LICENSE]()](/LICENSE)
[![Releases]()]()
[![Documentation Status]()]()
[![CII Best Practices]()]()

<img src="./docs/images/kappital-logo-only.png">

Kappital is an open source project that enables developers to manage cloud native applications across multiple clouds and edges
with no changes to developer's applications. Kappital defines Cloud Native Service Package specifications. By complying with spec,
Kappital enhances service capabilities, observability to applications.

## Why Kappital
* **Unified Cloud Native Service Management**
  - Web console and command-line client for centrally manages multiple clusters
  - Full lifecycle management,such as install,upgrade,update,unInstall,state management
* **Declarative Observability With Non-Intrusive**
  - Zero-Code for logging,monitoring and alarm
  - Support CNCF Observability project,such as prometheus, OpenTelemetry,Thanos
* **Multi-Deployment Engine**
  - Support Kubernetes Native deployment
  - Support Helm,Operator-Framework, and other deploy framework

### Architecture

<div  align="center">
<img src="./docs/images/kubeedge_arch.png" width = "85%" align="center">
</div>

The Kappital consists of the following components:
* Catalog
  - Directory Cloud-Native Service Package. Support OCI storage format.
  - Support HelmChart, Operator-Framework Bundle.
* Manager
  - Support full lifecycle of Cloud-Native Service such as install, upgrade, update and uninstall.

* Service-Engine

## Guides

Get start with this [doc](https://kubeedge.io/en/docs).

See our documentation on [kubeedge.io](https://kubeedge.io) for more details.

To learn deeply about KubeEdge, try some examples on [examples](https://github.com/kubeedge/examples).

## Roadmap

* [2021 Roadmap](./docs/roadmap.md#roadmap)

## Meeting

Regular Community Meeting:
- Europe Time: **Wednesdays at 16:30-17:30 Beijing Time** (biweekly, starting from Feb. 19th 2020).
  ([Convert to your timezone.](https://www.thetimezoneconverter.com/?t=16%3A30&tz=GMT%2B8&))
- Pacific Time: **Wednesdays at 10:00-11:00 Beijing Time** (biweekly, starting from Feb. 26th 2020).
  ([Convert to your timezone.](https://www.thetimezoneconverter.com/?t=10%3A00&tz=GMT%2B8&))

Resources:
- [Meeting notes and agenda](https://docs.google.com/document/d/1Sr5QS_Z04uPfRbA7PrXr3aPwCRpx7EtsyHq7mp6CnHs/edit)
- [Meeting recordings](https://www.youtube.com/playlist?list=PLQtlO1kVWGXkRGkjSrLGEPJODoPb8s5FM)
- [Meeting link](https://zoom.us/j/4167237304)
- [Meeting Calendar](https://calendar.google.com/calendar/embed?src=8rjk8o516vfte21qibvlae3lj4%40group.calendar.google.com) | [Subscribe](https://calendar.google.com/calendar?cid=OHJqazhvNTE2dmZ0ZTIxcWlidmxhZTNsajRAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)

## Contact

If you need support, start with the [troubleshooting guide](https://kubeedge.io/en/docs/developer/troubleshooting), and work your way through the process that we've outlined.

If you have questions, feel free to reach out to us in the following ways:

- [mailing list](https://groups.google.com/forum/#!forum/kubeedge)
- [slack](https://join.slack.com/t/kubeedge/shared_invite/enQtNjc0MTg2NTg2MTk0LWJmOTBmOGRkZWNhMTVkNGU1ZjkwNDY4MTY4YTAwNDAyMjRkMjdlMjIzYmMxODY1NGZjYzc4MWM5YmIxZjU1ZDI)
- [twitter](https://twitter.com/kubeedge)

## Contributing

If you're interested in being a contributor and want to get involved in
developing the KubeEdge code, please see [CONTRIBUTING](./CONTRIBUTING.md) for
details on submitting patches and the contribution workflow.

## Security

We encourage security researchers, industry organizations and users to proactively report suspected vulnerabilities to our security team (`cncf-kubeedge-security@lists.cncf.io`), the team will help diagnose the severity of the issue and determine how to address the issue as soon as possible.

For further details please see [Security Policy](https://github.com/kubeedge/community/blob/master/security-team/SECURITY.md) for our security process and how to report vulnerabilities.

## License

KubeEdge is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.