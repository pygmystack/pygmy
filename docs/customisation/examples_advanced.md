---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-05-05
---

# Advanced examples
## Modifications to services
```yaml
services:
  amazeeio-mailhog:
    HostConfig:
      PortBindings:
        1025/tcp:
        - HostPort: 2025
```

## Replacement services
### SSH Agent
```yaml
services:
  unofficial-ssh-agent:
    Config:
      Cmd:
        - ssh-agent
      Env:
        - "SSH_AUTH_SOCK=/.ssh-agent/socket"
      Image: nardeas/ssh-agent
      Labels:
        - pygmy.name: unofficial-ssh-agent
        - pygmy.enable: true
        - pygmy.output: false
        - pygmy.purpose: sshagent
        - pygmy.weight:  30
    HostConfig:
      AutoRemove: false
      IpcMode: private
      RestartPolicy:
        Name: always
        MaximumRetryCount: 0

  unofficial-ssh-agent-add-key:
    Config:
      Image: nardeas/ssh-agent
      Labels:
        - pygmy.name: unofficial-ssh-agent-add-key
        - pygmy.enable: true
        - pygmy.discrete: true
        - pygmy.output:  true
        - pygmy.purpose: addkeys
        - pygmy.weight: 31
    HostConfig:
      AutoRemove: true
      IpcMode: private
      VolumesFrom:
        - unofficial-ssh-agent

  unofficial-ssh-agent-show-keys:
    Config:
      Cmd:
        - ssh-add
        - -l
      Image: nardeas/ssh-agent
      Labels:
        - pygmy.name: unofficial-ssh-agent-show-keys
        - pygmy.enable: true
        - pygmy.discrete: true
        - pygmy.output: true
        - pygmy.purpose: showkeys
        - pygmy.weight: 32
    HostConfig:
      AutoRemove: true
      VolumesFrom:
        - unofficial-ssh-agent
```

### Traefik 1.x
```yaml
services:

  amazeeio-haproxy:
    Config:
      Labels:
        - pygmy.enable: false

  pygmy-traefik-1:
    Config:
      Image: library/traefik:v1.7.19
      Cmd:
        - --api
        - --docker
        - --docker.network=amazeeio-network
        - --docker.domain=docker.amazee.io
        - --docker.exposedbydefault=false
      ExposedPorts:
        80/tcp:
          HostPort: 80
        443/tcp:
          HostPort: 443
        8080/tcp:
          HostPort: 8080
      Labels:
        - pygmy: pygmy
        - pygmy.name: pygmy-traefik-1
        - pygmy.enable: true
        - pygmy.url: http://traefik.docker.amazee.io
        - traefik.enable: true
        - traefik.port: 8080
        - traefik.protocol: http
        - traefik.docker.domain: docker.amazee.io
        - traefik.docker.exposedByDefault: false
        - traefik.docker.network: amazeeio-network
        - traefik.frontend.rule: Host:traefik.docker.amazee.io
    HostConfig:
      Binds:
        - /var/run/docker.sock:/var/run/docker.sock
      PortBindings:
        443/tcp:
          - HostPort: 443
        80/tcp:
          - HostPort: 80
        8080/tcp:
          - HostPort: 8080
      RestartPolicy:
        Name: always
        MaximumRetryCount: 0
    NetworkConfig:
      Ports:
        80/tcp:
          - HostPort: 80
        8080/tcp:
          - HostPort: 8080

networks:
  amazeeio-network:
    Containers:
      pygmy-traefik-1:
        Name: pygmy-traefik-1

resolvers: []
```

### Traefik 2.x
```yaml
  amazeeio-haproxy:
    Config:
      Labels:
        - pygmy.enable: false

  pygmy-traefik-2:
    Config:
      Image: library/traefik:v2.1.3
      Cmd:
        - --api
        - --api.insecure=true
        - --providers.docker
        - --providers.docker.exposedbydefault=false
        - --providers.docker.defaultrule=Host(`{{ index .Labels "com.docker.compose.project" }}.docker.amazee.io`)
        - --entrypoints.web.address=:80
        - --entrypoints.websecure.address=:443
      ExposedPorts:
        80/tcp:
          HostPort: 80
        443/tcp:
          HostPort: 443
        8080/tcp:
          HostPort: 3080
      Labels:
        - pygmy.enable: true
        - pygmy.name: pygmy-traefik-2
        - pygmy.url: http://traefik.docker.amazee.io
        - traefik.docker.network: amazeeio-network
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.traefik.rule: Host(`traefik.docker.amazee.io`)
        - traefik.http.routers.traefik.service: api@internal
        - traefik.providers.docker.defaultport: 8080
    HostConfig:
      Binds:
        - /var/run/docker.sock:/var/run/docker.sock
      PortBindings:
        443/tcp:
          - HostPort: 443
        80/tcp:
          - HostPort: 80
        8080/tcp:
          - HostPort: 8080
      RestartPolicy:
        Name: always
        MaximumRetryCount: 0
    NetworkConfig:
      Ports:
        80/tcp:
          - HostPort: 80
        8080/tcp:
          - HostPort: 8080

networks:
  amazeeio-network:
    Containers:
      pygmy-traefik-2:
        Name: pygmy-traefik-2

resolvers: []
```

## Completely replacing everything
```yaml
services:
  amazeeio-haproxy:
    Config:
      Labels:
        - pygmy.enable: false

  amazeeio-dnsmasq:
    Config:
      Labels:
        - pygmy.enable: false

  amazeeio-mailhog:
    Config:
      Labels:
        - pygmy.enable: true
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.mailhog.rule: Host(`mailhog.docker.amazee.io`)

  unofficial-phpmyadmin:
    Config:
      Image: phpmyadmin/phpmyadmin
      Env:
        - "PMA_ARBITRARY=1"
      Labels:
        - pygmy.enable: true
        - pygmy.name: unofficial-phpmyadmin
        - pygmy.weight: 20
        - pygmy.url: http://phpmyadmin.docker.amazee.io
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.phpmyadmin.rule: Host(`phpmyadmin.docker.amazee.io`)
    HostConfig:
      PortBindings:
        80/tcp:
          - HostPort: 8770

  unofficial-portainer:
    Config:
      Image: portainer/portainer
      Labels:
        - pygmy.enable: true
        - pygmy.name: unofficial-portainer
        - pygmy.url: http://portainer.docker.amazee.io
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.portainer.rule: Host(`portainer.docker.amazee.io`)
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

  unofficial-traefik-2:
    Config:
      Image: library/traefik:v2.1.3
      Cmd:
        - --api
        - --api.insecure=true
        - --providers.docker
        - --providers.docker.exposedbydefault=false
        - --providers.docker.defaultrule=Host(`{{ index .Labels "com.docker.compose.project" }}.docker.amazee.io`)
        - --entrypoints.web.address=:80
        - --entrypoints.websecure.address=:443
      ExposedPorts:
        80/tcp:
          HostPort: 80
        443/tcp:
          HostPort: 443
        8080/tcp:
          HostPort: 3080
      Labels:
        - pygmy.name: unofficial-traefik-2
        - pygmy.enable: true
        - pygmy.url: http://traefik.docker.amazee.io
        - traefik.docker.network: amazeeio-network
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.traefik.rule: Host(`traefik.docker.amazee.io`)
        - traefik.http.routers.traefik.tls: true
        - traefik.http.routers.traefik.service: api@internal
        - traefik.providers.docker.defaultport: 8080
    HostConfig:
      Binds:
        - /var/run/docker.sock:/var/run/docker.sock
      PortBindings:
        443/tcp:
          - HostPort: 443
        80/tcp:
          - HostPort: 80
        8080/tcp:
          - HostPort: 8080
      RestartPolicy:
        Name: always
        MaximumRetryCount: 0
    NetworkConfig:
      Ports:
        80/tcp:
          - HostPort: 80
        8080/tcp:
          - HostPort: 8080

networks:
  amazeeio-network:
    Name: amazeeio-network
    Containers:
      amazeeio-haproxy: {}
      unofficial-traefik-2:
        Name: unofficial-traefik-2
      amazeeio-mailhog:
        Name: amazeeio-mailhog
      unofficial-portainer:
        Name: unofficial-portainer
      unofficial-phpmyadmin:
        Name: unofficial-phpmyadmin
    Labels:
      - pygmy.network: true

resolvers: []

volumes:
  portainer_data:
    Name: portainer_data
```