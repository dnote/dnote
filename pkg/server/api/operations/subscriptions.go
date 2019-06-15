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

package operations

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

// ErrSubscriptionActive is an error indicating that the subscription is active
// and therefore cannot be reactivated
var ErrSubscriptionActive = errors.New("The subscription is currently active")

// CancelSub cancels the subscription of the given user
func CancelSub(subscriptionID string, user database.User) error {
	updateParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := sub.Update(subscriptionID, updateParams)
	if err != nil {
		return errors.Wrap(err, "updating subscription on Stripe")
	}

	return nil
}

// ReactivateSub reactivates the subscription of the given user
func ReactivateSub(subscriptionID string, user database.User) error {
	s, err := sub.Get(subscriptionID, nil)
	if err != nil {
		return errors.Wrap(err, "fetching subscription")
	}

	if !s.CancelAtPeriodEnd {
		return ErrSubscriptionActive
	}

	updateParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	if _, err := sub.Update(subscriptionID, updateParams); err != nil {
		return errors.Wrap(err, "updating subscription on Stripe")
	}

	return nil
}

// MarkUnsubscribed marks the user unsubscribed
func MarkUnsubscribed(stripeCustomerID string) error {
	db := database.DBConn

	var user database.User
	if err := db.Where("stripe_customer_id = ?", stripeCustomerID).First(&user).Error; err != nil {
		return errors.Wrap(err, "finding user")
	}

	if err := db.Model(&user).Update("cloud", false).Error; err != nil {
		return errors.Wrap(err, "updating user")
	}

	return nil
}
