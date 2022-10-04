# REST API on Golang

## Idea

Listener is a web server and could send client requests to all known speakers when called as `/speakers`.
It collects all replies and return their concatenation.

TODO:

- Do it concurrently using goroutines inside listener.

- Do not mind about speaker dying and ignore the absence of reply.

## Dockerization

Side note: distroless vs ubi: 9.94 MB vs 241 MB
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
[unknown] said: UnknownSpeech
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

$ kubectl apply -f deployment/speaker-config.yaml
configmap/speaker-config created

$ kubectl apply -f deployment/speaker-winnie.yaml 
deployment.apps/winnie-deployment created
service/winnie-service created
$ kubectl apply -f deployment/speaker-piglet.yaml 
deployment.apps/piglet-deployment created
service/piglet-service created

$ kubectl get pod
NAME                                 READY   STATUS    RESTARTS   AGE
piglet-deployment-7fd7dc4689-8t4m9   1/1     Running   0          11s
winnie-deployment-b75885d67-2s6hh    1/1     Running   0          16s

$ kubectl logs -l app=speaker
2022/10/04 10:43:33 piglet is active
2022/10/04 10:43:30 winnie is active

$ kubectl apply -f deployment/listener-config.yaml 
configmap/listener-config created
$ kubectl apply -f deployment/listener.yaml 
deployment.apps/listener-deployment created
service/listener-service created

$ kubectl logs -l app=listener
2022/10/04 10:45:25 Listener is active

$ kubectl get svc
NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubernetes         ClusterIP   10.96.0.1        <none>        443/TCP          2m53s
listener-service   NodePort    10.108.218.193   <none>        3003:30103/TCP   31s
piglet-service     NodePort    10.105.89.198    <none>        3002:30102/TCP   2m20s
winnie-service     NodePort    10.103.46.135    <none>        3000:30100/TCP   2m25s

# Talk to all speakers
# Collect and return concat of replies
$ curl $(minikube ip):30103/speakers
Concatenated replies: 
[winnie] said: Hallo, Rabbit, isn't that you? 
[piglet] said: Isn't that Rabbit's voice? 

$ kubectl logs -l app=listener
2022/10/04 10:45:25 Listener is active
2022/10/04 10:46:08 [listener] Got a list of URLs [winnie-service piglet-service]
2022/10/04 10:46:08 [listener] Calling URL: http://winnie-service:3000
2022/10/04 10:46:08 [listener] Calling URL: http://piglet-service:3002

$ kubectl logs -l app=speaker
2022/10/04 10:43:30 winnie is active
2022/10/04 10:46:08 winnie got a request, going to tell: Hallo, Rabbit, isn't that you?
2022/10/04 10:43:33 piglet is active
2022/10/04 10:46:08 piglet got a request, going to tell: Isn't that Rabbit's voice?

# to restart
$ kubectl get deploy
$ kubectl rollout restart deploy listener-deployment
```
