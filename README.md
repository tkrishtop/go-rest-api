# REST API on Golang

## Testing in local

```console
# main terminal
$ go run pkg/speaker/speaker.go

# another terminal
$ curl localhost:3000
[Speaker] Hi there!
```

## Dockerization

Side note: distroless vs ubi: 241 MB vs 9.94 MB
```
$ podman images | grep webapp
quay.io/tkrishtop/webapp                                           distroless                                b8e1f29ab611  10 seconds ago  9.94 MB
quay.io/tkrishtop/webapp                                           ubi                                       a797dc42dbb5  17 minutes ago  241 MB
```

Build speaker image

```
$ podman build -f Dockerfile_speaker -t quay.io/tkrishtop/webapp:speaker .
$ podman run -rm -p 3000:3000 -d --name webapp -t quay.io/tkrishtop/webapp:speaker
$ curl localhost:3000
[Speaker] Hi there!
$ podman push quay.io/tkrishtop/webapp:speaker
```

## Deployment on k8s

For the reference: [Install minikube on Fedora36](https://www.tutorialworks.com/kubernetes-fedora-dev-setup/).

```
$ minikube start --driver=kvm2
$ minikube ip
192.168.39.2

# deploy speaker
$ kubectl apply -f deployment/webapp_speaker.yaml 
deployment.apps/speaker-deployment created
service/speaker-service created

$ kubectl get pod
NAME                                  READY   STATUS    RESTARTS   AGE
speaker-deployment-6b6c7f679d-6pg8j   1/1     Running   0          13s
speaker-deployment-6b6c7f679d-dscb4   1/1     Running   0          13s
speaker-deployment-6b6c7f679d-j84gp   1/1     Running   0          13s

$ kubectl get svc
NAME              TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
kubernetes        ClusterIP   10.96.0.1       <none>        443/TCP          73s
speaker-service   NodePort    10.109.33.238   <none>        3000:30100/TCP   34s

$ curl 192.168.39.2:30100
[Speaker] Hi there!

$ kubectl logs -l app=speaker
2022/07/14 22:06:21 Speaker is active
2022/07/14 22:06:38 Got request, going to speak
2022/07/14 22:06:18 Speaker is active
2022/07/14 22:06:20 Speaker is active

# to restart
$ kubectl get deploy
$ kubectl rollout restart deploy speaker-deployment
```

## Adding listener

```
$ go run pkg/listener/listener.go 
2022/07/15 00:27:47 Listener is active
2022/07/15 00:27:47 Sending a request
2022/07/15 00:27:47 Got response:  &{0xc0000ec280 {0 0} false <nil> 0x60aea0 0x60afa0}
2022/07/15 00:27:48 Sending a request
2022/07/15 00:27:48 Got response:  &{0xc00020e000 {0 0} false <nil> 0x60aea0 0x60afa0}
```
