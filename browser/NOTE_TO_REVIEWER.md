# Note to reviewer

This document contains instructions about how to reproduce the final build of this extension.

All releases are tagged and pushed to [the GitHub repository](https://github.com/dnote/dnote).

## Steps

To reproduce the obfuscated code for Firefox, please follow the steps below.

1.  Run `npm install` to install dependencies
2.  Run `./scripts/build_prod.sh` to build for Firefox and Chrome.

The obfuscated code will be under `/dist/firefox` and `/dist/chrome`.

## Further questions

Please contact sung@dnote.io
