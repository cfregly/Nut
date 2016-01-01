## Nut

Build LXC containers using Dockerfile like DSL

### Usage

```
nut -help
```

```
Usage: nut [options]

  Build containers using LXC runtime with pluggable build DSLs

Options:

  -help        Show usage
  -specfile    Local path to the specification file (defaults to dockerfle)
  -ephemeral   Destroy the container after creation
  -name        Name of the container (defaults to randomly generated UUID)
```

#### Artifact  Labels

Nut stores container metadata as mainfest.yml file inside the container
directory, right next to rootfs directory. Manifest data stores all labels,
maintainers and exposed ports details. Labels that starts with "nut_artifact"
are treated differently, their values are considered as build artifacts and
fetched from inside the container to current directory. Following is an example
of building ruby 2.2.3 debian packages for trusty using nut

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
nut -ephemeral
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

Since Nut uses LXC it has few key differences from Docker/Rocket style app containers. These include:
- Once built, they do not have image dependency. They can be run without
their parent images.
- They also dont require overlayfs or any other union file system
- Can be run as any normal user (e.g. service X can create and run containers as X).
- Since these are full system conatiners, you dont have to provide any specific entrypoint. You can daemonize your app
as you would do on host OS  (like using sys-v init script, or upstrat or systemd unit file etc).
- Since these are full system containers, you also wont have to reap zombie processes explicitly, init system will do it for you


### Development

We use vagrant for development purpose. Following will setup a development vagrant instance
as well as kick a test job.

```
vagrant up
vagrant reload
vagrant ssh -c "nut -specfile gopath/src/github.com/PagerDuty/nut/Dockerfile -ephemeral"
```

### LICENSE

[Apache 2](http://www.apache.org/licenses/LICENSE-2.0)
