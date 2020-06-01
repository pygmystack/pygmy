---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-05-05
---

# Simple Examples
!!! warning
    All of these examples have been tested in isolation. These are intended to be examples to show what is possible, and any example may not be compatible with each other.

    Using a custom configuration would mean upstream support would be either limited or refused depending on the situation. Usage is considered advanced and an advanced knowledge of Docker would be expected to resolve any issues you may have.

## Modifications to services
### HTTPS on `haproxy`
This example will open up the `haproxy` service to run on `443` so that you can use `pygmy` on HTTPS. It does not currently use a valid certificate so the traffic is considered insecure, but at least this will allow that insecure traffic which can sometimes be absolutely important for local development workflow.
```yaml
services:
  amazeeio-haproxy:
    Config:
      PortBindings:
        443/tcp:
        -
          HostPort: 443
```

### Image replacement
This example shows how an image on one of the default services can be replaced with another. In this example, the `haproxy` is replaced with a custom image which exposes the `haproxy` to port `8080` and not `80`.

Because `haproxy` is being run on port `8080` in this example, we can also note that other services will also be delivered from the new port.
```yaml
services:
  amazeeio-haproxy:
    Config:
      Image: fubarhouse/amazeeio-haproxy-8080
      Labels:
        - pygmy.url: http://docker.amazee.io:8080/stats
    HostConfig:
      PortBindings:
        8080/tcp:
          -
            HostPort: 8080
  
  amazeeio-mailhog:
    Config:
      Labels:
        - pygmy.url: http://mailhog.docker.amazee.io:8080/stats
```

## New services
### Cowsay
This useless example will add a single `cowsay` process in a Docker image to `pygmy up` and will output the stdout to the terminal. This is a fun, useless example of how something independent of `pygmy` can be integrated directly with `pygmy. 
```yaml
services:
  pygmy-cowsay:
    name: pygmy-cowsay
    Config:
      Image: mbentley/cowsay
      Cmd:
        - holy
        - ship
      Labels:
        - pygmy.enable: true
        - pygmy.discrete: true
        - pygmy.name: pygmy-cowsay
        - pygmy.output: true
        - pygmy.weight: 99
    HostConfig:
      AutoRemove: true
```

### PHPMyAdmin
Adding phpmyadmin is reasonably simple, and provides some local development value. You'll be able to manually connect to any database after retrieving the port information for the service you're wanting to connect to. This can be useful to visualise and to manage databases, though ultimately there are better ways of achieving this. 
```yaml
services:
  pygmy-phpmyadmin:
    Config:
      Image: phpmyadmin/phpmyadmin
      Env:
        - "AMAZEEIO=AMAZEEIO"
        - "AMAZEEIO_URL=phpmyadmin.docker.amazee.io"
        - "AMAZEEIO_HTTP_PORT=80"
        - "PMA_ARBITRARY=1"
      Labels:
        - pygmy.enable: true
        - pygmy.name: pygmy-phpmyadmin
        - pygmy.network: amazeeio-network
        - pygmy.weight: 20
        - pygmy.url: http://phpmyadmin.docker.amazee.io
    HostConfig:
      PortBindings:
        80/tcp:
          - HostPort: 8555
```

### Portainer
Portainer provides fantastic management and insight to the local Docker daemon/registry and can be very useful if you're not familiar or eager on the command line interface. Adding it is quite simple and at times can really benefit the developer.
```yaml
services:
  pygmy-portainer:
    Config:
      Image: portainer/portainer
      Env:
        - "AMAZEEIO=AMAZEEIO"
        - "AMAZEEIO_URL=portainer.docker.amazee.io"
        - "AMAZEEIO_HTTP_PORT=9000"
      Labels:
        - pygmy.enable: true
        - pygmy.name: pygmy-portainer
        - pygmy.network: amazeeio-network
        - pygmy.weight: 23
        - pygmy.url: http://portainer.docker.amazee.io
      ExposedPorts:
        9000/tcp: {}
    HostConfig:
      Binds:
        - /var/run/docker.sock:/var/run/docker.sock
        - portainer_data:/data
      PortBindings:
        8000/tcp:
          - HostPort: 8200
        9000/tcp:
          - HostPort: 8100

volumes:
  portainer_data:
    Name: portainer_data
```