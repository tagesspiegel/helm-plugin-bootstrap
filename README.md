# helm-plugin-bootstrap

A plugin for [Helm](https://helm.sh/) to add additional files to your Helm chart. This plugin adds resources such as `PodDisruptionBudget`, `NetworkPoliciy`, `ServiceMonitor` in the same way as the other files are created during helm create.

## Installation

```bash
helm plugin install https://github.com/tagesspiegel/helm-plugin-bootstrap
```

## Usage

```bash
helm create mychart
helm bootstrap ./mychart
```
