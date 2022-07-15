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
$ podman run --rm -p 3000:3000 -d --name webapp -t quay.io/tkrishtop/webapp:speaker
$ curl localhost:3000
[Speaker] Hi there!
$ podman push quay.io/tkrishtop/webapp:speaker
```

Build listener image

```
$ podman build -f Dockerfile_listener -t quay.io/tkrishtop/webapp:listener .
$ podman push quay.io/tkrishtop/webapp:listener
```

## Deployment on k8s

For the reference: [Install minikube on Fedora36](https://www.tutorialworks.com/kubernetes-fedora-dev-setup/).

```
$ minikube start --driver=kvm2
$ minikube ip
192.168.39.23

$ kubectl apply -f deployment/speaker.yaml
deployment.apps/speaker-deployment created
service/speaker-service created
$ kubectl apply -f deployment/listener-config.yaml
configmap/listener-config created
$ kubectl apply -f deployment/listener.yaml
deployment.apps/listener-deployment created
service/listener-service created

$ kubectl get pod
NAME                                  READY   STATUS    RESTARTS   AGE
listener-deployment-554f69bb5-k4drx   1/1     Running   0          65s
speaker-deployment-559b7954b4-hnggg   1/1     Running   0          66m
speaker-deployment-559b7954b4-j6jxt   1/1     Running   0          66m
speaker-deployment-559b7954b4-rlcnv   1/1     Running   0          66m

$ kubectl get svc
NAME               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
kubernetes         ClusterIP   10.96.0.1       <none>        443/TCP          66m
listener-service   NodePort    10.109.35.179   <none>        3001:30101/TCP   64m
speaker-service    NodePort    10.107.196.68   <none>        3000:30100/TCP   66m

$ curl 192.168.39.23:30100
[Speaker] Hi there!

$ kubectl logs -l app=speaker
2022/07/14 22:06:21 Speaker is active
2022/07/14 22:06:38 Got request, going to speak
2022/07/14 22:06:18 Speaker is active
2022/07/14 22:06:20 Speaker is active

$ kubectl logs -l app=listener
2022/07/14 23:56:41 Sending a request
2022/07/14 23:56:41 Got response:  &{0xc0002fb440 {0 0} false <nil> 0x606e20 0x606f20}
2022/07/14 23:56:42 Sending a request
2022/07/14 23:56:42 Got response:  &{0xc000346e40 {0 0} false <nil> 0x606e20 0x606f20}
2022/07/14 23:56:43 Sending a request

# to restart
$ kubectl get deploy
$ kubectl rollout restart deploy speaker-deployment
```
