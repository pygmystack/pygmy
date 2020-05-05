---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-05-05
---

# Simple Examples

## Modifications to services

### HTTPS on Pygmy
```yaml
services:
  amazeeio-haproxy:
    Config:
      PortBindings:
        443/tcp:
        - HostPort: 443
```

## New services
### Cowsay
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
        - pygmy.weight: 20
        - pygmy.url: http://phpmyadmin.docker.amazee.io
    HostConfig:
      PortBindings:
        80/tcp:
          - HostPort: 8080

networks:
  amazeeio-network:
    Containers:
      pygmy-phpmyadmin: 
        Name: pygmy-phpmyadmin
```

### Portainer
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
        - pygmy: pygmy
        - pygmy.enable: true
        - pygmy.name: pygmy-portainer
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

networks:
  amazeeio-network:
    Containers:
      pygmy-portainer: 
        Name: pygmy-portainer

volumes:
  portainer_data:
    Name: portainer_data
```