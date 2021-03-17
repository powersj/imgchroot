# imgchroot

*Quickly interact and customize a cloud image*

[![Build Status](https://travis-ci.com/powersj/imgchroot.svg?branch=master)](https://travis-ci.com/powersj/imgchroot/) [![Go Report Card](https://goreportcard.com/badge/github.com/powersj/imgchroot)](https://goreportcard.com/report/github.com/powersj/imgchroot) [![Go Reference](https://pkg.go.dev/badge/github.com/powersj/imgchroot.svg)](https://pkg.go.dev/github.com/powersj/imgchroot)

## Overview

imgchroot is a Go-based CLI to quickly customize cloud images in a chroot
without the need to boot the image or setup a user with credentials.

imgchroot mounts the image to a temporary directory, using the network
block device (NBD) protocol. It then runs the required operation, such as a
command, moving a file, or starts a shell in the chroot. Once the operation is
complete the image is unmounted all without needing to boot the image itself.

## CLI

Here are the primary functions available via the CLI:

1. **Chroot** commands will run commands against the image inside a chroot.
   This includes command execution, via a shell, or push and pull files from
   the image.
1. **Download** latest cloud images

### Chroot

imgchroot provides a number of different sub-commands to directly interact
with a cloud image. Click to see the sub-command's CLI page for more details:

* [exec](https://powersj.github.io/imgchroot/chroot/#exec): run a command on
  the image
* [info](https://powersj.github.io/imgchroot/chroot/#info): information about
  the image
* [pull](https://powersj.github.io/imgchroot/chroot/#pull): pull a file from
  the image
* [push](https://powersj.github.io/imgchroot/chroot/#push): push a file to the
  image
* [run](https://powersj.github.io/imgchroot/chroot/#run): transfer and run a
  file on the image
* [shell](https://powersj.github.io/imgchroot/chroot/#shell): start a shell on
  the image

### Download

imgchroot has the ability to find the latest cloud images as well. A user
needs to provide the distro (e.g. ubuntu) and release (e.g. focal) to download.
The the download will find the URL and download it.

See the [download](https://powersj.github.io/imgchroot/download) sub-command
for more information.

## Support

If you find a bug, have a question, or ideas for improvements please file an
[issue](https://github.com/powersj/imgchroot/issues/new) on GitHub.
