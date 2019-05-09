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

// proPlanIds are ids of plans that are for 'Dnote Pro'.
const proPlanIds = ['prod_BUCQYMoPGOcXLa'];

// getPlanLabel returns a label for the plan in the given subscription
export function getPlanLabel(subscription) {
  if (!subscription || subscription.items.length === 0) {
    return 'Free';
  }

  const item = subscription.items[0];

  if (proPlanIds.indexOf(item.product_id) > -1) {
    return 'Dnote Pro';
  }

  return 'Free';
}
