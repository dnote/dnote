# CHANAGELOG

All notable changes to the projects under this repository will be documented in this file.

* [Server](#server)
* [CLI](#cli)
* [Browser Extensions](#browser-extensions)

## Server

The following log documents the history of the server project.

### 1.0.1 - 2020-03-29

- Fix fresh install running migrations against tables that no longer exists.

### 1.0.0 - 2020-03-22

#### Fixed

- Fix unsubscribe link from the inactive reminder (#433)

#### Removed

- Remove the deprecated features related to digests and repetition rules (#432)
- Remove the migration for the deprecated, encrypted Dnote (#433)

### 0.5.0 - 2020-02-06

#### Changed

- **Deprecated** the digest and digest emails (#397)
- **Deprecated** the repetition rules (#397)

#### Fixed

- Fix refocusing to the end of the textarea input (#405)

### 0.4.0 - 2020-01-09

#### Added

- A web-based digest (#380)

#### Fixed

- Send inactive reminders with a correct email type (#385)
- Wrap words in note content (#389)

### 0.3.4 - 2019-12-24

#### Added

- Remind when the knowledge base stops growing (#375)
- Alert when a password is changed (#375)

#### Fixed

- Implement syntax highlighting for code blocks ($377)

### 0.3.3 - 2019-12-17

#### Added

- Send welcome email with login instructions upon reigstering (#352)
- Add an option to disable registration (#365)

#### Changed

- Send emails from the domain that hosts the application for on premise installations (#355)
- For on premise installations, automatically upgrade user accounts (#361)

### 0.3.2 - 2019-11-20

#### Fixed

- Fix server crash upon landing on a note page (#324).
- Allow to synchronize a large number of records (#321)

### 0.3.1 - 2019-11-12

#### Fixed

- Fix static files not being embedded in the binary. (#309)
- Fix mobile menu not covering the whole screen. (#308)

### 0.3.0 - 2019-11-12

#### Added

- Share notes (#300)
- Allow to recover from a missed repetition processing (#305)

### 0.2.1 - 2019-11-04

#### Upgrade Guide

* Please define the follwoing new environment variables:

  - `WebURL`: the URL to your Dnote server, without the trailing slash. (e.g. `https://my-server.com`) (Please see #290)
  - `SmtpPort`: the SMTP port. (e.g. `465`) optional - required *if you want to configure email*

#### Added

- Display version number in the settings (#293)
- Allow unsecure database connection in production (#276)

#### Fixed

- Allow to customize the app URL in the emails (#290)
- Allow to customize the SMTP port (#292)

### 0.2.0 - 2019-10-28

#### Added

- Specify spaced repetition rule (#280)

#### Changed

- Treat a linebreak as a new line in the preview (#261)
- Allow to have multiple editor states for adding and editing notes (#260)

#### Fixed

- Fix jumping focus on editor (#265)

### 0.1.1 - 2019-09-30

#### Fixed

- Fix asset loading (#257)


### 0.1.0 - 2019-09-30

#### Added

- Full-text search (#254)
- Password recovery (#254)
- Embedded notes in the digest emails (#254)

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Existing users need to go to `/classic` and follow the automated migration steps. (#254)
- **Breaking Change**: `v1` and `v2` API endpoints were removed, and `v3` API was added as a replacement.

#### Migration guide

- In your application, navigate to `/classic` and follow the automated migration steps.


## CLI

The following log documentes the history of the CLI project

### 0.11.0 - 2020-02-05

#### Added

- Allow to pass credentials through flags while logging in (#403)

### 0.10.0 - 2019-09-30

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Previous versions will no longer be able to interact with the web API, because `v1` and `v2` endpoints were replaced by a new `v3` endpoint to remove encryption.

#### Migration guide

- If you are using Dnote Pro, change the value of `apiEndpoint` in `~/.dnote/dnoterc` to `https://api.getdnote.com`.

## Browser Extensions

The following log documentes the history of the browser extensions project

### [Unreleased]

N/A

### 2.0.0 - 2019-10-29

- Allow to customize API and web URLs (#285)

### 1.1.1 - 2019-10-02

- Fix failing requests (#263)

### 1.1.0 - 2019-09-30

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Previous versions will no longer be able to interact with the web API, because `v1` and `v2` endpoints were replaced by a new `v3` endpoint to remove encryption.
