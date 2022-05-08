# hello-envoy [![envoy](https://github.com/int128/hello-envoy/actions/workflows/envoy.yaml/badge.svg)](https://github.com/int128/hello-envoy/actions/workflows/envoy.yaml)

This is an example of Envoy TCP Proxy with the dynamic configuration.
It is based on https://gist.github.com/int128/5b02bd1f9b7882aa6b48727b838e658e.

## :book: Table of Contents

- Static configuration (Docker)
- Dynamic configuration (Docker)
- Dynamic configuration (Kubernetes)

## Static configuration

To run a proxy:

```console
% docker run --rm envoyproxy/envoy:v1.22-latest --version

envoy  version: dcd329a2e95b54f754b17aceca3f72724294b502/1.22.0/Clean/RELEASE/BoringSSL

% docker run --rm -p 10000:10000 -v $PWD/envoy_static.yaml:/envoy_static.yaml:ro envoyproxy/envoy:v1.22-latest -c /envoy_static.yaml
...
[2022-04-24 07:13:37.518][1][info][config] [source/server/listener_manager_impl.cc:789] all dependencies initialized. starting workers
```

To connect to the remote host via the proxy:

```console
% curl -v http://localhost:10000/get
*   Trying 127.0.0.1:10000...
* Connected to localhost (127.0.0.1) port 10000 (#0)
> GET /get HTTP/1.1
> Host: localhost:10000
> User-Agent: curl/7.79.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sun, 24 Apr 2022 07:15:25 GMT
< Content-Type: application/json
< Content-Length: 250
< Connection: keep-alive
< Server: gunicorn/19.9.0
< Access-Control-Allow-Origin: *
< Access-Control-Allow-Credentials: true
<
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "localhost",
    "User-Agent": "curl/7.79.1",
    "X-Amzn-Trace-Id": "Root=1-6264f90d-79e432f97032afa8435446e0"
  },
  "origin": "192.168.0.0",
  "url": "http://localhost/get"
}
* Connection #0 to host localhost left intact
```

Check the [access log](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/observability/access_logging) in stdout.

```console
[2022-05-08T09:37:42.496Z] "- - -" 0 - 82 480 529 - "-" "-" "-" "-" "72.44.59.17:80"
```

## Dynamic configuration

### Set up Envoy

To run a proxy:

```console
% docker run --rm envoyproxy/envoy:v1.22-latest --version

envoy  version: dcd329a2e95b54f754b17aceca3f72724294b502/1.22.0/Clean/RELEASE/BoringSSL

% docker run --rm -p 10000:10000 -p 19901:19901 -v $PWD/envoy_dynamic:/etc/envoy envoyproxy/envoy:v1.22-latest -c /etc/envoy/bootstrap.yaml
...
[2022-04-24 07:39:54.948][1][info][config] [source/server/listener_manager_impl.cc:789] all dependencies initialized. starting workers
```

Make sure httpbin returns a response.

```console
% curl -v http://localhost:10000/get
...
  "url": "http://localhost/get"
```

Make sure the ready endpoint returns `LIVE`.

```console
% curl -v http://localhost:19901/ready
...
< HTTP/1.1 200 OK
< content-type: text/plain; charset=UTF-8
< cache-control: no-cache, max-age=0
< x-content-type-options: nosniff
< date: Sun, 08 May 2022 11:55:47 GMT
< server: envoy
< x-envoy-upstream-service-time: 2
< transfer-encoding: chunked
<
LIVE
```

### Change the config without restart

Change the xDS config in the container.

```console
% docker ps -a
dc1323f80ea2

% docker exec dc1323f80ea2 sed -i -e 's/httpbin.org/google.com/' /etc/envoy/cds.yaml
```

Check the log of Envoy container.

```console
[2022-04-24 07:58:50.477][1][info][upstream] [source/common/upstream/cds_api_helper.cc:30] cds: add 1 cluster(s), remove 0 cluster(s)
[2022-04-24 07:58:50.479][1][info][upstream] [source/common/upstream/cds_api_helper.cc:67] cds: added/updated 1 cluster(s), skipped 0 unmodified cluster(s)
```

Make sure Google returns a response.

```console
% curl -v http://localhost:10000/get
...
  <a href=//www.google.com/><span id=logo aria-label=Google></span></a>
```

## Dynamic configuration (Kubernetes)

### Set up

Create a cluster.

```console
% kind create cluster
 ✓ Ensuring node image (kindest/node:v1.23.4) 🖼
 ✓ Preparing nodes 📦
...

% kubectl version
Client Version: version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.6", GitCommit:"ad3338546da947756e8a88aa6822e9c11e7eac22", GitTreeState:"clean", BuildDate:"2022-04-14T08:41:58Z", GoVersion:"go1.18.1", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.4", GitCommit:"e6c093d87ea4cbb530a7b2ae91e54c0842d8308a", GitTreeState:"clean", BuildDate:"2022-03-06T21:32:53Z", GoVersion:"go1.17.7", Compiler:"gc", Platform:"linux/amd64"}
```

Deploy the Envoy.

```console
% kubectl create cm envoy --from-file=envoy_dynamic
configmap/envoy created

% kubectl apply -f envoy_deployment.yaml
deployment.apps/envoy created

% kubectl get deployments envoy
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
envoy   1/1     1            1           54s

% kubectl logs -l app=envoy
...
[2022-04-24 09:06:38.547][1][info][config] [source/server/listener_manager_impl.cc:789] all dependencies initialized. starting workers
```

Make sure httpbin returns a response.

```console
% kubectl port-forward "$(kubectl get pods -l app=envoy -oname)" 10000:10000
Forwarding from 127.0.0.1:10000 -> 10000
Forwarding from [::1]:10000 -> 10000

% curl -v http://localhost:10000/get
...
  "url": "http://localhost/get"
```

### Change the config without restart

Change the xDS config in the ConfigMap.

```console
% kubectl edit cm envoy
configmap/envoy edited
```

Check the log of Envoy container.

```console
% kubectl logs -l app=envoy
...
[2022-04-24 09:16:19.350][1][info][upstream] [source/common/upstream/cds_api_helper.cc:30] cds: add 1 cluster(s), remove 0 cluster(s)
[2022-04-24 09:16:19.351][1][info][upstream] [source/common/upstream/cds_api_helper.cc:67] cds: added/updated 1 cluster(s), skipped 0 unmodified cluster(s)
```

Make sure Google returns a response.

```console
% curl -v http://localhost:10000/get
...
  <a href=//www.google.com/><span id=logo aria-label=Google></span></a>
```

### Clean up

Delete the cluster.

```console
% kind delete cluster
Deleting cluster "kind" ...
```

## :link: References

- https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/listeners/tcp_proxy
- https://www.envoyproxy.io/docs/envoy/latest/start/quick-start/run-envoy
- https://www.envoyproxy.io/docs/envoy/latest/start/quick-start/configuration-dynamic-filesystem

If you are looking for a HTTP(S) proxy, see https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/http instead.
