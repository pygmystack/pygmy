!!! warning
    This section is outdated and needs an update

# Local Drupal Docker Development

amazee.io supports development workflows which involve local development sites. We provide a Drupal Docker development environment that runs on your local computer.
It uses the exact same configuration for **all** services like on the amazee.io servers. This means:

* If the site runs locally, it also runs on production
* You can use the exact same `settings.php` file for local and production

**And the best:** You don't need to have any amazee.io account or site running in order to use the local development environment! Just install it, and experience all the benefits of amazee.io for free.

The Docker based Drupal Development environment consists of two parts:

### Part I: Shared Docker Containers

The shared docker containers for HAProxy and the SSH Agent, these are used by all other containers in order to properly work. They are started with `pygmy` for Linux & OS X.

### [Part II: Drupal Docker Containers](./drupal_site_containers.md)

The Docker Containers which will run Drupal. These are made to be copied into a Drupal root directory and to be started from there with `docker-compose`.[ Read how they are used](./drupal_site_containers.md)

## What it includes

The amazee.io Local Docker Drupal Development environment equips you with all the tools you need to develop your Drupal site locally:

* **Webserver:** Nginx
* **Frontend Caching:** Varnish
* **FastCGI Process Manager:** PHP-FPM
* **Server-side Scripting Language:** PHP
* **Database:** MariaDB
* **Search:** Apache Solr
* **Dependency Manager for PHP:** Composer
* NodeJS / NPM

For more information about software components used in the amazee.io Stack head over to the [Components](../architecture/components.md) overview page.

## How this works

Docker is super awesome and the perfect tool for local development. There are some hurdles though \(no worries, we have a solution for all of them\):

#### Exposed ports

If multiple Docker containers are exposing the same port it assigned a random port to the exposed port. In our case, this would mean, that each Drupal Container which would like to listen on Port 80 would get a random port like 34564 assigned. As they are random assigned it would be a lot of hassle of figuring out which port that the Drupal is found, additionally, Drupal doesn't like to run on another Port then 80 or 443 so much.

#### SSH Keys

It is possible to add mount ssh private keys into Docker containers, but this is again cumbersome, especially when you have a passphrase protected key \(as you should!\). You would need to enter the passphrase for each container that you start. Not a lot of fun.

### The Solution

amazee.io implemented a Drupal Docker Development environment which handles all these issues nicely for you. It allows you to:

* Access all sites via the Port 80 or 443 with just different URLs like site1.docker.amazee.io and site2.docker.amazee.io
* Add your SSH Key once to the system and can forget about it, no need to add it to each container

The environment starts 3 containers:

* [andyshinn/dnsmasq](https://hub.docker.com/r/andyshinn/dnsmasq/) Docker container which will listen on port 53 and resolve all DNS requests from `*.docker.amazee.io` to `127.0.0.1` \(so basically a better way then filling your `/etc/hosts` file by hand\)
* [amazeeio/haproxy](https://hub.docker.com/r/amazeeio/haproxy/) Docker container which will listen on port 80 and 443. It additionally listens to the Docker socket, realize when you start a new Drupal Container and adapt fully automatically it's haproxy configuration \(thanks to the awesome tool [docker-gen](https://github.com/jwilder/docker-gen)\). It forwards HTTP and HTTPs requests to the correct Drupal Container. With that we can access all Drupal Containers via a single Port.
* [amazeeio/ssh-agent](https://hub.docker.com/r/amazeeio/ssh-agent/) Docker container which will keeps an ssh-agent at hand for the other Drupal Containers. With that the Drupal Containers do not need to handle ssh-agenting themselves

#### Schema for Linux \(native Docker\)

```
                                            +--------------------------------------------------------------------+
                                            |Docker                                                              |
                                            |                                                                    |
                                            |          HAProxy knows which                                       |
                                            |          *.docker.amazee.io is                                     |
                                            |          handled by which container  +---------------------+       |
                                            |                                      |                     |       |
                                            |                              +-------+ Drupal Container 1  <--+    |
                                            |                              |       |                     |  |    |
+--------------------+                      |     +------------------+     |       +---------------------+  |    |
|                    |                      |     |                  |     |                                |    |
|                    |                      |     |     HAProxy      +-----+                                |    |
|                    +---------------------------->                  |     |       +---------------------+  |    |
|                    |                      |     | Published Ports  |     |       |                     |  |    |
|                    |  any HTTP/HTTPS      |     | 80/443           |     +-------+ Drupal Container 2  <--+    |
|                    |  request             |     +------------------+             |                     |  |    |
|      Browser       |                      |                                      +---------------------+  |    |
|                    |                      |                                                               |    |
|                    |                      |     +------------------+                                      |    |
|                    +---------------------------->                  |                                      |    |
|                    <----------------------------+   dns masq       |                                      |    |
|                    |                      |     |                  |                                      |    |
|                    |  Resolves            |     |                  |                                      |    |
|                    |  *.docker.amaze.io   |     |                  |                                      |    |
+--------------------+  to IP of Haproxy    |     +------------------+                                      |    |
                                            |                                                               |    |
                                            |                                                               |    |
                                            |                                                               |    |
                                            |                                                               |    |
                                            |                                                               |    |
+--------------------+                      |     +------------------+                                      |    |
|                    |                      |     |                  |                                      |    |
| pygmy              +---------------------------->    ssh agent     +--------------------------------------+    |
|                    |                      |     |                  |                                           |
+--------------------+  injects ssh-key     |     +------------------+  Exposes ssh-agent via                    |
                        into agent          |                           /tmp/amazeeio_ssh-agent/socket           |
                                            |                                                                    |
                                            +--------------------------------------------------------------------+
```

