# User guide

## Installing

### Github Release

Visit the [releases page](https://github.com/insidieux/pinchy/releases/latest) to download one of the pre-built binaries
for your platform.

### Docker

Use the [Docker image](https://hub.docker.com/repository/docker/insidieux/pinchy)

```shell
docker pull insidieux/pinchy
```

### go get

Alternatively, you can use the go get method:

```shell
go get github.com/insidieux/pinchy/cmd/pinchy
```

Ensure that `$GOPATH/bin` is added to your `$PATH`.

## Usage

### Binary

```shell
pinchy %source% %registry% %mode% [flags] 
```

### Docker

```shell
docker run insidieux/pinchy:latest %source% %registry% %mode% [flags]
```

### Docker-compose run

Example docker-compose file be found in [deployment](./../deployments/docker-compose/pinchy/docker-compose.yml)
directory

## Modes

`once` mode run sync process only single time

`watch` mode run sync process repeatedly with constant `schedule.interval`

## Command common flags

```
--logger.level string     Log level (default "info")
--manager.exit-on-error   Stop manager process on first error and by pass it to command line
```

### Watch mode

```
--scheduler.interval duration   Interval between manager runs (1s, 1m, 5m, 1h and others) (default 1m0s)
```

### Source and Registry flags

Flags for chosen `source` and `registry` are described in a related documentation for sources and registry types.

## Available source types

- [file]

[file]: ./source/file.md

## Available registry types

- [consul]

[consul]: ./registry/consul.md

## Examples

### Once

```
pinchy \
    file \
    consul \
    once \
    --source.path /etc/pinchy/services.yml \
    --registry.address http://127.0.0.1:8500
```

### Watch

```
pinchy \
    file \
    consul \
    watch \
    --source.path /etc/pinchy/services.yml \
    --registry.address http://127.0.0.1:8500 \
    --scheduler.interval 5s
```
