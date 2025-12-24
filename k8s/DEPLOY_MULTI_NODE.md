# Deploy to Multi-Node Kubernetes Cluster

## Option 1: Push to Docker Hub (Recommended)

```bash
# Tag the image
docker tag todo-microservice:latest YOUR_DOCKERHUB_USERNAME/todo-microservice:latest

# Login to Docker Hub
docker login

# Push the image
docker push YOUR_DOCKERHUB_USERNAME/todo-microservice:latest
```

Update `k8s/deployment.yaml` image to:
```yaml
image: YOUR_DOCKERHUB_USERNAME/todo-microservice:latest
```

Then deploy:
```bash
kubectl apply -f k8s/deployment.yaml
```

## Option 2: Load Image on All Nodes

For kind cluster:
```bash
kind load docker-image todo-microservice:latest --name YOUR_CLUSTER_NAME
```

For other multi-node setups, save and load on each node:
```bash
# Save image
docker save todo-microservice:latest > todo-microservice.tar

# Copy to each worker node and load
# On each worker node:
docker load < todo-microservice.tar
```

## Option 3: Use Local Registry

If your cluster is kind/minikube, you can set up a local registry that all nodes can access.
