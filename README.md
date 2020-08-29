## Admission Webhook Demo
#### MutatingWebhook
The webhook is for injecting a sidecar container when a pod is creating.

#### ValidatingWebhook
The webhook is for checking something.

## Attention
If you deal with the `subResource`, you should make sure that your kubernetes version is greater than 1.15.0 due to the [bug](https://github.com/kubernetes/kubernetes/issues/67221), and the bug has be [fixed](https://github.com/kubernetes/kubernetes/pull/76849) after 1.15.0

If your kubernetes version is less than 1.15.0, it will be failed with below error when execute `kubectl scale`:
```
Error from server (InternalError): Internal error occurred: converting (apps.Deployment).Replicas to (v1beta1.Scale).Replicas: Selector not present in src
```

Besides, `HPA` will also not work, because `HPA` send a scale object to the apiserver.
