#!/bin/bash
echo "==> Install Operator..."
kubectl apply -f https://raw.githubusercontent.com/ansible/awx-operator/devel/deploy/awx-operator.yaml

echo "==> Wait Operator started..."
kubectl wait --for=condition=ready pod -l name=awx-operator

echo "==> Starting AWX Test installation..."
kubectl create ns ansible-awx
# kubectl delete AWX awx -n ansible-awx
KUBENODE_IP=$(kubectl get node -ojson | jq '.items[0].status.addresses[0].address' -r)
INGRESS_DOMAIN="${KUBENODE_IP}.sslip.io"

cat <<EOF | kubectl apply -f -
apiVersion: awx.ansible.com/v1beta1
kind: AWX
metadata:
  name: awx
  namespace: ansible-awx
spec:
  deployment_type: awx
  tower_admin_user: test
  tower_admin_email: test@example.com
  tower_admin_password: changeme
  tower_broadcast_websocket_secret: changeme
  tower_ingress_type: Ingress
  tower_hostname: awx.${INGRESS_DOMAIN}
EOF

echo "==> Waiting Operator Started the AWX Deployment..."
sleep 45

echo "==> Waiting AWX full Started..."
kubectl wait --for=condition=ready pod -l app=awx -n ansible-awx --timeout=800s
