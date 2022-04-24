# hello-envoy

This is an example of Envoy TCP Proxy from localhost:10000 to httpbin.org:80.
See also https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/listeners/tcp_proxy.

It is based on https://gist.github.com/int128/5b02bd1f9b7882aa6b48727b838e658e.

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

## Dynamic configuration

### Set up Envoy

To run a proxy:

```console
% docker run --rm envoyproxy/envoy:v1.22-latest --version

envoy  version: dcd329a2e95b54f754b17aceca3f72724294b502/1.22.0/Clean/RELEASE/BoringSSL

% docker run --rm -p 10000:10000 -v $PWD/envoy_dynamic:/etc/envoy envoyproxy/envoy:v1.22-latest -c /etc/envoy/envoy.yaml
...
[2022-04-24 07:39:54.948][1][info][config] [source/server/listener_manager_impl.cc:789] all dependencies initialized. starting workers
```

Make sure httpbin returns a response.

```console
% curl -v http://localhost:10000/get
...
  "url": "http://localhost/get"
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
