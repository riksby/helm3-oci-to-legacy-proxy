# Helm 3 OCI-based registry proxy with authentication

Listens on port 80 and serves .tar.gz files from the specified [OCI-based Chart registry](https://helm.sh/docs/topics/registries/).

| OCI-based image               | request                                                      |
| ----------------------------- | ------------------------------------------------------------ |
| `ghcr.io/myorg/myimage:1.0.0` | http://helm3-oci-to-legacy-proxy/ghcr.io/myorg/myimage:1.0.0 |
| `myrepo/mychart:2.7.0`        | http://helm3-oci-to-legacy-proxy/myrepo/mychart:2.7.0        |

To authenticate just mount your `$HOME/.docker/config.json` into the container at `/` see below.
For kubernetes you can reuse your `kubernetes.io/dockerconfigjson` imagePullSecret and mount that.

## Image

[ghcr.io/riksby/helm3-oci-to-legacy-proxy](https://github.com/orgs/riksby/packages/container/package/helm3-oci-to-legacy-proxy)

## Usage

```sh
docker run --rm -it -p 8080:80 -v $HOME/.docker/config.json:/.docker/config.json ghcr.io/riksby/helm3-oci-to-legacy-proxy
```
