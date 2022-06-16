# sidecar-auth-proxy
authorization proxy that intercepts all incoming API calls to backend and validate via a RoundTripper Reverse Proxy
This is an extensible proxy where various handlers running pre-flight or handshake checks can be implement
such as authz, rbac, iam permission checks etc.

