# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
    for node in (1..2)
      config.vm.define "vm#{node}" do |node|
        node.vm.box = "ubuntu/trusty64"

        # Vagrant boot needs more time on AppVeyor
        node.vm.boot_timeout = 1800

        # Prevent SharedFoldersEnableSymlinksCreate errors
        node.vm.synced_folder ".", "/vagrant", disabled: true

        # Redirect console to file
        node.vm.provider :virtualbox do |vb|
          vb.name = "node#{node}"
          vb.customize [
            "modifyvm", :id, "--uart1", "0x3F8", "1"
          ]

          vb.customize [
            "modifyvm", :id, "--uartmode1", "file", "/tmp/vm#{node}.log"
          ]
        end

        # Configure network
        node.vm.network :private_network, ip: "192.168.0.#{node}"
      end
    end
end
