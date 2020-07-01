# Subcription Watch
This is Floor Gang bot, it watches for non-paying members

## Usage
```
$ go build
$ ./subwatch
$ ... edit config.yml ...
$ ./subwatch
```

## Bot Usage
The bot will then emit to the channel (given in the config.yml) when someone 
is missing the required roles (`roles` in config.yml). 
