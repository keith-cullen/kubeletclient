# Kubeletclient

A library for accessing kubelet instances over Unix socket and HTTPS interfaces.

## HTTPS Test Instructions

1. Get the name of the node

        $ kubectl get nodes
        save the node name to http/node.txt

2. Create a temporary pod and a service account with appropriate RBAC permissions

        $ docker run -d -p 5000:5000 --restart always --name registry registry:2
        $ docker pull alpine:latest
        $ docker tag alpine:latest localhost:5000/alpine:latest
        $ docker push localhost:5000/alpine:latest
        $ curl http://localhost:5000/v2/alpine/tags/list
        $ kubectl apply -f http/kubeletclient-pod.yaml

3. Get the access token for the service account from the temporary pod

        $ kubectl exec -n kube-system kubeletclient-pod -c alpine -- cat /var/run/secrets/kubernetes.io/serviceaccount/token > http/token.txt

3. Run the tests

        $ cd http
        $ sudo go test

4. Remove the temporary pod

        $ kubectl delete pod kubeletclient-pod -n kube-system
