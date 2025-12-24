docker build -t todo-microservice:latest .
docker tag todo-microservice:latest YOUR_DOCKER_USERNAME/todo-microservice:latest
docker push YOUR_DOCKER_USERNAME/todo-microservice:latest

kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

kubectl get pods -n todo-microservice
kubectl get svc -n todo-microservice
