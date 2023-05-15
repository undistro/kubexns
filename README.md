# kubexns

[![Go Reference](https://pkg.go.dev/badge/github.com/undistro/kubexns.svg)](https://pkg.go.dev/github.com/undistro/kubexns)
[![test](https://github.com/undistro/kubexns/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/undistro/kubexns/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/undistro/kubexns)](https://goreportcard.com/report/github.com/undistro/kubexns)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/undistro/kubexns?sort=semver&color=brightgreen)
![GitHub](https://img.shields.io/github/license/undistro/kubexns?color=brightgreen)

Kubexns (short for "Kubernetes Cross Namespaces") is a container solution 
that enables the mapping of `ConfigMaps` or `Secrets` from different namespaces 
as volumes in Kubernetes Pods using an `initContainer`.

## Why?

By default, Kubernetes restricts Pods to mount `ConfigMaps` or `Secrets` within the same namespace.

## Usage

You can find a complete example in the [example.yaml](examples/example.yaml) file.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
  namespace: app
  labels:
    app: myapp
spec:
  serviceAccountName: myapp     # it must have permission to `get` and `list` `configmaps` and `secrets`
  volumes:
    - name: global-config       # shared volume between init and application container
      emptyDir: {}
  initContainers:
    - name: global-config
      image: ghcr.io/undistro/kubexns:v0.1.0
      volumeMounts:
        - mountPath: "/tmp/.config"
          name: global-config
      env:
        - name: DIR
          value: "/tmp/.config"
        - name: CONFIGMAPS
          value: "config/global-config"     # mount the ConfigMap `global-config` from `config` namespace
        - name: SECRETS_SELECTOR
          value: "foo=bar"                  # match secrets by label selector
  containers:
    - name: app
      image: bash:latest
      imagePullPolicy: IfNotPresent
      command: ["watch"]
      args: ["ls", "-lha", "/tmp/.config"]
      volumeMounts:
        - mountPath: "/tmp/.config"
          name: global-config
  restartPolicy: Always
```

### Environment variables

| name                  | description                                                                                                                                           | default |
|-----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `DIR`                 | The directory path where the files should be placed.                                                                                                  | `/tmp`  |
| `DEFAULT_MODE`        | The mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. | `0644`  |
| `CONFIGMAPS`          | A comma-separated list of `ConfigMaps` namespaced names (`ns1/cm,ns2/cm`)                                                                             | -       |
| `SECRETS`             | A comma-separated list of `Secrets` namespaced names (`ns1/sec,ns2/sec`)                                                                              | -       |
| `CONFIGMAPS_SELECTOR` | A label selector to match `ConfigMaps` (`foo=bar`)                                                                                                    | -       |
| `SECRETS_SELECTOR`    | A label selector to match `Secrets` (`foo=bar`)                                                                                                       | -       |
| `IGNORE_NOT_FOUND`    | Specifies when the not found errors should be ignored                                                                                                 | 'false' |

# Contributing

We appreciate your contribution.
Please refer to our [contributing guideline](https://github.com/undistro/kubexns/blob/main/CONTRIBUTING.md) for further information.
This project adheres to the Contributor Covenant [code of conduct](https://github.com/undistro/kubexns/blob/main/CODE_OF_CONDUCT.md).

# License

Kubexns is available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
