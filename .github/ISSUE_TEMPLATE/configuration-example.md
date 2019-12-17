---
name: Configuration example
about: Request an example configuration for a new container.
title: Request for example pygmy ______ integration
labels: ''
assignees: fubarhouse

---

**Background**

Reason for asking and what you're interested in achieving with this. This will serve as a reference to other users making similar requests if yours is more suitable.

**Implementation**

__Docker__

The docker command is required for this work. Please identify the full docker command in which you normally use to run the container. This is required and your request will likely be delayed without a reasonable explanation of why it is missing.

```
docker run -it --rm mbentley/cowsay holy ship!
```

__Variables__

Variables which would need injection to the container at start-up.

```
MYVAR1=hello
MYVAR2=world!
```

__Volumes__

If any named volumes need to be present for this container please name them here - pygmy will create them but managing and removing them will be your responsibility.

```
- myservice_filesystem
- myservice_database
```
