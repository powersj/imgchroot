# Download

The download subcommand will download an image matching your system's
architecture:

```shell
imgchroot download <distro> <release>
```

Currently the command only supports `ubuntu` as a valid distro.

The download can resume an previous attempt if it was cancelled or connection
is lost. Download will not re-download an image if it already exists.

Below is example output:

```shell
$ imgchroot download ubuntu bionic
downloading https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img
 346 of  346 MB (100.0%) @   9.2 MB/sec
image saved as bionic-server-cloudimg-amd64.img
```
