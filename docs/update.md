# Updates

## Updating pygmy
As `pygmy-go` is an active project, you should also take care of updating pygmy.

Use the [same instructions](./installation.md) to update Pygmy as to install it. 

## Updating the images

`pygmy-go` can update shared docker containers for you:

    pygmy-go update && pygmy-go restart

## Dealing with unexpected behaviour

If you see anything unexpected after upgrading, the recommended advice is to clean up the environment _and_ remove the docker network.

Any applications which use the network `amazeeio-network` such as a docker-compose Drupal project - should not be running.

You can alternatively run `docker network rm amazeeio-network --force` to force the docker network to be removed.

```console
$ pygmy-go clean
$ docker network rm amazeeio-network
```