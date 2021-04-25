# Config Data

## Ordering

```text
* earlycmd
* archives
    sources
    ppa
* packages
    install
    remove
    upgrade
    proposed
* unattended-upgrades
* snap
* files
* groups
* users
* cmd
```

## Full Example

```yaml
version: 1

early-commands:
    - ["ls"]
    - ["uname", "-a"]
commands:
    - ["ls"]
    - ["uname", "-a"]
files:
    - encoding: b64
      content: CiMgVGhpcyBmaWxlIGNvbnRyb2xzIHRoZSBzdGF0ZSBvZiBTRUxpbnV4...
      owner: root:root
      path: /etc/sysconfig/selinux
      permissions: '0644'
    - content: |
        # My new /etc/sysconfig/samba file

        SMBDOPTIONS="-D"
      path: /etc/sysconfig/samba
groups:
    - tester
archives:
    sources:
        - url: "http://us.archive.ubuntu.com/ubuntu/ focal main"
          comment: "Added by..."
          remove: true
    ppa:
        - name: ppa:cloud-init-dev/daily
        - name: ppa:ubuntu-advantage/security-benchmarks
          creds: powersj:aaaaaa
          key: 0xA166877412DAC26E73CEBF3FF6C280178D13028C
          remove: true
          comment: "Added by..."
packages:
    install:
        - vim
    remove:
        - emacs
    proposed:
        - vim
    pin:
        - name: vim
          pin: release a=$LB_DISTRIBUTION-proposed
          priority: 400
unattended-upgrades:
    origins:
        - ${distro_id}:${distro_codename}-security
    blacklist:
        - vim
        - libc6
snap:
    - name: eksctl
      channel: 1.18.9/stable
users:
    - username: myaccount
      password: mypasswd
      ssh-import-id:
        - lp:powersj
        - gh:powersj
      ssh-authorized-keys:
        -  ssh-rsa AAAAB3NzaC1...POt5Q8zWd9qG7PBl9+eiH5qV7NZ mykey@host
      sudoer: true
      groups: [tester]
```
