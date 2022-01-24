# -*- mode: ruby -*-
# vi: set ft=ruby :

CLUSTER_SIZE = 3

Vagrant.configure("2") do |config|
  config.vm.boot_timeout = 1800

  # Prevent SharedFoldersEnableSymlinksCreate errors
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Configure network
  config.vm.network :private_network, type: "dhcp"

  (1..CLUSTER_SIZE).each do |vid|
    config.vm.define "vm#{vid}" do |node|
      node.vm.network :private_network, ip: "192.168.0.#{vid}",
                      virtualbox__intnet: true

      node.vm.hostname = "node#{vid}"
      node.vm.provider :virtualbox do |vb|
        node.vm.box = "generic/ubuntu1804"

        vb.memory = 2048
        vb.cpus = 2

        # Enable using uart port 1
        vb.customize [
          "modifyvm", :id, "--uart1", "0x3F8", "1"
        ]

        # Redirect console to file
        vb.customize [
          "modifyvm", :id, "--uartmode1", "file", "/tmp/vm#{vid}.log"
        ]
      end  

      node.vm.provision "ansible" do |ansible|
        if vid == 1 then
          ansible.playbook = "kubernetes/setup/master-playbook.yml"
            ansible.extra_vars = {
              node_ip: "192.168.0.#{vid}",
            }
        else
          ansible.playbook = "kubernetes/setup/worker-playbook.yml"
          ansible.extra_vars = {
            node_ip: "192.168.0.#{vid}",
          }
        end
      end

      if vid == CLUSTER_SIZE then
        node.trigger.after :up do |trigger|
          trigger.info = "Start playground"

          Dir.glob('kubernetes/playground/*.yml') do |yaml|
            node.vm.provision "ansible" do |ansible|
              ansible.playbook = yaml
              ansible.extra_vars = {
                node_ip: "192.168.0.#{vid}",
              }
            end
          end
        end
      end
    end
  end
end
