---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-05-05
---

# Advanced examples
!!! warning
    All of these examples have been tested in isolation. These are intended to be examples to show what is possible, and any example may not be compatible with each other.

    Using a custom configuration would mean upstream support would be either limited or refused depending on the situation. Usage is considered advanced and an advanced knowledge of Docker would be expected to resolve any issues you may have.

## Modifications to services
### Custom mailhub port for mailhog
Switching the mailhub port will also mean the container sending the mail will need to be configured to use the new port as well. This will sllow for multiple instances of pygmy as required outside isolation that docker-compsoe provides.
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
Replacing the SSH agent for any reason is an involved process, because the base container config will need a new configuration item, and the same with shower/adder containers. This change would also break Windows compatibility because the images in use by default also support Windows when configured correctly. This is not a task that is justified other than a proof of concept to prove that it can be done.
```yaml
services:
  amazeeio-ssh-agent:
    Config:
      Labels:
        - pygmy.enable: false

  amazeeio-ssh-agent-add-key:
    Config:
      Labels:
        - pygmy.enable: false

  amazeeio-ssh-agent-show-keys:
    Config:
      Labels:
        - pygmy.enable: false

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
Pygmy-go was created from the desire to integrate Traefik 1.x with the Ruby version of `pygmy` It's a replacement to `haproxy` which can provide a `sudo`less and resolver-free experience. This one doesn't require additional configuration to get going, unlike the configuration for Traefik 2.x.
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
Traefik 2.0 is made for Kubernernetes, so with it comes a lot of difference in configuration for pygmy. This one will require additional configuration. This one is tailored for a single service in a docker-compose project where the labels `traefik.enable` is `true`, and `traefik.port` is the port Traefik will be looking on.
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
defaults: false
services:

  pygmy-ssh-agent:
    Config:
      Cmd:
        - ssh-agent
      Env:
        - "SSH_AUTH_SOCK=/.ssh-agent/socket"
      Image: nardeas/ssh-agent
      Labels:
        - pygmy.name: pygmy-ssh-agent
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

  pygmy-ssh-agent-add-key:
    Config:
      Image: nardeas/ssh-agent
      Labels:
        - pygmy.name: pygmy-ssh-agent-add-key
        - pygmy.enable: true
        - pygmy.discrete: true
        - pygmy.output:  true
        - pygmy.purpose: addkeys
        - pygmy.weight: 31
    HostConfig:
      AutoRemove: true
      IpcMode: private
      VolumesFrom:
        - pygmy-ssh-agent

  pygmy-ssh-agent-show-keys:
    Config:
      Cmd:
        - ssh-add
        - -l
      Image: nardeas/ssh-agent
      Labels:
        - pygmy.name: pygmy-ssh-agent-show-keys
        - pygmy.enable: true
        - pygmy.discrete: true
        - pygmy.output: true
        - pygmy.purpose: showkeys
        - pygmy.weight: 32
    HostConfig:
      AutoRemove: true
      VolumesFrom:
        - pygmy-ssh-agent

  amazeeio-mailhog:
    Config:
      Labels:
        - pygmy.enable: true
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.mailhog.rule: Host(`mailhog.docker.amazee.io`)

  pygmy-phpmyadmin:
    Config:
      Image: phpmyadmin/phpmyadmin
      Env:
        - "PMA_ARBITRARY=1"
      Labels:
        - pygmy.enable: true
        - pygmy.name: pygmy-phpmyadmin
        - pygmy.weight: 20
        - pygmy.url: http://phpmyadmin.docker.amazee.io
        - traefik.enable: true
        - traefik.port: 80
        - traefik.http.routers.phpmyadmin.rule: Host(`phpmyadmin.docker.amazee.io`)
    HostConfig:
      PortBindings:
        80/tcp:
          - HostPort: 8770

  pygmy-portainer:
    Config:
      Image: portainer/portainer
      Labels:
        - pygmy.enable: true
        - pygmy.name: pygmy-portainer
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

  pygmy-traefik:
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
        - pygmy.name: pygmy-traefik
        - pygmy.enable: true
        - pygmy.url: http://traefik.docker.amazee.io
        - traefik.docker.network: pygmy-network
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
  pygmy-network:
    Name: pygmy-network
    Containers:
      amazeeio-haproxy: {}
      unofficial-traefik-2:
        Name: pygmy-traefik
      amazeeio-mailhog:
        Name: amazeeio-mailhog
      unofficial-portainer:
        Name: pygmy-portainer
      unofficial-phpmyadmin:
        Name: pygmy-phpmyadmin
    Labels:
      - pygmy.network: true

volumes:
  portainer_data:
    Name: portainer_data
```