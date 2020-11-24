echo "deploying atlant service..."
kubectl apply -f ./atlant-deployment.yaml
kubectl apply -f ./atlant-service.yaml

echo "deploying envoy service..."
kubectl apply -f ./envoy-service.yaml
EXTERNAL_IP=
while [ -z "$EXTERNAL_IP" ]
do
  EXTERNAL_IP=$(kubectl get service envoy -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
done
echo "Assigned external IP: $EXTERNAL_IP"

echo "generating envoy certificates..."
openssl req -x509 -nodes -newkey rsa:2048 \
  -days 365 -keyout ./privkey.pem -out ./cert.pem -subj \
  "/CN=$(kubectl get service envoy -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"

kubectl create secret tls envoy-certs \
    --key privkey.pem --cert cert.pem \
    --dry-run -o yaml | kubectl apply -f -

echo "deploying envoy service..."
kubectl apply -f ./envoy-configmap.yaml
kubectl apply -f ./envoy-deployment.yaml
echo "deployment done"