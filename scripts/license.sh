#!/usr/bin/env bash
set -eux

function remove_notice {
  # Remove old copyright notice - matches /* Copyright ... */ including the trailing newline
  # The 's' flag makes . match newlines, allowing multi-line matching
  # The \n? matches an optional newline after the closing */
  perl -i -0pe 's/\/\* Copyright.*?\*\/\n?//s' "$1"
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

license="/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */"

dir=$(dirname "${BASH_SOURCE[0]}")
basedir="$dir/.."
pkgPath="$basedir/pkg"

# Apply license to all source files
allFiles=$(find "$pkgPath" -type f \( -name "*.go" -o -name "*.js" -o -name "*.ts" -o -name "*.tsx" -o -name "*.scss" -o -name "*.css"  \) ! -path "**/vendor/*" ! -path "**/node_modules/*" ! -path "**/dist/*")

for file in $allFiles; do
  remove_notice "$file"
  add_notice "$file" "$license"
done
