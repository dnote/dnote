#!/usr/bin/env bash
set -eux

function has_license {
  # Check if file already has a copyright notice
  grep -q "Copyright.*Dnote Authors" "$1"
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

year=$(date +%Y)
license="/* Copyright $year Dnote Authors
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
  if ! has_license "$file"; then
    add_notice "$file" "$license"
  fi
done
