# Assets

Assets from BN6 should be placed in this folder. Due to copyright restrictions, we cannot distribute these files.

You must use two tools to extract this content. Eventually they'll be merged into one, but for now you'll have to run them separately.

## Using `murkland/bnrom/bndumper`

This will dump all graphical assets.

1. Run `go get github.com/murkland/bnrom/bndumper`

1. Run `bndumper path_to_rom.gba` in this directory. This should dump all the graphical assets into this folder.

## Using `murkland/agdbump`

This will dump all audio assets.

1. Run `git clone https://github.com/murkland/agbdump`.

1. Run `git submodule init && git submodule update`

1. Run `make -j4` to build `agbdump`.

1. Make a directory in the assets folder called `sounds` and run `agdbump path_to_rom.gba` in the `sounds` directory.

1. Encode all `.wav` files to `.ogg`. If `oggenc` is installed, you can just run `oggenc *.wav`.
