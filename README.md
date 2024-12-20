# Iracluster for Go

### Dependencies

#### VCPKG
vcpkg.json
```json
{
  "name": "iracluster-dependencies",
  "version": "0.1.0",
  "dependencies": [
    "boost-asio",
    "boost-beast",
    "boost-convert",
    "boost-coroutine",
    "boost-crc",
    "boost-date-time",
    "boost-filesystem",
    "boost-stacktrace",
    "boost-system",
    "boost-thread",
    "boost-timer",
    "boost-uuid",
    "brotli",
    "cnats",
    "curl",
    "openssl",
    "rapidjson",
    "spdlog",
    "sqlite3",
  ],
  "overrides": [
    {"name": "fmt", "version": "9.1.0#1"},
    {"name": "openssl", "version": "3.0.8#2"}
  ]
}
```
vcpkg-configuration.json
```json
{
  "default-registry": {
    "kind": "git",
    "baseline": "d033613d9021107e4a7b52c5fac1f87ae8a6fcc6",
    "repository": "https://github.com/microsoft/vcpkg"
  },
  "registries": [
    {
      "kind": "artifact",
      "location": "https://github.com/microsoft/vcpkg-ce-catalog/archive/refs/heads/main.zip",
      "name": "microsoft"
    }
  ]
}
```

### Install C++ IraCluster library

```shell
cp libiracommon.a libiracluster.a /usr/local/lib
```

### Install go-iracluster
```shell
go get github.com/epicira/go-iracluster
```

### Build applications with go-iracluster
Debian
```shell
export CGO_LDFLAGS="-Lvcpkg_installed/x64-linux/lib -lnats_static -lbrotlicommon -lbrotlidec -lbrotlienc -lsqlite3 -lcurl -lz -lspdlog -lboost_system -lboost_coroutine -lboost_stacktrace_backtrace -lboost_thread -lboost_timer -lboost_date_time -lboost_filesystem -lfmt -lssl -lcrypto -lsodium -L/usr/local/lib -liracluster -L/usr/lib/x86_64-linux-gnu/ -lm -lstdc++"
```

```shell
go build -o ./bin/example .
```
OR

The below command is just for reference. Not recommended as glibc cannot be fully made static if certain features are used.
```shell
go build --ldflags '-s -w -linkmode external -extldflags=-static'  -o ./bin/example .
```
