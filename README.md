bin
===

`bin` is a binary installer that can extract binaries from archives and
automatically select binary for a platform.

TODO
====

- [x] `direct` extractor (for kubectl and sops)
- [ ] Add logger everywhere
- [ ] Tracking of installed binaries
    * Don't download from URL if already did it (override with force)
    * Implement new command `list` to show installed binaries with versions and
      URL used to install.
    * Implement new command `uninstall` command.
- [ ] Install via short github URL like `bin install sharkdp/fd`
    * Should discover latest release and choose correct asset by name
- [ ] Implement MacOS install
