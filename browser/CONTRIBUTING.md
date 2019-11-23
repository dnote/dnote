# Contributing

Use the following commands to set up, build, and release.

## Set up

* `npm install` to install dependencies.

## Developing locally

* `npm run watch:firefox`
* `npm run watch:chrome`

## Releasing

* Set a new version in `package.json`
* Run `./scripts/build_prod.sh`
  * A gulp task `manifest` will copy the version from `package.json` to `manifest.json`
