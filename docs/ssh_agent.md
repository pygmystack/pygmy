# SSH Agent

Per default your SSH Key at `~/.ssh/id_rsa` is added to the Docker containers from `pygmy`

If you need another key, read the documentation of [`pygmy`](linux_pygmy.md) about this.

## How it works
1. `pygmy` starts `amazeeio/ssh-agent` container with a volume `/tmp/amazeeio_ssh-agent`
2. `pygmy` adds a default SSH key from the host into this volume
3. `docker-compose.yml` should have volume inclusion specified for CLI container:
  ```
  volumes_from:
    - container:amazeeio-ssh-agent
  ```
4. When CLI container starts, the volume is mounted and an entrypoint script adds SHH key into agent.
  @see https://github.com/amazeeio/lagoon/blob/master/images/php/cli/10-ssh-agent.sh

Running `ssh-add -L` within CLI container should show that the SSH key is correctly loaded.

## Troubleshooting
### SSH Key issues

As everything on amazee.io works with key authentication sometimes you might run into issues where the drush aliases aren't displayed or you can't connect to the servers.

    Could not load API JWT Token, error was: 'lagoon@ssh.lagoon.amazeeio.cloud: Permission denied (publickey).'

Or for legacy systems:

     drupal@example.amazee.io:~/public_html/docroot (staging)$ drush @master ssh
     Permission denied (publickey).

1. Check if you see the SSH Key inside your container with `ssh-add -L` <br>
   If you get `Could not open a connection to your authentication agent.` or `The agent has no identities.` head straight to **step 3.**
2. Check if you see your SSH Key in `pygmy status`
3. If you don't see the key in `pymgy status` run `pygmy addkey`. You should see `Successfully added ssh key` if the key addition was successful.
4. After that you need to recreate the containers `docker-compose up -d --force`
5. When the containers are recreated you should be able to see your ssh key with `ssh-add -L`
6. If you still get the `Permission denied (publickey)` error get in touch with our engineers to check if the key is configured correctly on the hosting side.
