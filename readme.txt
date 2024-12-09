go mod init dummy-post-office-server

go mod tidy

#Makefile
docker build -t golang-app .
docker tag golang-app localhost:32000/dummy-post-office-golang-app:latest
docker push localhost:32000/dummy-post-office-golang-app:latest
microk8s kubectl rollout restart deploy dummy-post-office-golang-app -n dummy-post-office-namespace

sudo vim /etc/hosts
127.0.0.1       post.mydomain.com



# desde el pod
curl http://localhost:8080/status
# ip del pod
curl http://10.1.69.40:8080/status
# ip del servicio
microk8s kubectl get services -n dummy-post-office-namespace | grep dummy-post-office
curl http://10.152.183.94:8080/status
# nombre del servicio corto
curl http://dummy-post-office-golang-app-service:8080/status
# nombre del servicio largo
curl http://dummy-post-office-golang-app-service.dummy-post-office-namespace.svc.cluster.local:8080/status
# desde fuera
curl -k https://post.mydomain.com/post-office/status
# desde el clúster
microk8s kubectl run curlpod --image=curlimages/curl:latest -it --rm -- /bin/sh
curl http://dummy-post-office-golang-app-service:8080/status