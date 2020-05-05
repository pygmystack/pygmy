---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-05-05
---

# Simple Examples

## PHPMyAdmin
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