# `cidls`

`ls`, but with CIDs.

A cli command similar to `ls`, except with CIDs computed for each file name alongside the output. The current implementation is functional, but in an MVP state.

## Usage

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
