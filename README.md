# Role watcher
A Discord bot feature that will send a notification if a member doesn't have any of the roles set in ``config.yml``

[Trello card - claimed by @Elian0213](https://trello.com/c/0k7DwbSX)

# Usage
The commands to use this service are as follows:
```
# Create or destroy an role id
.[syntax-command] [add | delete] "role id"
```

# Configuration
Configuration is possible in ``config.yml``

```yaml
"token": "",
"prefix": ".role-watcher",
"notificationchannel": "",
"roles":
    - "role id"
```

# Setup
Download golang if you haven't already at https://golang.org/dl/ after that install the packages 

```
$ go get
$ go build 
```
