GitHook is a kubernetes CRDs (kubebuilder based) that can trigger tekton pipeline from git webhooks.

Supported git provider:
- Gogs
- Gitlab (soon)
- Github (soon)

## Prerequisite
- Kubernetes cluster (tested on 1.14, 1.15)
  See how to setup cluster [here](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)
- Knative serving (istio based) (tested on 0.6)
  See instruction [here](https://knative.dev/docs/install)
  > Note: Knative service endpoint must be accessible from git webhook
- Tekon build pipeline (tested on 0.4)
  See instruction [here](https://github.com/tektoncd/pipeline/blob/master/docs/install.md)

## Installation
- Install crds using command line
```sh
kubectl apply -f https://gitlab.com/pongsatt/githook/-/jobs/243861313/artifacts/raw/release.yaml
```
- Verify if controller is running
```sh
kubectl -n githook-system get pod
```

You should get output like this.
```sh
NAME                                          READY   STATUS    RESTARTS   AGE
githook-controller-manager-7869dc5b76-7gsrx   2/2     Running   0          42m
```

## Sample
- Create git accessToken so the program is able to registry a webhook
- Install a sample tekton pipeline
```sh
kubectl apply -f config/samples/0-simple_tekton_pipeline.yaml
```
- Create git secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: gitsecret
type: Opaque
stringData:
  accessToken: "xxxx" #replace this with your repository access token
  secretToken: "mysecret1234" #replace this with your own value
```
- Install githook resource into cluster
```yaml
apiVersion: tools.pongzt.com/v1alpha1
kind: GitHook
metadata:
  name: githook-sample
spec:
  gitProvider: gogs
  eventTypes:
  - push
  - issue_comment
  - pull_request
  projectUrl: "https://gitlab.com/<your test path>/<your projectname>" #replace this with your test repo
  accessToken:
    secretKeyRef:
      name: gitsecret
      key: accessToken
  secretToken:
    secretKeyRef:
      name: gitsecret
      key: secretToken
  runspec:
    pipelineRef:
      name: simple-pipeline
    serviceAccount: default
```

## How it works
- A new GitHook resource is applied to the cluster
- Controller creates new knative service to receive git webhook and wait until it is ready
- Controller registers a webhook to git repository specified in GitHook resource
- When an event specified in GitHook resource happens, knative service will create new pipelinerun based on spec in GitHook resource
  > Note: Pipeline resource named "git-source" is injected by service using webhook information