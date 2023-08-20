# Htpasswd Operator

[![Github CI](https://github.com/mowirth/htpasswd-operator/actions/workflows/build.yaml/badge.svg)](https://github.com/mowirth/htpasswd-operator/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/go.flangaapis.com/htpasswd-kubernetes-operator)](https://goreportcard.com/report/github.com/mowirth/htpasswd-operator)

A small operator to automatically generate bcrypt htpasswd files from CRDs.
This allows you to handle your htpasswd accounts as Kubernetes resources, and reuse existing secrets stored in kubernetes.

htpasswd-operator features both a setup for an init container and an operator to automatically watch and deploy changes in the
CRD. If used as an operator, secrets and configmaps referenced by any htpasswd-user are automatically updated when a change occurs.

If you need multiple operators or want to limit the access or exclude some credentials, you can restrict the access to the users for each operator using Cluster Roles and Service Accounts.

The image is available at DockerHub: [mowirth/htpasswd-operator](https://hub.docker.com/r/mowirth/htpasswd-operator)

### Configuration

| Env Variable  | Description                                                                    | Default    |
|---------------|--------------------------------------------------------------------------------|------------|
| HTPASSWD_FILE | Set the output file. The directory must already exist.                         | ./htpasswd |
| KUBE_CONFIG   | Path to the kube config. Can be empty.                                         | ""         |
| PID_FILE      | Path to the PID file. Can be empty.                                            | ""         |
| LOG_LEVEL     | Set Log Level, use TRACE, DEBUG, INFO, WARN, ERROR, FATAL for specific logging | INFO       |

### Setup

An example setup for a docker registry and the operator as an init container can be found in the `kubernetes/`
directory.

##### Custom Resource Definition

The CRD is available at `kubernetes/htpasswd-account.crd.yaml` and must be applied prior to creating any users.

The operator requires sufficient permission to list the CRD resources and the secrets and configmaps from where the
users should be retrieved.

If one of the referenced secrets or configMaps is removed or updated, the htpasswd entry is removed or updated automatically

##### Adding Users

The syntax for the htpasswd user is very similar to other deployments.
We support mixing configMapValues, secretKeyValues and normal values, so you can directly define your username and load
the password from the secret.

```yaml
---
kind: HtpasswdUser
apiVersion: flanga.io/v1
metadata:
  name: htpasswd-admin
spec:
  username:
    value: admin
  password:
    value: adminPassword
```

Additionally, it is possible to retrieve the user from a configMap:

```yaml
kind: HtpasswdUser
apiVersion: flanga.io/v1
metadata:
  name: htpasswd-configuser
spec:
  username:
    configMapKeyRef:
      name: htpasswd-configusers
      key: test-username
  password:
    configMapKeyRef:
      name: htpasswd-configusers
      key: test-password
```

Finally, it is also possible to load your user from a secret:

```yaml
kind: HtpasswdUser
apiVersion: flanga.io/v1
metadata:
  name: htpasswd-testuser1
spec:
  username:
    secretKeyRef:
      name: htpasswd-users
      key: test1-username
  password:
    secretKeyRef:
      name: htpasswd-users
      key: test1-password
```

##### Reload Signal
If you want to send a SIGHUP signal after regenerating the htpasswd file, you can configure a PIDFile over the config.
If a change is written to disk, the operator sends a SIGHUP signal to the process, signaling a reload.

Please note that the reload is debounced by a second to prevent unnecessary restarts during many config changes.

### Modes

We support running htpasswd-operator as an init-container, standalone and as a watcher.

Compared to the watcher mode, the operator exists abnormally if something fails (or normally when the file was created), where the watch image continues to
listen for changes.

##### Standalone

While running htpasswd-operator as standalone is not recommended, since it requires manual work, it is possible to use this tool without kubernetes
deployment.
To do so, deploy the CRD and the user definitions in the required kubernetes context.
Afterward, run the generate command and use the resulting htpasswd file generated in the output directory.
Please note that this requires permissions to access secrets, configMaps and CRDs.

##### As an Init Container

When running as an init-container, all existing user definitions are loaded, resolved into the matching data and stored.
It is worth noting that htpasswd-operator does not successfully exit before all users are loaded to prevent missing
secrets that are loaded later.

##### As a watcher

When running as a watcher, htpasswd-operator is deployed as a sidecar container next to the required application.
When a change on a HtpasswdUser or one of its referenced Secrets or ConfigMaps is detected, the htpasswd file is automatically updated if necessary.
Additionally, it is possible to configure a PIDFile, where the operator can send a `SIGHUP` command to signal the
process that it should reload the config.


