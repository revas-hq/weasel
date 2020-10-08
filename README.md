# Cloud Object Static Webserver

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

> Inspired by <https://github.com/google/weasel> using go-kit architecture.

A simple frontend web server that serves content from a storage cloud solution.

## Features

- HTTPS on custom domain with forced TLS and STS responses
- [WIP] Robust redirect from naked custom domain
- Deployment flow using existing cloud storage solutions
- Efficient cache system with various providers
- [WIP] SPA website mode (404 redirect with router history mode)
- HTTP/2 push with object metadata (infrastructure should support push like Google Frontend on Google App Engine)
- [WIP] Etag If-None-Match optimization middleware

### Cloud storage solutions

- Google Cloud Storage

### Cache solutions

https://developers.google.com/web/fundamentals/performance/optimizing-content-efficiency/http-caching?hl=it

- inmem with bigcache
- [WIP] memcache server

## Misconfiguration error

http.StatusTeapot is used in case of misconfiguration error.
