### Query the number of pods created by StatefulSets:

```
sum(kube_pod_owner{namespace="monitoring", owner_kind="StatefulSet"}) by (owner_name)
```

### Alternative Approach (Using kube_pod_info)

```
sum(kube_pod_info{namespace="monitoring", created_by_kind="StatefulSet"}) by (created_by_name)
```
