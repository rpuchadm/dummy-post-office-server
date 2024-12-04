go mod init dummy-post-office-server

go mod tidy

docker build -t golang-app .

docker tag golang-app localhost:32000/dummy-post-office-golang-app:latest

docker push localhost:32000/dummy-post-office-golang-app:latest

microk8s kubectl rollout restart deploy dummy-post-office-golang-app-container -n dummy-post-office-namespace