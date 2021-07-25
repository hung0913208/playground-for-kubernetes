# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.boot_timeout = 1800

  # Prevent SharedFoldersEnableSymlinksCreate errors
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Configure network
  config.vm.network :private_network, type: "dhcp"

  (1..2).each do |vid|
    config.vm.define "vm#{vid}" do |node|
      node.vm.network :private_network, ip: "192.168.0.#{vid}",
                      virtualbox__intnet: true

      node.vm.provider :virtualbox do |vb|
        node.vm.box = "generic/ubuntu1804"

        # Enable using uart port 1
        vb.customize [
          "modifyvm", :id, "--uart1", "0x3F8", "1"
        ]

        # Redirect console to file
        vb.customize [
          "modifyvm", :id, "--uartmode1", "file", "/tmp/vm#{vid}.log"
        ]
      end
    end
  end
end
