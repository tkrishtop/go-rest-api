# REST API on Golang

## Testing in local

```console
# main terminal
go run main.go

# another terminal
$ curl localhost:3000
[Worker 1] Hi there!
```

![homepage](readme/homepage.png)

## Dockerization

```
# main terminal
$ podman build -t quay.io/tkrishtop/webapp:distroless .
$ podman run -rm -p 3000:3000 -d --name webapp -t quay.io/tkrishtop/webapp:distroless
$ curl localhost:3000
[Worker 1] Hi there!
```

Let's compare ubi and distroless images by size:
```
$ podman images | grep webapp
quay.io/tkrishtop/webapp                                           distroless                                b8e1f29ab611  10 seconds ago  9.94 MB
quay.io/tkrishtop/webapp                                           ubi                                       a797dc42dbb5  17 minutes ago  241 MB
```

The image is pushed in quay.io

```
podman pull quay.io/tkrishtop/webapp:distroless
```

## Deployment on k8s

For the reference: [Install minikube on Fedora36](https://www.tutorialworks.com/kubernetes-fedora-dev-setup/).

```
$ minikube start --driver=kvm2
$ kubectl get node -o wide
NAME       STATUS   ROLES                  AGE   VERSION   INTERNAL-IP      EXTERNAL-IP   OS-IMAGE              KERNEL-VERSION   CONTAINER-RUNTIME
minikube   Ready    control-plane,master   20m   v1.23.3   192.168.39.131   <none>        Buildroot 2021.02.4   4.19.202         docker://20.10.12

$ minikube ip
192.168.39.131

$ kubectl apply -f deployment/webapp.yaml 
deployment.apps/webapp-deployment created
service/webapp-service created

$ kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
webapp-deployment-8b946c8f5-4bszb   1/1     Running   0          18m
webapp-deployment-8b946c8f5-9jw7v   1/1     Running   0          18m
webapp-deployment-8b946c8f5-lpbqx   1/1     Running   0          18m

$ kubectl get svc
NAME             TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
kubernetes       ClusterIP   10.96.0.1     <none>        443/TCP          22m
webapp-service   NodePort    10.111.56.3   <none>        3000:30100/TCP   19m


$ curl 192.168.39.131:30100
[Worker 1] Hi there!
```

