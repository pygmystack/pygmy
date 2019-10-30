# Update pygmy

As `pygmy` is released regularly you should also take care of updating pygmy.

As ruby gems does not remove old versions of gems when updating, you should remove old version after a new version has been installed.

First update:

    gem update pygmy

    Updating installed gems
    Updating pygmy
    Successfully installed pygmy-0.9.4
    Parsing documentation for pygmy-0.9.4
    Installing ri documentation for pygmy-0.9.4
    Installing darkfish documentation for pygmy-0.9.4
    Done installing documentation for pygmy after 0 seconds
    Parsing documentation for pygmy-0.9.4
    Done installing documentation for pygmy after 0 seconds
    Gems updated: pygmy

A new version!

Now uninstall the old one:

    gem uninstall pygmy

    Select gem to uninstall:
     1. pygmy-0.9.3
     2. pygmy-0.9.4
     3. All versions
    > 1
    Successfully uninstalled pygmy-0.9.3

check the correct version:

    pygmy -v

    Pygmy - Version: 0.9.4


## I see errors after the upgrade

If you see following error after the pygmy upgrade you might have installed 2 pygmy versions next to eachother
```
> pygmy version
/Users/user1/.rvm/gems/ruby-2.2.1/gems/pygmy-0.9.11/lib/pygmy/version.rb:2: warning: previous definition of VERSION was here
/Users/user1/.rvm/gems/ruby-2.2.1/gems/pygmy-0.9.11/lib/pygmy/version.rb:3: warning: previous definition of DATE was here
/Users/user1/.rvm/gems/ruby-2.2.1/gems/pygmy-0.9.10/lib/pygmy/version.rb:2: warning: already initialized constant Pygmy::VERSION
/Users/user1/.rvm/gems/ruby-2.2.1/gems/pygmy-0.9.10/lib/pygmy/version.rb:3: warning: already initialized constant Pygmy::DATE
```

Run `gem uninstall pygmy` and remove the old versions.

## Update Docker Containers with `pygmy`

`pygmy` can update shared docker containers for you:

    pygmy update

After it updated all containers, it will recreate them as well.

