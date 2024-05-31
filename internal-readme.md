## Commands


```
kubebuilder create api --group webapp --version --version v1alpha1 --kind DeclarativeLabels
```


## Design

```
apiVersion: webapp.dmpe.github.io/v1alpha1
kind: DeclarativeLabels
metadata:
    name: declarativelabels-test
spec:
    period: 60 # required, number in seconds
    minNodes: 1 # also required, how many nodes should be labelled
    nodeLabels: # required node will be labelled with these labels
        a: b
        c: d
        e: f

    # future TODO
    excludeNodeNames:
        - 'a'
        - 'b'
    includeNodeNames:
        - '*'
status:
    lastClusterCheck: TimeStamp
```

## Test with KIND cluster

```
make && make generate && make manifests
make docker-build
make docker-push
kind create cluster --config kind.config
kubectl cluster-info --context kind-declarativecluster
make install & make deploy
kubectl apply -f test/k8s/desiredLabel.yaml
kind delete cluster -n mycluster
```

