#!/bin/bash
# dev.sh starts a development environment for the api.
# usage: DOTENV_PATH=path_to_dotenv_file EMAIL_TEMPLATE_DIR=path_to_email_templates ./dev.sh
set -eux

basePath="$GOPATH/src/github.com/dnote/dnote/server"

# load env
export $(cat "$DOTENV_PATH" | xargs)

PORT=5000 CompileDaemon \
  -directory="$basePath"/api \
  -command="$basePath/api/api --emailTemplateDir $basePath/mailer/templates/src"
