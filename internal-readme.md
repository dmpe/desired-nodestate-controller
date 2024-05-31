## Commands


```
kubebuilder create api --group webapp --version --version v1alpha1 --kind DeclarativeLabels
```


## Design

```
kind: DeclarativeLabels
apiVersion: v1alpha1
metadata:
    name: declarativeLabels-test
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