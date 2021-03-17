# Chroot Commands

The following describes and outlines the various chroot related subcommands.

!!! info
    All of these commands require root privilege to run

## Exec

Execute a command in the chroot. Everything after the `--` is passed
and executed

```shell
imgchroot exec <img> -- <cmd>
```

Here is an example:

```shell
sudo ./imgchroot exec focal-server-cloudimg-amd64.img -- cat /etc/cloud/build.info
build_name: server
serial: 20210308
```

## Info

Print out information about the image and its' partitions. Optionally a user
can get the information in JSON:

```shell
imgchroot info <img> [--json]
```

Below is an example of an Ubuntu 20.04 LTS (focal) image:

```shell
$ sudo ./imgchroot info focal-server-cloudimg-amd64.img
focal-server-cloudimg-amd64.img
type: qcow2
size: 529 MiB
virtual size: 2252 MiB
partition table: gpt
sector size: 512
partitions:
  - name: nbd0p1
    type: Linux filesystem
    label: cloudimg-rootfs
    filesystem: ext4
  - name: nbd0p14
    type: BIOS boot
    label:
    filesystem:
  - name: nbd0p15
    type: EFI System
    label: UEFI
    filesystem: vfat
```

## Pull

To pull a file from the chroot use the pull subcommand:

```shell
imgchroot pull <img> <src> <dst>`
```

The below example demonstrates pulling the /etc/hosts file from the chroot
and saves it as hosts:

```shell
$ sudo ./imgchroot pull focal-server-cloudimg-amd64.img /etc/hosts hosts
$ cat hosts
127.0.0.1 localhost

# The following lines are desirable for IPv6 capable hosts
::1 ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
ff02::3 ip6-allhosts
```

## Push

To push a file to the chroot use the push subcommand:

```shell
imgchroot push <img> <src> <dst>
```

The following example shows how to put a file on the chroot and verify it
afterwards:

```shell
$ echo "myserver" > hostname
$ sudo ./imgchroot push focal-server-cloudimg-amd64.img hostname /etc/hostname
$ sudo ./imgchroot exec focal-server-cloudimg-amd64.img -- cat /etc/hostname
myserver
```

## Run

The run command is a short-cut subcommand, which will push the file, execute
it, and remove the file from the chroot:

```shell
imgchroot run <img> <script>
```

The following example runs a few commands:

```shell
$ cat myscript.sh
#!/bin/bash
whoami
pwd
ls
$ sudo ./imgchroot run focal-server-cloudimg-amd64.img myscript.sh
root
/
bin   dev  home  lib32  libx32      media  opt   root  sbin  srv  tmp  var
boot  etc  lib   lib64  lost+found  mnt    proc  run   snap  sys  usr
```

## Shell

The shell subcommand is used to launch a shell, by default bash, on the chroot:

```shell
imgchroot shell <img>
```

This is a quick way to run various commands or explore the chroot:

```shell
$ sudo ./imgchroot shell focal-server-cloudimg-amd64.img
root:/# echo $SHELL
/bin/bash
root:/# exit
exit
```
