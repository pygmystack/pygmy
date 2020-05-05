---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-05-05
---

# Advanced examples
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

### Traefik 2.x
```yaml
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
      amazeeio-traefik-2:
        Name: amazeeio-traefik-2
```