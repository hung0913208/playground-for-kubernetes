# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.boot_timeout = 1800

  # Prevent SharedFoldersEnableSymlinksCreate errors
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Configure network
  config.vm.network :private_network, type: "dhcp"

  for i in (1..2)
    config.vm.define "vm#{i}" do |node|
      node.vm.box = "ubuntu/trusty64"

      node.vm.provider :virtualbox do |vb|
        vb.customize [
          "modifyvm", :id, "--uart1", "0x3F8", "1"
        ]

        # Redirect console to file
        vb.customize [
          "modifyvm", :id, "--uartmode1", "file", "/tmp/vm#{i}.log"
        ]
      end
    end
  end
end
