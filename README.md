# Kuberneters-native Golang web application with concurrency

## Idea

Listener is a web server that could send client requests to many speakers. Being called at `/speakers` endpoint, it collects the replies from all speakers listed in deployment/listener-config.yaml and returns their concatenation.

![flowchart](images/flowchart.png)

TODO:

- Implement global context to gracefully stop the application, handle errors and timeouts 


## Dockerization

Side note: distroless 9.94 MB vs ubi 241 MB
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

$ kubectl logs -l app=listener -f
2022/10/08 15:11:11 Going to read config file:  /config/speakers-config.yaml
2022/10/08 15:11:11 Going to unmarchall configuration
2022/10/08 15:11:11 List of speakers to call:  [{winnie winnie-service :3000} {piglet piglet-service :3002}]
2022/10/08 15:11:11 Calling speaker winnie
2022/10/08 15:11:11 Calling speaker piglet
2022/10/08 15:11:11 Calling URL: http://piglet-service:3002
2022/10/08 15:11:11 Calling URL: http://winnie-service:3000
2022/10/08 15:11:11 There is an error while calling url, ignore it and return empty reply: Get "http://piglet-service:3002": dial tcp 10.110.191.158:3002: connect: connection refused

$ curl $(minikube ip):30103/speakers
Concatenated replies: 
[winnie] said: Hallo, Rabbit, isn't that you? 
[piglet] said: Isn't that Rabbit's voice?

$ kubectl logs -l app=listener -f
2022/10/08 15:11:19 Going to read config file:  /config/speakers-config.yaml
2022/10/08 15:11:19 Going to unmarchall configuration
2022/10/08 15:11:19 List of speakers to call:  [{winnie winnie-service :3000} {piglet piglet-service :3002}]
2022/10/08 15:11:19 Calling speaker winnie
2022/10/08 15:11:19 Calling speaker piglet
2022/10/08 15:11:19 Calling URL: http://piglet-service:3002
2022/10/08 15:11:19 Calling URL: http://winnie-service:3000

# to restart
$ kubectl get deploy
$ kubectl rollout restart deploy listener-deployment
```
