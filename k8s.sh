export $(grep '^REGISTRY_HOST=' .env)
export $(grep '^HOST=' .env)
export WORKING_DIR="$(pwd)"


if [[ "$1" == "--build" ]]; then
  IMAGE_NAME="${REGISTRY_HOST}/gonotify:latest"
    docker build -t "$IMAGE_NAME" .
    docker push "$IMAGE_NAME"   
fi

envsubst '${REGISTRY_HOST} ${HOST} ${WORKING_DIR}' < k8s.yaml > k8s_temp.yaml

kubectl apply -k .

rm k8s_temp.yaml

