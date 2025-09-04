# terraform-provider-objectscale


### Installation from source code

1. [Install Rust](https://www.rust-lang.org/tools/install)
2. In the root of this repo run:
    ```shell
    make client-build
    make
    ```

### Release

To generate the release files:
```shell
make release
gpg --detach-sign terraform-provider-objectscale_*_SHA256SUMS
```
