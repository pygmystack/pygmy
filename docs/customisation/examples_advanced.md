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
