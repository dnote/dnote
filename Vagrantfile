# -*- mode: ruby -*-

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.synced_folder '.', '/go/src/github.com/dnote/dnote'
  config.vm.network "forwarded_port", guest: 3000, host: 3000
  config.vm.network "forwarded_port", guest: 8080, host: 8080
  config.vm.network "forwarded_port", guest: 5432, host: 5433

  config.vm.provision 'shell', path: './scripts/vagrant/install_utils.sh'
  config.vm.provision 'shell', path: './scripts/vagrant/install_go.sh', privileged: false
  config.vm.provision 'shell', path: './scripts/vagrant/install_node.sh', privileged: false
  config.vm.provision 'shell', path: './scripts/vagrant/install_postgres.sh', privileged: false
  config.vm.provision 'shell', path: './scripts/vagrant/bootstrap.sh', privileged: false

  config.vm.provider "virtualbox" do |v|
    v.memory = 4000
    v.cpus = 2
  end
end
