# Got Bindings?

A tiny app that can be useful when testing the claiming and binding of Services on Tanzu Application Platform (TAP).

## Example Usage

```
# clone the repo
$ git clone https://github.com/teddyking/gotbindings.git
$ cd gotbindings
```

```
# deploy to TAP
$ tanzu apps workload create gotbindings \
  --app gotbindings \
  --local-path . \
  --type web \
  --source-image <SOURCE IMAGE>

# wait for app to deploy successfully
$ tanzu app workload get gotbindings
```

```
# confirm that the app don't got bindings
curl -s <WORKLOAD URL> | jq .

{
  "gotBindings": false,
  "bindings": []
}
```

```
# create a test claim for a test Secret
$ kubectl create secret generic service-instance-creds
$ tanzu service claim create test-claim \
  --resource-name service-instance-creds \
  --resource-kind Secret \
  --resource-api-version v1
```

```
# bind the test claim to gotbindings
$ tanzu apps workload update gotbindings \
  --service-ref test-service-instance=services.apps.tanzu.vmware.com/v1alpha1:ResourceClaim:test-claim
```

```
# wait for app to deploy successfully
$ tanzu app workload get gotbindings
```

```
# confirm that the app got bindings
curl -s <WORKLOAD URL> | jq .

{
  "gotBindings": true,
  "bindings": [
    {
      "name": "test-service-instance"
    }
  ]
}
```

