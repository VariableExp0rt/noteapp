#!/bin/sh
set -o errexit

start_time=$(date +%s)

# create registry container unless it already exists
reg_name='registry'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name kind-dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    listenAddress: "127.0.0.1"
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    listenAddress: "127.0.0.1"
    protocol: TCP
  extraMounts:
  - hostPath: /dev
    containerPath: /dev

EOF

# connect the registry to the cluster network
# (the network may already be connected)
docker network connect "kind" "${reg_name}" || true

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

# Install tektoncd latest release
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml

# Create ConfigMap for Tektoncd requests for storage
kubectl create configmap config-artifact-pvc \
                         --from-literal=size=10Gi \
                         --from-literal=storageClassName=manual \
                         -o yaml -n tekton-pipelines \
                         --dry-run=client | kubectl replace -f -

PATHVAR=${TEKTON_MANIFESTS}
FILES=$(find ${PATHVAR} -name \*.yaml)
for f in $FILES
do
  kubectl apply -f $f || exit
done


# Start the taskrun now that the other resources are in place
tkn task start build-docker-image-from-git-source --use-param-defaults -i docker-source=noteapp-git -o builtImage=noteapp-image --dry-run | kubectl create -f -

end=$(date +%s)
printf "Execution time was `expr $end - $start` seconds. Use the following command to monitor your build;\n->\ttkn taskrun logs --last -f\n\n"
printf "Optional: if you would like to verify that your build was successful and pushed to the repository specified in the Task output, use the following command;\n"
printf "docker run -it -p 8080:8080 --name registry-web --link registry -e REGISTRY_URL=http://registry:5000/v2 -e REGISTRY_NAME=localhost:5000 hyper/docker-registry-web"