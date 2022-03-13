---
title: Examples of Pygmy customisation
summary: Examples of the customisation options available in Pygmy
authors:
    - Karl Hepworth
date: 2020-01-28
---
# Introduction

The following are examples of how somebody can utilise pygmy to customise their environment using a `~/.pygmy-yml` file. This file will have a schema which can be imported and the services match the Docker API.

**Standard schema for `~/.pygmy-yml`**
```yaml
# Defaults is a boolean which indicates all default settings should be inherited.
defaults: true

# Resolvers is the Resolv configuration, you can disable this by setting it to [].
resolvers:
  -	Data:   "Contents of the resolvr file/section"
   	File:   "filename.conf"
   	Folder: "/folderpath"
   	Name:   "Human-readable name"

# Services is a hashmap of 
services:

  # The hashmap entry denotes the service name - such as "amazeeio-dnsmasq".
  mycontainer:

    # Config is derrived from the Docker API, intended for container configuration.
    # See https://godoc.org/github.com/docker/docker/api/types/container#Config for the full spec.
    Config:

      # This field is MANDATORY as the value will by default be empty.
      Image: imagename

      # Labels is a key/value pair of labels which will evaluate to string
      # equivelants. For example, Booleans can evaluate to 0/1 or true/false
      # string values depending on the Docker runtime you're using.
      Labels:

        # To enable Pygmy to the configuration, you will need this label.
        # This field is MANDATORY as the value will by default be false.
        pygmy-enable: true

        # You need to give this container a name
        # This field is MANDATORY as the value will by default be empty.
        pygmy-name: mycontainer

        # If you are customising an existing service, you can optionally
        # inherit the defaults if the global defaults are disabled.
        # Setting this value on a non-standard service will do nothing.
        pygmy-defaults: true

        # To display the output when the container starts:
        pygmy-output: true

        # To hide the container from the status messages:
        pygmy-discrete: true

        # To test an endpoint:
        pygmy-url: http://mycontainer.docker.amazee.io

        # To identify the purpose of a container - this is rather specialised so please ignore.
        pygmy-purpose: sshagent

        # To set a weight between 10 and 99 to control the order containers are started:
        pygmy-weight: 50

    # HostConfig is derived from the Docker API, intended for host configuration.
    # See https://godoc.org/github.com/docker/docker/api/types/container#HostConfig for the full spec.
    HostConfig: []

    # NetworkConfig is derived from the Docker API, intended for network configuration.
    # See https://godoc.org/github.com/docker/docker/api/types/network#NetworkingConfig for the full spec.
    NetworkConfig:
      
      # Hashmap value ideally should be the network name - but could be anything.
      # Results may vary, so try what works.
      amazeeio-network:

        # Every network needs a name.
        Name: amazeeio-network

        # An array of Containers.
        Containers:

          # Container name will tell Pygmy to integrate the container of the specified name should be connected to the docker network.
          Name: amazeeio-haproxy

        Labels:

          # Mandatory for network creation/usage via pygmy-
          pygmy-network: true

# networks is a hashmap of the API for a NetworkResource.
# See https://godoc.org/github.com/docker/docker/api/types#NetworkResource for the full spec.
networks: []

# volumes is a hashmap of the API for Volumes
# See https://godoc.org/github.com/docker/docker/api/types#Volume for the full spec.
volumes: []

# keys is all of the SSH key paths which you're utilising.
keys:
  - /home/user1/.ssh/id_rsa
  - /home/user2/.ssh/id_rsa
```

## Applied examples

A suite of examples with a specific purpose are on the way. 