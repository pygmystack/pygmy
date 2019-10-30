
## Start
To start `pygmy` run following command

    pygmy up

`pygmy` will now start all the required Docker containers and add the ssh key.

If you are on Ubuntu you might need to run pygmy with `pygmy up --no-resolver`

**All done?** Head over to [Drupal Docker Containers](./drupal_site_containers.md) to learn how to work with docker containers.

# Command line usage

```
pygmy help

Commands:
  pygmy addkey [~/.ssh/id_rsa]  # Add additional ssh-key
  pygmy down                    # Stop and destroy all pygmy services
  pygmy help [COMMAND]          # Describe available commands or one specific command
  pygmy restart                 # Stop and restart all pygmy services
  pygmy status                  # Report status of the pygmy services
  pygmy stop                    # Stop all pygmy services
  pygmy up                      # Bring up pygmy services (dnsmasq, haproxy, mailhog, resolv, ssh-agent)
  pygmy update                  # Pulls Docker Images and recreates the Containers
  pygmy version                 # Check current installed version of pygmy
```



## Adding ssh keys

Call the `addkey` command with the **absolute** path to the key you would like to add. In case this they is passphrase protected, it will ask for your passphrase.

    pygmy addkey /Users/amazeeio/.ssh/my_other_key

    Enter passphrase for /Users/amazeeio/.ssh/my_other_key:
    Identity added: /Users/amazeeio/.ssh/my_other_key (/Users/amazeeio/.ssh/my_other_key)

## Checking the status

Run `pygmy status` and `pygmy` will tell you how it feels right now and which ssh-keys it currently has in it's stomach:

    pygmy status

    [*] Dnsmasq: Running as docker container amazeeio-dnsmasq
    [*] Haproxy: Running as docker container amazeeio-haproxy
    [*] Resolv is properly configured
    [*] ssh-agent: Running as docker container amazeeio-ssh-agent, loaded keys:
    4096 SHA256:QWzGNs1r2dfdfX2PHdPi5sdMxdsuddUbPSi7HsrRAwG43sHI /Users/amazeeio/.ssh/my_other_key (RSA)


## `pygmy down` vs `pygmy stop`

`pygmy` behaves like Docker, it's a whale in the end!
During regular development `pygmy stop` is perfectly fine, it will keep the Docker containers still alive, just in stopped state.
If you like to cleanup though, use `pygmy down` to really remove the Docker containers.

## Access HAProxy statistic page and logs  

HAProxy service has statistics web page already enabled. To access the page, just point the browser to [http://docker.amazee.io/stats](http://docker.amazee.io/stats).  

To watch at haproxy container logs, use the `docker logs amazeeio-haproxy` command with standard `docker logs` options like `-f` to follow.
