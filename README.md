Factory-cli
===========

The cli hub for doing things with factory.

Factory-cli is a GAR client, handling GRPC protocol.

Just like git, naming trick is used : `factory user` calls `factory-user`.

Factory CLI talks to multiple services and handles redirections.

Commands
--------

    factory user inspect

    factory group ls
    factory group inspect GROUP [GROUP...]

    factory project ls
    factory project PROJECT service ls
    factory project PROJECT stop
    factory project PROJECT pause
    factory project PROJECT unpause
    factory project PROJECT service SERVICE logs
    factory project PROJECT service SERVICE scale
    factory project PROJECT container ls

    factory container CONTAINER logs

