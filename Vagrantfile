# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.boot_timeout = 1800

  # Prevent SharedFoldersEnableSymlinksCreate errors
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Configure network
  config.vm.network :private_network, type: "dhcp"

  config.vm.define "vm1" do |node|
    node.vm.box = "ubuntu/trusty64"
    node.vm.network :private_network, ip: "192.168.0.1",
                    virtualbox__intnet: true

    node.vm.provider :virtualbox do |vb|
      vb.customize [
        "modifyvm", :id, "--uart1", "0x3F8", "1"
      ]

      # Redirect console to file
      vb.customize [
        "modifyvm", :id, "--uartmode1", "file", "/tmp/vm1.log"
      ]
    end
  end

  config.vm.define "vm2" do |node|
    node.vm.box = "ubuntu/trusty64"
    node.vm.network :private_network, ip: "192.168.0.2",
                    virtualbox__intnet: true

    node.vm.provider :virtualbox do |vb|
      vb.customize [
        "modifyvm", :id, "--uart1", "0x3F8", "1"
      ]

      # Redirect console to file
      vb.customize [
        "modifyvm", :id, "--uartmode1", "file", "/tmp/vm2.log"
      ]
    end
  end
end
