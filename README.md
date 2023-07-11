# Htpasswd Basic Auth Kubernetes Operator

A small project to automatically generate bcrypt htpasswd files from CRDs.
This allows user to handle their htpasswd accounts as Kubernetes resources, and reuse existing secrets without manual
interaction.

It features both a setup for an init container and an operator to automatically watch and deploy changes in the
CRD. Please note that this does not include the used secrets and configMaps yet.

Furthermore, at the moment, users are added based on their namespace, so all users created in on namespace are included
in the htpasswd file.
If you need multiple operators, please use different namespaces for your deployments.

Finally, please note that this project is still in development.
Please open a new issue if you encounter any bugs.

### Configuration

| Env Variable  | Description                                                             | Default    |
|---------------|-------------------------------------------------------------------------|------------|
| HTPASSWD_FILE | Set the output file. The directory must already exist.                  | ./htpasswd |
| KUBE_CONFIG   | Path to the kube config. Can be empty.                                  | ""         |
| PID_FILE      | Path to the PID file. Can be empty.                                     | ""         |
| LOG_LEVEL     | Set Log Level, use DEBUG, INFO, WARN, ERROR, FATAL for specific logging | INFO       |

### Setup

An example setup with a docker registry and the operator as an init container can be found in the `kubernetes/`
directory.

##### Custom Resource Definition

The CRD is available at `kubernetes/htpasswd-account.crd.yaml` and must be applied prior to creating any users.

The operator requires sufficient permission to list the CRD resources and the secrets and configmaps from where the
users should be retrieved.

##### Adding Users

The syntax for the htpasswd user is very similar to other deployments.
We support mixing configMapValues, secretKeyvalues and normal values, so you can directly define your username and load
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
---
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

### Modes

We support running htpasswd-operator as an init-container, standalone and as a watcher.

For standalone and init-containers, we recommend to use the generator image: [mowirth/htpasswd-generator]().
Compared to the watcher image, the container exists abnormally if something fails (or normally when the file was created), where the watch image continues to
listen for changes.

The watcher image is available at DockerHub: [mowirth/htpasswd-operator]()


##### Standalone

While this is not recommended, as it still requires manual work, it is possible to use this tool without kubernetes
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
When a change on a HtpasswdUser is detected, the htpasswd file is automatically updated if necessary.
Additionally, it is possible to configure a PIDFile, where the operator can send a `SIGHUP` command to signal the
process that it should reload the config (very useful for mosquitto).


