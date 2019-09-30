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

import { getHttpClient, HttpClientConfig } from '../helpers/http';

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    createSubscription: ({ source, country }) => {
      const payload = {
        source,
        country
      };

      return client.post('/subscriptions', payload);
    },

    getSubscription: () => {
      return client.get('/subscriptions');
    },

    cancelSubscription: ({ subscriptionId }) => {
      const data = {
        op: 'cancel',
        stripe_subscription_id: subscriptionId
      };

      return client.patch('/subscriptions', data);
    },

    reactivateSubscription: ({ subscriptionId }) => {
      const data = {
        op: 'reactivate',
        stripe_subscription_id: subscriptionId
      };

      return client.patch('/subscriptions', data);
    },

    getSource: () => {
      return client.get('/stripe_source');
    },

    updateSource: ({ source, country }) => {
      const payload = {
        source,
        country
      };

      return client.patch('/stripe_source', payload);
    }
  };
}
