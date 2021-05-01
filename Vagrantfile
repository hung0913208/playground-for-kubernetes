# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.boot_timeout = 1800

  # Prevent SharedFoldersEnableSymlinksCreate errors
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Configure network
  config.vm.network :private_network, type: "dhcp"

  for i in (1..2)
    ip = "192.168.0.#{i}"
    name = "vm#{i}"
    console = "/tmp/vm#{i}.log"

    config.vm.define name do |node|
      node.vm.box = "ubuntu/trusty64"

      node.vm.network :private_network, ip: ip, virtualbox__intnet: true

      node.vm.provider :virtualbox do |vb|
        vb.customize [
          "modifyvm", :id, "--uart1", "0x3F8", "1"
        ]

        # Redirect console to file
        vb.customize [
          "modifyvm", :id, "--uartmode1", "file", console
        ]
      end
    end
  end
end
