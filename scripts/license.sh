#!/usr/bin/env bash
set -eux

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

gpl="/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */"

agpl="/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024 Monomax Software Pty Ltd
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

dir=$(dirname "${BASH_SOURCE[0]}")
basedir="$dir/.."
pkgPath="$basedir/pkg"
serverPath="$basedir/pkg/server"

gplFiles=$(find "$pkgPath" -type f \( -name "*.go" -o -name "*.js" -o -name "*.ts" -o -name "*.tsx" -o -name "*.scss" -o -name "*.css"  \) ! -path "**/vendor/*" ! -path "**/node_modules/*" ! -path "$serverPath/*")

for file in $gplFiles; do
  remove_notice "$file"
  add_notice "$file" "$gpl"
done

agplFiles=$(find "$serverPath" -type f \( -name "*.go" -o -name "*.js" -o -name "*.ts" -o -name "*.tsx" -o -name "*.scss" -o -name "*.css" \) ! -path "**/vendor/*" ! -path "**/node_modules/*" ! -path "**/dist/*")

for file in $agplFiles; do
  remove_notice "$file"
  add_notice "$file" "$agpl"
done
