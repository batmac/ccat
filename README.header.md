# ccat
[![Go](https://github.com/batmac/ccat/actions/workflows/go.yml/badge.svg)](https://github.com/batmac/ccat/actions/workflows/go.yml)
![GitHub](https://img.shields.io/github/license/batmac/ccat)
[![Go Report Card](https://goreportcard.com/badge/github.com/batmac/ccat)](https://goreportcard.com/report/github.com/batmac/ccat)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbatmac%2Fccat.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbatmac%2Fccat?ref=badge_shield)
[![codecov](https://codecov.io/gh/batmac/ccat/branch/main/graph/badge.svg?token=PCD6DM6S75)](https://codecov.io/gh/batmac/ccat)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/batmac/ccat/main.svg)](https://results.pre-commit.ci/latest/github/batmac/ccat/main)

cat on steroids.
Leveraging great go modules to ease my CLI life.

## install
- get the latest release from https://github.com/batmac/ccat/releases/latest
- or `go install github.com/batmac/ccat/cmd/ccat@latest`

then, you can update to the latest github release with `ccat --selfupdate`

## build
you need a recent go, available build tags:
- `libcurl`: build with the libcurl opener.
- `fileonly`: build with the local file opener only.
- `nomd`: build without the markdown interpreter (glamour).
- `nohl`: build without the syntax-highlighter.
- `crappy`: build with some crappy (but useful) openers/mutators (needs a recent go version).

for instance:
`go build --tags libcurl,crappy .`

## examples
```
cat -X "kubectl get all" -i -w -t ready,Running --bg
NAME           READY   STATUS    RESTARTS   AGE
pod/busybox3   1/1     Running   1          17m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   21d
```
( "READY and "Running" have different background colors )

```
$ echo Hello: | cat -m y2j,j2y,base64,hexdump
00000000  53 47 56 73 62 47 38 36  49 47 35 31 62 47 77 4b  |SGVsbG86IG51bGwK|
```
```
$ cat -m zstd LICENSE | cat -m base64
KLUv/WQwA+UUAGKyjCQQi1gA+Mu1uPH4trE1CdEhIptwP1oghhOOOP5rTRkCAKACAAdroEXGmA2Uz10inntYY97kIsI70zHae3CBb+DnAi4lAWNzxHPSRXalSceYEc4rc2FxPri8gX/QzYX9sT/32JzZHMbds2JzJmwIlmDMlYRjz/LGa8+BueftcBkb80Dn00G42FDA5RFeRujMBVvNmX8e5sjP1xLa4ftzH5ivJRriudg89+A789ibhIq4UDyyZ7G4GZni7BE2l7K96ZqlNKZ73IHNx5wltsQbeNlcc+bhPNOmS7jnE77PxeGI55jFszx3hWnuUSTcivBStqhQaFeY5xDeeOwNfGM+FiNjfC6gMx3jCoB0lwzgTGNcAoyKS/cWGWM2UL45c60B6ko+0c6WZuyo8Cl6249qdXyaqIRxSxnqfO0rw3XQV2mudYDt7dyRvdAhthvuNN3Hzc9QXbj4N6oSbktVj4U7cqWG33ycO0JVyW/+cW3JhdY0VXzqEkSWXpkcFb5Sc+VoKuOOrki4J0clqZErk0dHN3JFqbj6zQu/SrkjPUFsR9UE2b5ib2eoqYwr9VxdKXfDX6V822uiNRIAkACvUg44rlM/bm9nmm2NVNCNHtVMG7oY00eWvlK0IT8uPm4p/3HXRHKh00Ns26bZcW8vpqrkE0t/85UKoagE/XFfK/W4EtVsKQX5xFGpKh337W2MO0JVxVr2ZkuX14w7S+jhBgYFBheAAp4lSFMTyV3WqthWqJ9PlwgmAGYTiQo1Zi5cSUcnOBoEX3ev5YUMcVnvUrgx5HaCOLWhPcddUjOGZmX2yZFNAyB+ELmdEtsBPE17UsXU4fTV1jdNHpPl8KiAkC34tsZqWhxaa3RlIgn25jJkslDF+VBC8cW5mikplMUF0g==
```

## docker image
multi-arch container images with tags `libcurl,crappy` are automatically built by [Github Actions](https://github.com/batmac/ccat/actions/workflows/docker-images.yml) and made available on [Dockerhub](https://hub.docker.com/r/batmac/ccat) and [Github Packages](https://github.com/batmac/ccat/pkgs/container/ccat). <br/>
for instance:
```
docker run --rm ghcr.io/batmac/ccat:latest -h
```
or
```
kubectl run -i --tty ccat --image=batmac/ccat:latest -- /bin/sh
```

## help
