# Gomrun

It is a simple CLI tool for run multiply go projects. Like a docker-compose.

## Install

```
go install github.com/negasus/gomrun@latest
```

## Run

```
$ gomrun 
$ gomrun --config /path/to/config.yml
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

> todo
