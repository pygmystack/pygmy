# Update pygmy

As `pygmy` is an active project, you should also take care of updating pygmy-

Use the [same instructions](./installation.md) to update Pygmy as to install it. 

## I see errors or unexpected behaviour after the upgrade

If you see anything unexpected after upgrading, the recommended advice is to clean up the environment _and_ remove the docker network.

Any applications which use the network `amazeeio-network` such as a docker-compose Drupal project - should not be running. You can alternatively run `docker network rm amazeeio-network --force`.

```console
$ pygmy clean
$ docker network rm amazeeio-network
```

## Update Docker Containers with `pygmy`

`pygmy` can update shared docker containers for you:

    pygmy update && pygmy restart
