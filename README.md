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
