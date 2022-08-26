# Got Bindings?

A tiny app that can be useful when testing the claiming and binding of Services on Tanzu Application Platform (TAP).

## Example Usage

```
# on a TAP cluster ...
$ kubectl apply -f https://raw.githubusercontent.com/teddyking/gotbindings/main/config/tap.yml

# wait for Workload to become Ready, then visit/curl the URL endpoint
# add or remove serviceClaims to/from the Workload as desired
```

