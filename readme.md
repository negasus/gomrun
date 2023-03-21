# Gomrun

It is a simple CLI tool for build and run multiply go projects. Like a docker-compose.

## FAQ

Q: I can use docker-compose with `build` section. Why do I need `gomrun`?

A: For example, in development you can use `replace` directive in your `go.mod` file.
It can be difficult to build such a project with docker-compose. Also, `gomrun` provide some useful features: `envset`, run binary without build.

Q: I want some feature.

A: It's cool. Issues and PRs are welcome

Q: What's mean `gomrun`?

A: **Go** **M**ultiply **Run**

## Install

```
go install github.com/negasus/gomrun@latest
```

## Run

```
$ gomrun [<service name> <service name> ...]
$ gomrun --config /path/to/config.yml [<service name> <service name> ...]
```

Default config path: `.gomrun.yml`

## Config

```yaml
services:
  s1:
    cmd: echo
    args: ['foo', 'bar']
    
  my_service_1:
    build:
      context: /Users/user/srv1
    envset: ['common']

  my_service_2:
    build:
      context: /Users/user/srv2
      path: ./cmd/main
    args: ['-verbose']
    work_dir: /Users/users/srv2
    envset: ['db', 'common']
    environment:
      GOCACHE: /Users/user/Library/Caches/go-build
      ADDRESS: 127.0.0.1:4403

envset:
  common:
    DEBUG: true
  db:
    POSTGRES_HOST: 127.0.0.1
    POSTGRES_PORT: 5434
    POSTGRES_USERNAME: postgres
    POSTGRES_DATABASE: postgres
    POSTGRES_PASSWORD: secret
    POSTGRES_SSL_MODE: disable

```

### Service

```yaml
serviceName:
  cmd: /path/to/binary
  build:
    context: /Users/user/srv2
    path: ./cmd/main
  args: ['-verbose', '-c', 'config.json']
  work_dir: /Users/users/srv2
  envset: ['db', 'common']
  environment:
    ADDRESS: 127.0.0.1:4403
```

You must use `cmd` or `build` sections. Not both.

`cmd` allows you run some binary

`build` allows to build go project and run it

You can define `args` strings array for run arguments

`envset` allows you to reuse predefined set of environment variables

`environment` allows to define custom env variables for this service
