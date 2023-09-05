# `cidls`

`ls`, but with CIDs.

A cli command similar to `ls`, except with CIDs computed for each file name alongside the output. The current implementation is functional, but in an MVP state.

## Quick install

Quickly install this cli command by running
```bash
curl -fsSL https://raw.githubusercontent.com/orvn/cidls/main/install.sh | sudo bash
```

That's it! OS detection is automatic and the binary will be moved to an executable path. You'll likely be prompted for a password.

## Usage

`cidls` `[path]` `[CID version]`

- Run the command like `cidls ~/some/path 0`
- The path is optional, default path is the current working directory
- The second argument accepts `0` for CID v0 or `1` for CID v1


### Flags

- `cidls -v` for version information
- `cidls -h` for help

(more flags to come for setting multihash and multicodec)

### Build from source

1. Clone this this repo
2. Compile from source with `go build`
3. Try it out by running `./cidls` for the current directory or `./cidls ~/some/path` to target any directory


#### todo

- Compile and test on different OSs (currently only tested on macOS)
- Add to path to run as cli command
- Support different types of CIDs
- Create a caching system to avoid re-processing files too much
- Make compatible with BSD-style `LSCOLORS` variable (e.g., `exgxcxdxbxegedabagacad`)
- Add options for a few common ls flags
