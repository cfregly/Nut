# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = 'ubuntu/trusty64'

  config.vm.synced_folder '.', '/home/vagrant/gopath/src/github.com/PagerDuty/nut'

  $script = <<-SCRIPT
    sudo chown -R vagrant:vagrant /home/vagrant/gopath
    wget -q https://storage.googleapis.com/golang/go1.5.2.linux-amd64.tar.gz
    tar -zxvf go1.5.2.linux-amd64.tar.gz
    echo 'export GOROOT="/home/vagrant/go"' >  /tmp/gopath.sh
    echo 'export GOPATH="/home/vagrant/gopath"' >> /tmp/gopath.sh
    echo 'export PATH="/home/vagrant/go/bin:/home/vagrant/gopath/bin:\$PATH"' >> /tmp/gopath.sh
    sudo mv /tmp/gopath.sh /etc/profile.d/gopath.sh
    sudo chmod 0755 /etc/profile.d/gopath.sh
    source /etc/profile.d/gopath.sh
    cd /home/vagrant/gopath/src/github.com/PagerDuty/nut && go get ./... && go install ./...
  SCRIPT

  config.vm.provision 'shell', privileged: true, inline: <<-SHELL
    apt-get update
    apt-get install -y \
      build-essential \
      git \
      lxc \
      lxc-dev \
      pkg-config
  SHELL

  config.vm.provision 'shell', privileged: false, inline: $script

  lxc_config = <<-CONFIG
    lxc.network.type = veth
    lxc.network.link = lxcbr0
    lxc.network.flags = up
    lxc.network.hwaddr = 00:16:3e:xx:xx:xx
    lxc.id_map = u 0 100000 65536
    lxc.id_map = g 0 100000 65536
  CONFIG

  lxc_usernet = 'vagrant veth lxcbr0 10'

  config.vm.provision 'shell', privileged: false, inline: <<-SHELL
    mkdir -p /home/vagrant/.config/lxc /home/vagrant/.local/share/lxc /home/vagrant/.cache/lxc
    echo "#{lxc_config}" | tee /home/vagrant/.config/lxc/default.conf
    echo "#{lxc_usernet}" | sudo tee /etc/lxc/lxc-usernet
    echo "#{lxc_usernet}" | tee /home/vagrant/.config/lxc/lxc.conf
  SHELL
end
