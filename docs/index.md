# Pygmy

`pygmy` is the single tool needed to get the local [amazee.io](https://amazee.io) Docker Drupal Development Environment running on your Linux based system. It built to work with [Docker for Mac](https://docs.docker.com/docker-for-mac/)! (quite a lot for such a [small whale](https://en.wikipedia.org/wiki/Pygmy_sperm_whale) 🐳)

**What `pygmy` will handle for you:**

* Starting the necessary Docker Containers for the amazee.io Drupal Docker Development
* If on Linux: Adds `nameserver 127.0.0.1` to your `/etc/resolv.conf` file, so that your local Linux can resolve `*.docker.amazee.io` via the dnsmasq container
* If on Mac with Docker for Mac: Creates the file `/etc/resolver/docker.amazee.io` which tells OS X to forward DNS requests for `*.docker.amazee.io` to the dnsmasq container
* Tries to add the ssh key in `~/.ssh/id_rsa` to the ssh-agent container (no worries if that is the wrong key, you can add more any time)
* Starts a local mail Mail Transfer Agent (MTA) in order to test and view mails

