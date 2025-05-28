# Kubernetes Authorization Listing - KAL

KAL can be used to list every permission of a Kubernetes user, service account token, kubeconfig authentication, or a JWT token.

This CLI connects to the provided Kubernetes Cluster, list all resources, and for each resource tests if the provided authentication has access in the resource. The test is performed using the [SelfSubjectAccessReview](https://kubernetes.io/docs/reference/kubernetes-api/authorization-resources/self-subject-access-review-v1/) request.

## Installation

### Go Install

```sh
go install -v github.com/ing-bank/kal@latest
```

### Compile from source

```sh
git clone https://github.com/ing-bank/kal.git
cd kal; go install
```

## Quick Start

### User authentication options

#### 1. Automatic

KAL searches for authentication credentials in the following order:

1. Provided in `-token` argument
2. Search for a kubeconfig file (default location `~/.kube/config`)
3. Assume it is running inside a POD and using the credentials in the `/var/run/secrets/kubernetes.io/serviceaccount/` folder

#### 2. Manual authentication

Provide the authentication token as a CLI argument.

```sh
kal -token '<your_jwt_token>'
```

#### 3. Custom kubeconfig location

Provide the custom kubeconfig file location.

```sh
kal -c /path/to/kubeconfig.yaml
```


### Execution

#### 1. Listing permissions of default namespace

Command:

```sh
kal
```

Expected output:
```sh
############################
#                          #
# ██╗  ██╗ █████╗ ██╗      #
# ██║ ██╔╝██╔══██╗██║      #
# █████╔╝ ███████║██║      #
# ██╔═██╗ ██╔══██║██║      #
# ██║  ██╗██║  ██║███████╗ #
# ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝ #
# Kubernetes Authz Listing #
############################

[!] legal disclaimer: Usage of kal for attacking targets without prior mutual consent is illegal. It is the end user\'s responsibility to obey all applicable local, state and federal laws. Developers assume no liability and are not responsible for any misuse or damage caused by this program


[INF] running from namespace = default
[INF] found 105 resources and sub-resources
bindings/v1 [create,get,bind,patch,escalate,deletecollection,list,impersonate,watch,update,delete,approve] [default]
componentstatuses/v1 [create,get,delete,deletecollection,escalate,impersonate,update,patch,approve,watch,bind,list] [CLUSTER_WIDE]
...[snip]...
prioritylevelconfigurations.flowcontrol.apiserver.k8s.io/v1 [escalate,impersonate,list,approve,watch,deletecollection,get,patch,update,delete,bind,create] [CLUSTER_WIDE]
prioritylevelconfigurations.flowcontrol.apiserver.k8s.io/v1/status [escalate,impersonate,patch,watch,list,create,get,delete,update,approve,deletecollection,bind] [CLUSTER_WIDE]
flowschemas.flowcontrol.apiserver.k8s.io/v1beta3 [patch,approve,create,escalate,list,deletecollection,impersonate,delete,watch,update,bind,get] [CLUSTER_WIDE]
flowschemas.flowcontrol.apiserver.k8s.io/v1beta3/status [escalate,patch,deletecollection,update,get,bind,impersonate,delete,approve,watch,list,create] [CLUSTER_WIDE]
prioritylevelconfigurations.flowcontrol.apiserver.k8s.io/v1beta3 [escalate,impersonate,approve,update,get,create,list,deletecollection,patch,watch,delete,bind] [CLUSTER_WIDE]
prioritylevelconfigurations.flowcontrol.apiserver.k8s.io/v1beta3/status [get,create,list,escalate,impersonate,patch,bind,update,delete,approve,watch,deletecollection] [CLUSTER_WIDE]
```

#### 2. Custom namespace

```sh
kal -namespace <namespace>
```

#### 3. No Rate Limit

Removes the rate limit restraints enforced by `k8s.io/client-go/kubernetes` package.

```sh
kal -no-rate-limit
```

#### 4. List permissions with User Impersonation

Impersonate a user and list its permissions.

```sh
kal -as '<user>'
```

### Output Options

#### Verbose & Silent

Select the verbosity of the output.

```sh
kal -verbose/-silent
```

#### Show all results

This option show all results, even not allowed commands.

```sh
kal -all
```

#### JSON output

```sh
kal -json
```

#### Show permission reason

Command: 
```sh
kal -show-reason
```

Expected output:

```sh
[ERR] could not create a kubernetes custom client error=invalid configuration for kubernetes custom client
[INF] running from namespace = default
[INF] found 105 resources and sub-resources
bindings/v1 [delete,patch,bind,create,update,watch,get,list,deletecollection,impersonate,approve,escalate] [default] [RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins"]
componentstatuses/v1 [get,escalate,list,delete,approve,patch,update,bind,watch,impersonate,deletecollection,create] [CLUSTER_WIDE] [RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins"]
...[snip]...
prioritylevelconfigurations.flowcontrol.apiserver.k8s.io/v1beta3 [create,patch,update,deletecollection,escalate,get,delete,bind,watch,impersonate,list,approve] [CLUSTER_WIDE] [RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins"]
prioritylevelconfigurations.flowcontrol.apiserver.k8s.io/v1beta3/status [create,escalate,list,update,delete,deletecollection,bind,patch,get,approve,watch,impersonate] [CLUSTER_WIDE] [RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins";RBAC: allowed by ClusterRoleBinding "kubeadm:cluster-admins" of ClusterRole "cluster-admin" to Group "kubeadm:cluster-admins"]
```

## Internals

This section explains how KAL works under the hood.

### Inspiration

Based on the article of [Raesene - Fun with Kubernetes Authorization Auditing](https://raesene.github.io/blog/2024/04/22/Fun-with-Kubernetes-Authz/), sometimes the command `kubectl auth can-i --list` can omit some permissions specially if they are from a custom resource. In this case, KAL overcomes this "issue" by listing all available resources and testing if the current authorization has permission to execute certain API verb in the resource.


#### API Verbs

[Kuberntes Authorization Request Verbs](https://kubernetes.io/docs/reference/access-authn-authz/authorization/#determine-the-request-verb)

- create
- get
- list
- watch
- update
- patch
- delete
- deletecollection
- impersonate
- bind
- approve
- escalate

#### Api Resources

Listing all API resources.

```sh
kubectl auth can-i --list -o wide
```

## Contributing

Contributions are more than welcome! Please see our [contribution guidelines first](./CONTRIBUTING.md).

## Use as a library

KAL can be used a a library by instantiating the [pkg/runner](./pkg/runner/) package, it contains the required setup.

```go
import "github.com/ing-bank/kal/pkg/runner"

func main() {
    kalRunner := runner.FromOptions()
}
```

## License

You can check our licensing scheme [here](./LICENSE).
