#!/bin/bash

function remove_notice {
  sed -i -e '/\/\* Copyright/,/\*\//d' "$1"

  # remove leading newline
  sed -i '/./,$!d' "$1"
}

function add_notice {
  ed "$1" <<END
0i
$2

.
w
q
END
}

gpl="/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */"

agpl="/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */"

cliPath="$GOPATH"/src/github.com/dnote/dnote/cli
cliFiles=$(find "$cliPath" -type f -name "*.go" ! -path "**/vendor/*")

for file in $cliFiles; do
  remove_notice "$file"
  add_notice "$file" "$gpl"
done

 serverPath="$GOPATH"/src/github.com/dnote/dnote/server
 webPath="$GOPATH"/src/github.com/dnote/dnote/web
 agplFiles=$(find "$serverPath" "$webPath" -type f \( -name "*.go" -o -name "*.js" \) ! -path "**/vendor/*" ! -path "**/node_modules/*")
 
 for file in $agplFiles; do
   remove_notice "$file"
   add_notice "$file" "$agpl"
 done
