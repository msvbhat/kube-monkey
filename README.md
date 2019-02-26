# kube-monkey

kube-monkey is a tool to check the resiliency of the system. It deletes
random pods repeatedly at specific intervals.

## Description

This is a tool inspired after the
[Chaos Monkey](https://en.wikipedia.org/wiki/Chaos_engineering#Chaos_Monkey).
But this simply kills the random pods in the Kubernetes cluster. There are few
ways to control which pods can be killed and at what intervals etc. Those are
described [below](https://github.com/msvbhat/kube-monkey#deploy).

## Dependencies and Building

The tool is written in go and uses official
[client-go](https://github.com/kubernetes/client-go/) library. Also to expose
the health check API it uses [mux](https://github.com/gorilla/mux) project.

### Installing dependencies

Once you clone the repo, run the below command at the root of the repo.

```bash
make dep
```

This installs all the dependencies of the repo. Note that this project does
not use the dependency management or vendoring yet. So the behaviour might
be different for you. I will add dependency management soon in future.

### Buidling

To build a static binary for linux systems run the below command at the root
if the repository.

```bash
make build
```

This creates a binary called `kube-monkey` at the root of the repository.

## Deploy

This tool is built assuming that it would be running inside the kubernetes
cluster as a pod. So it is important that the contianer is running inside
a kubernetes pod for authenticating with the Kubernetes API server.

And since it discovers and deletes other pods, it needs to be running with
proper `serviceaccounts` with required permissions. To create the required
service accounts with required permissions, run the below command.

Please note that below command assumes that there is a Kubernetes cluster
running and `kubectl` is configured to communicate with the cluster.

```bash
kubectl create -f k8s-deploy/rbac.yaml
```
And then to deploy `kube-monkey` as a deployment, run the below command. And
note that the image is pulled from the docker repo `msvbhat/kube-monkey`. If
you have built another docker image probaly with custom built binary, please
update it in the [file](k8s-deploy/kube-monkey.yaml).

```bash
kubectl create -f k8s-deploy/kube-monkey.yaml
```

By default the 50% of the pods are killed every 2 minutes. The pods running
in `kube-system` namespaces are whitelisted by default. To control this
behaviour, please use the below env variables in the deployment manifest.

1. `NAMESPACE_WHITELIST` - This is a space seperated list of Namespaces that
    should be whitelisted from killing pods. That means the pods running in
    these namespaces will not be considered for deleting. And the namespace
    kube-system is always whitelisted.

2. `DELETE_PERCENTAGE` - This is the percentage of pods that should be
    deleted. To not delete any pod, specify 0 and to delete all pods
    specify 100. But note that this percentage is applied to the pods that
    are eligible for deletion i.e. this percentage is applied to the pods
    that are not running in whitelisted namespaces.

3. `KM_SCHEDULE` - This is the schedule for kube monkey to delete pods. This
    follows the cron syntax. To understand more about the cron syntax that is
    allowed, please check
    [docs](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format)

## Considerations and Limitations

This has been only tested with the minikube. But is supposed to run in any
Kubernetes cluster.

The project doesn't have unit tests yet. Unit tests will be added soon.

Currently the /metrics endpoint is a dummy endpoint. It doesn't return any
metrics but only returns 200 OK.

## Planned Enhancements

1. Introdure Active Health Check instead of Passive one
1. Define metrics for exporting and add metrics endpoint
1. Add sophisticated method of specifying pods with labels etc
1. Also add blacklisting namespaces
1. Use cli args instead of env variables
1. Send events to pods for visibility
