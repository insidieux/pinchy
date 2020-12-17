# User guide

## Installing

### Github Release

Visit the [releases page](https://github.com/insidieux/pinchy/releases/latest) to download one of the pre-built binaries for your platform.

### Docker

Use the [Docker image](https://hub.docker.com/repository/docker/insidieux/pinchy)

```shell
docker pull insidieux/pinchy
```

or 

```shell
echo PASSWORD_FILE | docker login docker.pkg.github.com --username USERNAME --password-stdin
docker pull docker.pkg.github.com/insidieux/pinchy
```

### go get

Alternatively, you can use the go get method:

```shell
go get github.com/insidiuex/pinchy/cmd/pinchy
```

Ensure that `$GOPATH/bin` is added to your `$PATH`.

## Usage

### Binary

```shell
pinchy ...
```

### Docker

```shell
docker run insidieux/pinchy
```

or

```shell
docker run docker.pkg.github.com/insidieux/pinchy/pinchy
```

### Available source types

- [file]

[file]: ./source/file.md

### Available registry types

- [consul]

[consul]: ./registry/consul.md
