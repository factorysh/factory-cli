# Factory-cli

[![Build Status](https://drone.bearstech.com/api/badges/factorysh/factory-cli/status.svg)](https://drone.bearstech.com/factorysh/factory-cli)

The cli hub for doing things with factory.

Factory CLI talks to multiple services and handles redirections.

## Commands

```
./bin/factory

 _
| |             __            _
| | /| /| /|   / _| __ _  ___| |_ ___  _ __ _   _
| |/ |/ |/ |  | |_ / _' |/ __| __/ _ \| '__| | | |
|          |  |  _| (_| | (__| || (_) | |  | |_| |
+----------+  |_|  \__,_|\___|\__\___/|_|   \__, |
                                             |___/

Full documentation:
  https://github.com/factorysh/factory-cli/blob/master/README.fr.adoc

Usage:
  factory [command]

Available Commands:
  container   Do something on a container
  help        Help about any command
  infos       Show project's infos
  journal     Show journal
  volume      Do something on a volume

Flags:
      --config string    Config file (default is $HOME/.factory-cli.yaml)
  -g, --gitlab string    Gitlab server url (default "github.com")
  -h, --help             help for factory
  -p, --project string   Gitlab project path (default "factorysh/factory-cli")
  -t, --token string     Gitlab token
  -v, --verbose          Verbose output

Use "factory [command] --help" for more information about a command.
```

## Licence

GPLv3 ©2019 Bearstech
