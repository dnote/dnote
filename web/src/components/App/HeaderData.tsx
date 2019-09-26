/* Copyright (C) 2019 Monomax Software Pty Ltd
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

import React from 'react';
import Helmet from 'react-helmet';

import config from 'web/libs/config';

const title = 'Dnote - A Simple Personal Knowledge Base';
const description =
  'Dnote is a personal knowledge base with an automated spaced repetition. Give a home to your knowledge and let your knowledge flow.';

export default () => {
  return (
    <Helmet defaultTitle={title} titleTemplate="%s | Dnote">
      <meta name="description" content={description} />
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:site" content="@dnoteio" />
      <meta name="twitter:title" content={title} />
      <meta name="twitter:description" content={description} />
      <meta
        name="twitter:image"
        content={`${config.cdnUrl}/brands/logo-text-horizontal-large.png`}
      />
      <meta name="twitter:image:alt" content="A logo of dnote" />
      <meta name="og:title" content="Dnote" />
      <meta name="og:site_name" content="Dnote" />
      <meta
        name="og:image"
        content={`${config.cdnUrl}/brands/logo-text-horizontal-large.png`}
      />
      <meta name="og:description" content={description} />
    </Helmet>
  );
};
