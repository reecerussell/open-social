#!/bin/bash

echo "Applying Kubernetes Dashboard"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
kubectl apply -f ./dashboard/admin.yaml

echo "Applying Ingress Controller (GCLOUD)"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.43.0/deploy/static/provider/cloud/deploy.yaml

echo "Applying Kustomize"
kubectl apply -k ./overlays/dev

echo "Done."
