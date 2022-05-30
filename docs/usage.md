
## Start
To start `pygmy` run following command

    pygmy up

`pygmy` will now start all the required Docker containers and add the ssh key.

If you are on Ubuntu you might need to run pygmy with `pygmy up --no-resolver`

**All done?** Head over to [Drupal Docker Containers](./drupal_site_containers.md) to learn how to work with docker containers.

# Command line usage

```
Amazeeio's local development tool,
        
Runs DNSMasq, HAProxy, MailHog and an SSH Agent in local containers for local development.

Usage:
  pygmy [command]

Available Commands:
  addkey      Add/re-add an SSH key to the agent
  clean       Stop and remove all pygmy services regardless of state
  down        Stop and remove all pygmy services
  export      Export validated configuration to a given path
  help        Help about any command
  restart     Restart all pygmy containers.
  status      Report status of the pygmy services
  up          Bring up pygmy services (dnsmasq, haproxy, mailhog, resolv, ssh-agent)
  update      Pulls Docker Images and recreates the Containers
  version     # Check current installed version of pygmy

Flags:
      --config string   config file (default is $HOME/.pygmy.yml)
  -h, --help            help for pygmy
  -t, --toggle          Help message for toggle

Use "pygmy [command] --help" for more information about a command.
```



## Adding ssh keys

Call the `addkey` command with the **absolute** path to the key you would like to add. In case this they is passphrase protected, it will ask for your passphrase.

    pygmy addkey /Users/amazeeio/.ssh/my_other_key

    Enter passphrase for /Users/amazeeio/.ssh/my_other_key:
    Identity added: /Users/amazeeio/.ssh/my_other_key (/Users/amazeeio/.ssh/my_other_key)

## Checking the status

Run `pygmy status` and `pygmy` will tell you how it feels right now and which ssh-keys it currently has in it's stomach:

    pygmy status

    [*] pygmy-ssh-agent: Running as container pygmy-ssh-agent
    [*] mailhog.pygmy.site: Running as container mailhog.pygmy.site
    [*] pygmy-haproxy: Running as container pygmy-haproxy
    [*] pygmy-dnsmasq: Running as container pygmy-dnsmasq
    [*] pygmy-haproxy is connected to network pygmy-network
    [*] Resolv MacOS Resolver is properly connected
    �ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDNxWpKZcU/D+t7ToRGPNEXbvojrFtxKH99ZuaOJ7cs9KurVJyiEHyBEUZAPt0j9SO5yzdVEM//rVoZIwZeypW9C7CYgTpRoA/k1BnE1xvtoQT+528GmjQG542NBFo2KdO+LWqx19kClvoN7haGDtYKbS6MWUYEwD0ey69cquFDKC+A5NKx3z065gn9UZqLIeXjHCJ+v5PCSWXL3CFn57UlN824j1OFAECrjfNNfFEVmDJqa2Da6o9DhN+W1wyZJCklRPCiRlK5m3p9x1ClPKALUGQ0hvpjz36QSsXqS88MJPHsZvsv2PuW6xXNW8PSBCHcK6no5lYV/4hk8jcDQd2P6dpwvDiti+bTcfDH3jrVNqFati7ku37xIc3jWGn7CkCpMy008ai4kFMq2W2w6gOy0HncQ7z8AE8BdndxyEFYCLJviWOjW1SjSesPJpc9dxgmSmp/2qa6u0UZzFFHxJklIHepJAvcoHghs5Te2oMHwriRdpKqXiW+eJyudWCOzEeJljr73/Caft+CgZ7+kmmiy0hlqVAD6xkyBsuEF8+MdONfBHarpY8qZdLehavGd0DJW36nDnPvefDxoidJ0qYtjF8ElpNkeguAnsUFEwHkoc3Ur/NDcrkdGTKS8wb5AtkdwbDOCQTR00ABfAcYUFwOAvXodoQLrvm2ibp5l7/Y/Q== user@localhost
     - http://mailhog.pygmy.site (mailhog.pygmy.site)
     - http://pygmy.site/stats (pygmy-haproxy)

## `pygmy down` vs `pygmy clean`

`pygmy` behaves like Docker, it's a whale in the end!

During regular development `pygmy stop` is perfectly fine, it will remove the Docker containers still alive.

If you like to cleanup though, use `pygmy clean` to kill and remove all of the Docker containers, even if they're not alive.

## Access HAProxy statistic page and logs  

HAProxy service has statistics web page already enabled. To access the page, just point the browser to [http://pygmy.site/stats](http://pygmy.site/stats).  

To watch at haproxy container logs, use the `docker logs pygmy-haproxy` command with standard `docker logs` options like `-f` to follow.
