/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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
 */

package tmpl

var noteMetaTags = `
<meta name="description" content="{{ .Description }}" />
<meta name="twitter:card" content="summary" />
<meta name="twitter:title" content="{{ .Title }}" />
<meta name="twitter:description" content="{{ .Description }}" />
<meta name="twitter:image" content="https://dnote-asset.s3.amazonaws.com/images/logo-text-vertical.png" />
<meta name="og:image" content="https://dnote-asset.s3.amazonaws.com/images/logo-text-vertical.png" />
<meta name="og:title" content="{{ .Title }}" />
<meta name="og:description" content="{{ .Description }}" />`
