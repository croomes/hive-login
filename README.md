# hive-login

An example project for serving a login page using the kubebuilder webhook server.

## Instructions

Run the server:

```console
$ make run
go fmt ./...
go vet ./...
go run ./main.go
2021-05-25T07:48:08.207Z        INFO    controller-runtime.metrics      metrics server is starting to listen    {"addr": ":8080"}
2021-05-25T07:48:08.208Z        INFO    controller-runtime.webhook      registering webhook     {"path": "/foo/"}
2021-05-25T07:48:08.208Z        INFO    controller-runtime.webhook      registering webhook     {"path": "/baz/"}
2021-05-25T07:48:08.208Z        INFO    setup   starting manager
2021-05-25T07:48:08.209Z        INFO    controller-runtime.manager      starting metrics server {"path": "/metrics"}
2021-05-25T07:48:08.209Z        INFO    controller-runtime.webhook.webhooks     starting webhook server
2021-05-25T07:48:08.209Z        INFO    controller-runtime.certwatcher  Updated current TLS certificate
2021-05-25T07:48:08.209Z        INFO    controller-runtime.webhook      serving webhook server  {"host": "", "port": 9443}
2021-05-25T07:48:08.209Z        INFO    controller-runtime.certwatcher  Starting certificate watcher
...
```

Access the page:

```console
$ curl -k https://localhost:9443/foo/
<!doctype html>
<html>
<head>
    <title>kb webserver site</title>
    <link rel="stylesheet" href="style.css">
    <script src="script.js"></script>
</head>
<body>
<div>
    <h1>kb webserver site</h1>
    <button onclick="popAlert()">Click me</button>
</div>
</body>
</html>
```
