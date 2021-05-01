# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/trusty64"
    config.vm.define 'ubuntu'

    # Vagrant boot needs more time on AppVeyor
    config.vm.boot_timeout = 1800

    # Prevent SharedFoldersEnableSymlinksCreate errors
    config.vm.synced_folder ".", "/vagrant", disabled: true

    # Configure to redirect console log to file
    config.vm.provider :virtualbox do |vb|
      vb.name = 'ubuntu'
      vb.customize [
        "modifyvm",
        :id,
        "--uart1",
        "0x3F8",
        "1"
      ]
    end

    for node in (1..2)
      config.vm.define "vm#{node}" do |node|
        node.vm.customize [
          "modifyvm", :id, "--uartmode1", "file", "/tmp/vm#{node}.log"
        ]

        node.vm.network :private_network, ip: "192.168.0.#{node}"
      end
    end
end
