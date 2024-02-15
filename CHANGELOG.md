# CHANAGELOG

All notable changes to the projects under this repository will be documented in this file.

- [Server](#server)
- [CLI](#cli)

## CLI

- `search` not matches partial words.
- `view` Accepts a list of ids and shows the content of all them.

## Server

The following log documents the history of the server project.

### Unreleased

None

### 2.1.1 2023-03-04

#### Fixed

- Added the missing CSS and JS in the server release

### 2.1.0 2023-03-04

#### Changed

- `OnPremise` environment variable is deprecated and is replaced with `OnPremises`
- Upgrade Go from 1.17 to 1.20.

### 2.0.0 2022-05-09

#### Removed

- The web interface for managing notes and books (#594)

### 1.0.4 2020-05-23

#### Removed

- Simplify the bundle by removing unnecessary payment logic

#### Fixed

- Fix timestamp in the note content view
- Invalidate existing sessions when password is changed

### 1.0.3 2020-05-03

#### Fixed

- Fix timeline grouping notes by added time rather than updated time.

#### Changed

- Sort notes by last activity to make it easier to see the most recently accessed information.

### 1.0.2 2020-05-03

#### Changed

- Support arm64.

### 1.0.1 - 2020-03-29

- Fix fresh install running migrations against tables that no longer exists.

### 1.0.0 - 2020-03-22

#### Fixed

- Fix unsubscribe link from the inactive reminder (#433)

#### Removed

- Remove the deprecated features related to digests and repetition rules (#432)
- Remove the migration for the deprecated, encrypted Dnote (#433)

#### Changed

- Please set `OnPremise` environment to `true` in order to automatically use the Pro version.

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

- Implement syntax highlighting for code blocks (\$377)

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

- Please define the follwoing new environment variables:

  - `WebURL`: the URL to your Dnote server, without the trailing slash. (e.g. `https://my-server.com`) (Please see #290)
  - `SmtpPort`: the SMTP port. (e.g. `465`) optional - required _if you want to configure email_

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

### Unreleased

None

### 0.15.1 - 2024-02-03

- Upgrade `color` dependency (#660).
- Use Go 1.21 (#658).

### 0.15.0 - 2023-05-27

- Add `enableUpgradeCheck` configuration to allow to opt out of automatic update check.

### 0.14.0 - 2023-03-10

- Remove `autocomplete` subcommand that was accidentally added by a dependency (#637)

### 0.13.0 - 2023-02-10

- Allow to add note from stdin.

```
echo "test" | dnote add mybook

dnote add mybook << EOF
test line 1
test line 2
EOF
```

### 0.12.0 - 2020-01-03

#### Upgrade guide

- **On Linux or macOS** Please move your Dnote files to new directories based on the XDG base directory specfication. **On Windows**, no action is required.

```
# Move the database file
mv ~/.dnote/dnote.db ~/.local/share/dnote/dnote.db

# Move the config file
mv ~/.dnote/dnoterc ~/.config/dnote/dnoterc

# Delete ~/.dnote. (it is safe to delete DNOTE_TMPCONTENT.md files, if they exist.)
rm -rf ~/.dnote
```

If `~/.dnote` directory exists, dnote will continue to use that directory for backward compatibility until the next major release.

#### Added

- Add `--content-only` flag to print the note content only (#528)

#### Changed

- Use XDG base directory on Linux and macOS (#527)

### 0.11.1 - 2020-04-25

#### Fixed

- Fix upgrade URL (#453)

#### Changed

- Display hostname of the self-hosted instance while logging in (#454)
- Display helpful error if endpoint is misconfigured (#455)

### 0.11.0 - 2020-02-05

#### Added

- Allow to pass credentials through flags while logging in (#403)

### 0.10.0 - 2019-09-30

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Previous versions will no longer be able to interact with the web API, because `v1` and `v2` endpoints were replaced by a new `v3` endpoint to remove encryption.

#### Migration guide

- If you are using Dnote Pro, change the value of `apiEndpoint` in `~/.dnote/dnoterc` to `https://api.getdnote.com`.
