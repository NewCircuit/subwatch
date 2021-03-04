# Subscription Watch
This bot removes non-subscribed members, since Discord keeps members if they
have any other roles upon unsubscribing.

## Usage
```
$ cd ./cmd/subwatch
$ go build
$ ./subwatch
$ ... edit config.yml ...
$ ./subwatch
```

## Bot Usage
The bot will then emit to the channel (given in the config.yml) when someone 
is missing the required roles (`roles` in config.yml). 

## Versioning
`x.y.z`
 - X: A major patch: you should update
 - Y: A bug fix: update if this affects you
 - Z: A small change: can be ignored
