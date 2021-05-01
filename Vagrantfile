Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/trusty64"
    config.vm.define 'ubuntu'

    # Vagrant boot needs more time on AppVeyor (see https://help.appveyor.com/discussions/problems/1247-vagrant-not-working-inside-appveyor)
    config.vm.boot_timeout = 1800

    # Prevent SharedFoldersEnableSymlinksCreate errors
    config.vm.synced_folder ".", "/vagrant", disabled: true

    config.vm.provider :virtualbox do |vb|
        vb.name = 'ubuntu'
        vb.customize [
          "modifyvm",
          :id,
          "--uart1",
          "0x3F8",
          "1"
        ]

        vb.customize [
          "modifyvm",
          :id,
          "--uartmode1",
          "file",
          "/tmp/vm1.log"
        ]
    end

end
