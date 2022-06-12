# -*- mode: ruby -*-
# vi: set ft=ruby :

INSTALL_K8S = false
INSTALL_PXE = false
INSTALL_NETDEV = true
CLUSTER_SIZE = 3

Vagrant.configure("2") do |config|
  config.vm.boot_timeout = 1800

  # Prevent SharedFoldersEnableSymlinksCreate errors
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Configure network
  config.vm.network :private_network, type: "dhcp"

  (1..CLUSTER_SIZE).each do |vid|
    config.vm.define "vm#{vid}" do |node|
      node.vm.hostname = "node#{vid}"
      node.vm.provider :virtualbox do |vb|
        if INSTALL_PXE == true then
          if vid == 1 then 
            node.vm.box = "generic/ubuntu1804"
          else
            node.vm.box = "pace/empty"
          end
        else
          node.vm.box = "generic/ubuntu1804"
        end

        if INSTALL_K8S then
          vb.memory = 2048
          vb.cpus = 2
        else
          vb.memory = 512
          vb.cpus = 2
        end

        # Enable using uart port 1
        vb.customize [ "modifyvm", :id, "--uart1", "0x3F8", "1" ]

        # Redirect console to file
        vb.customize [ "modifyvm", :id, "--uartmode1", "file", "/tmp/vm#{vid}.log" ]

        if INSTALL_K8S == true or INSTALL_NETDEV == true then
          node.vm.network :private_network,
            ip: "192.168.0.#{vid}",
            virtualbox__intnet: true
        elsif vid == 1 then
          node.vm.network :private_network,
            ip: "192.168.0.#{vid}",
            virtualbox__intnet: true
        else
          vb.customize ['modifyvm', :id, '--boot1', 'net']
          vb.customize ['modifyvm', :id, '--boot2', 'disk']
          vb.customize ['modifyvm', :id, '--biospxedebug', 'on']
          vb.customize ['modifyvm', :id, '--cableconnected2', 'on']
          vb.customize ['modifyvm', :id, '--nicbootprio2', '1']
          vb.customize ['modifyvm', :id, "--nictype2", '82540EM'] 

          node.vm.network :private_network,
            ip: "192.168.0.0",
            auto_config: false,
            virtualbox__intnet: true

          config.vm.provision :shell, run: 'always', inline: "ip route list 0/0 | xargs ip route del; ip route add default via 192.168.0.1"
        end
      end

      if INSTALL_K8S == true then
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
              
            node.vm.provision "ansible" do |ansible|
              ansible.playbook = "kubernetes/wait.yml"
              ansible.extra_vars = {
                node_ip: "192.168.0.#{vid}",
              }
            end
          end
        end
      end # install-k8s

      if INSTALL_NETDEV == true then 
        Dir.glob('netdev/setup/*.yml') do |yaml|
          node.vm.provision "ansible" do |ansible|
            ansible.extra_vars = {
              node_ip: "192.168.0.#{vid}",
            }
          end
        end

      end # install-netdev
    end # node template
  end
end
