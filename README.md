GitHook is a kubernetes CRDs (kubebuilder based) that can trigger tekton pipeline from git webhooks.

Supported git provider:
- Gogs
- Github
- Gitlab

## Prerequisite
- Kubernetes cluster (tested on 1.14, 1.15)
  See how to setup cluster [here](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)
- Knative serving (istio based) (tested on 0.6)
  See instruction [here](https://knative.dev/docs/install)
  > Note: Knative service endpoint must be accessible from git webhook
- Tekon build pipeline (tested on 0.4)
  See instruction [here](https://github.com/tektoncd/pipeline/blob/master/docs/install.md)

## Installation
- Install crds and service account needed to run the pipeline using command line
```sh
kubectl apply -f https://gitlab.com/pongsatt/githook/-/jobs/244356986/artifacts/raw/release.yaml
kubectl apply -f https://gitlab.com/pongsatt/githook/raw/master/config/tektonrole.yaml
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
In this sample, we will apply githook resource for gitlab. When push event happen to the sample project, it will trigger a simple tekton pipeline which just print a message to the log. See more advance example pipeline [here](https://github.com/tektoncd/pipeline/tree/master/examples).

- Install a simple tekton pipeline
```sh
kubectl apply -f https://gitlab.com/pongsatt/githook/raw/master/config/samples/0-simple_tekton_pipeline.yaml
```
- Create a sample gitlab project
- Create gitlab access token for the controller to registry webhook (following this [instruction](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html))
- Apply below git secret to the cluster. Replace "xxxx" with token from step above

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
- Apply githook resource into cluster (replace projectUrl with your own project from above)
```yaml
apiVersion: tools.pongzt.com/v1alpha1
kind: GitHook
metadata:
  name: githook-sample
spec:
  gitProvider: gitlab
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
- Wait for a while and verify if the webhook is registered on your project ([gitlab webhook](https://docs.gitlab.com/ee/user/project/integrations/webhooks.html))
- Push something to your repository and run command below to verify pipeline
```sh
kubectl get pipelinerun
```
You should see something like:
```sh
NAME                   SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
githook-sample-29ldn   True        Succeeded   20h         20h
```

## How it works
- A new GitHook resource is applied to the cluster
- Controller creates new knative service to receive git webhook and wait until it is ready
- Controller registers a webhook to git repository specified in GitHook resource
- When an event specified in GitHook resource happens, knative service will create new pipelinerun based on spec in GitHook resource
  > Note: Pipeline resource named "git-source" is injected by service using webhook information