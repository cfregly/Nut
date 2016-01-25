## Nut

Build [LXC](http://linuxcontainers.org/) containers using [Dockerfile](https://docs.docker.com/engine/reference/builder/) like [DSL](https://en.wikipedia.org/wiki/Domain-specific_language)

### Introduction

Nut is a minimal golang based CLI for building LXC containers. It allows user create containers
using [Dockerfile](https://docs.docker.com/engine/reference/builder/) like syntax and publishing them in s3. Nut is intended to be using in CI/CD
infrastructure to build container images as artifacts.

### Usage

```
nut build -help
```

```
Usage of build:
  -ephemeral
      Destroy the container after creating it
  -export string
      File path for the container tarball
  -export-sudo
      Use sudo while invoking tar
  -name string
      Name of the resulting container (defaults to randomly generated UUID)
  -publish string
      Publish tarball in s3 (assumes -stop, -export)
  -specfile string
      Container build specification file (default "Dockerfile")
  -stop
      Stop container after build
  -volume string
      Mount host directory inside container. Format: '[host_directory:]container_directory[:mount options]
```

#### Artifact  Labels

Nut stores container metadata as mainfest.yml file inside the container
directory, right next to rootfs directory. Manifest data stores all labels,
maintainers, exposed ports and entry point details. Labels that starts with "nut_artifact"
are treated differently, their values are considered as build artifacts and
fetched from inside the container to current directory. Following is an example
of building ruby 2.2.3 debian packages for trusty using nut

Example:

- Dockerfile
```sh
FROM trusty
MAINTAINER ranjib@pagerduty.com
RUN apt-get update -y
RUN apt-get install -y build-essential curl git-core libcurl4-openssl-dev libffi-dev libreadline-dev libsqlite3-dev libssl-dev libtool libxml2-dev libxslt1-dev libyaml-dev openssh-server python-software-properties sqlite3 wget zlib1g-dev
RUN mkdir -p /opt/ruby
RUN git clone https://github.com/sstephenson/ruby-build.git /opt/ruby-build
RUN /opt/ruby-build/bin/ruby-build 2.2.3 /opt/ruby/2.2.3
RUN /opt/ruby/2.2.3/bin/gem install bundler fpm --no-ri --no-rdoc
RUN /opt/ruby/2.2.3/bin/fpm -s dir -t deb -n ruby-2.2.3 -v 1.0.0 /opt/ruby/2.2.3
LABEL nut_artifact_ruby=/root/ruby-2.2.3_1.0.0_amd64.deb
```
And then nut can be invoked as:
```
nut build -ephemeral
```
Upon invocation nut will clone a new container from `trusty`, execute the RUN statement, which in turn will build ruby debian package, and then copy theresulting debian from /root/ruby-2.2.3_1.0.0_amd64.deb to current directory.

Since vanilla LXC is not aware of image repositories, all containers are created from cloning existing container(s).
A trusty (ubuntu 14.04) container can be created as
```
lxc-create -n trusty -t download -- -d ubuntu -a amd64 -r trusty
```
This in turn, can be used inside a Dockerfile DSL via the `FROM` instruction.
Nut converts the image name from  `org/repo:version` to 'org-repo_version' as the container
from which the new container will be built.
For example `FROM pagerduty/ruby:2.2.3` will instruct Nut to create a container by cloning
an existing container named `pagerduty-ruby_2.2.3`.


### Differences of LXC vs Docker runtime

Since Nut uses LXC it provides system container, which has few key differences from Docker/Rocket style app containers. These include:
- Container images/tarballs are independent artifacts, they do not have image or unionfs layer dependencies, and can be run without their parent images.
- Nut does not rely on overlayfs or any other union file system (whatever lxc supports)
- Nut is built and tested as unprivileged user. Hence `nut` can be run as any normal user (e.g. service X can create and run containers as X).
- LXC being a system conatiner runtime, does not force you to provide any specific entrypoint. You can model your app as standard linux service instead (init.d script, systmed unit file etc.
as you would do on host OS  (like using sys-v init script, or upstrat or systemd unit file etc). This also means you dont have to run an explicit process supervisor (like supervisord)
- Nut runs cron and a handful of additional services (exact list of services depends on the container distro), which means log rotation and all other periodic tasks for long lived containers do not involve anything special.


### Development

We use vagrant for development purpose. Following will setup a development vagrant instance
as well as kick a test job.

```
vagrant up
vagrant reload
vagrant ssh -c "nut build -specfile gopath/src/github.com/PagerDuty/nut/Dockerfile -ephemeral"
```

### LICENSE

[Apache 2](http://www.apache.org/licenses/LICENSE-2.0)
