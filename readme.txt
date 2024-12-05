go mod init dummy-post-office-server

go mod tidy

docker build -t golang-app .

docker tag golang-app localhost:32000/dummy-post-office-golang-app:latest

docker push localhost:32000/dummy-post-office-golang-app:latest

microk8s kubectl rollout restart deploy dummy-post-office-golang-app -n dummy-post-office-namespace

sudo vim /etc/hosts
127.0.0.1       post.mydomain.com