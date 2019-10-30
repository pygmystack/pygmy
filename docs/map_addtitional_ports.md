!!! warning
    This section is outdated and needs an update

# Map additional ports

As it explained in the [Connect to MySQL](./connect_to_mysql_from_external.md) section, Docker maps MySQL's 3306 port to a random port on docker.amazee.io. This happens because port 3306 is set in the `docker-compose.yml` file:
```
    ports:
      - "3306"
```

If you need to map other ports, simply add them to the `ports` section and restart the container.

### Example for Solr

If you use one of amazee.io Drupal containers with Solr included, your Solr URL most likely looks like this: http://127.0.0.1:8149/solr/drupal/

In this case, to play with Solr queries:
- add `"8149"` to the `ports` section of `docker-compose.yml` file
- restart the container with `docker-compose stop && docker-compose up -d`
- get the port number with `docker-compose port drupal 8149`
- start playing at http://docker.amazee.io:&lt;PORT_NUMBER&gt;/solr/drupal/admin/
