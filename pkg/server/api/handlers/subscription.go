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

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/api/operations"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentsource"
	"github.com/stripe/stripe-go/source"
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/webhook"
)

var proPlanID = "plan_EpgsEvY27pajfo"

func getOrCreateStripeCustomer(tx *gorm.DB, user database.User) (*stripe.Customer, error) {
	if user.StripeCustomerID != "" {
		c, err := customer.Get(user.StripeCustomerID, nil)
		if err != nil {
			return nil, errors.Wrap(err, "getting customer")
		}

		return c, nil
	}

	var account database.Account
	if err := tx.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		return nil, errors.Wrap(err, "finding account")
	}

	customerParams := &stripe.CustomerParams{
		Email: &account.Email.String,
	}
	c, err := customer.New(customerParams)
	if err != nil {
		return nil, errors.Wrap(err, "creating customer")
	}

	user.StripeCustomerID = c.ID
	if err := tx.Save(&user).Error; err != nil {
		return nil, errors.Wrap(err, "updating user")
	}

	return c, nil
}

func addCustomerSource(customerID, sourceID string) (*stripe.PaymentSource, error) {
	params := &stripe.CustomerSourceParams{
		Customer: stripe.String(customerID),
		Source: &stripe.SourceParams{
			Token: stripe.String(sourceID),
		},
	}

	src, err := paymentsource.New(params)
	if err != nil {
		return nil, errors.Wrap(err, "creating source for customer")
	}

	return src, nil
}

func removeCustomerSource(customerID, sourceID string) (*stripe.Source, error) {
	params := &stripe.SourceObjectDetachParams{
		Customer: stripe.String(customerID),
	}
	s, err := source.Detach(sourceID, params)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func createCustomerSubscription(customerID, planID string) (*stripe.Subscription, error) {
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(planID),
			},
		},
	}

	s, err := sub.New(subParams)
	if err != nil {
		return nil, errors.Wrap(err, "creating subscription for customer")
	}

	return s, nil
}

type createSubPayload struct {
	Source  stripe.Source `json:"source"`
	Country string        `json:"country"`
}

// createSub creates a subscription for a the current user
func (a *App) createSub(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var payload createSubPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		handleError(w, "decoding params", err, http.StatusBadRequest)
		return
	}

	db := database.DBConn
	tx := db.Begin()

	if err := tx.Model(&user).
		Update(map[string]interface{}{
			"cloud":           true,
			"billing_country": payload.Country,
		}).Error; err != nil {
		tx.Rollback()
		handleError(w, "updating user", err, http.StatusInternalServerError)
		return
	}

	customer, err := getOrCreateStripeCustomer(tx, user)
	if err != nil {
		tx.Rollback()
		handleError(w, "getting customer", err, http.StatusInternalServerError)
		return
	}

	if _, err = addCustomerSource(customer.ID, payload.Source.ID); err != nil {
		tx.Rollback()
		handleError(w, "attaching source", err, http.StatusInternalServerError)
		return
	}

	if _, err := createCustomerSubscription(customer.ID, proPlanID); err != nil {
		tx.Rollback()
		handleError(w, "creating subscription", err, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		handleError(w, "committing a subscription transaction", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type updateSubPayload struct {
	StripeSubcriptionID string       `json:"stripe_subscription_id"`
	Op                  string       `json:"op"`
	Body                *interface{} `json:"body"`
}

var (
	updateSubOpCancel     = "cancel"
	updateSubOpReactivate = "reactivate"
)

var validUpdateSubOp = []string{
	updateSubOpCancel,
	updateSubOpReactivate,
}

func validateUpdateSubPayload(p updateSubPayload) error {
	var isOpValid bool

	for _, op := range validUpdateSubOp {
		if p.Op == op {
			isOpValid = true
			break
		}
	}

	if !isOpValid {
		return errors.Errorf("Invalid operation %s", p.Op)
	}

	if p.StripeSubcriptionID == "" {
		return errors.New("stripe_subscription_id is required")
	}

	return nil
}

func (a *App) updateSub(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}
	if user.StripeCustomerID == "" {
		handleError(w, "Customer does not exist", nil, http.StatusForbidden)
		return
	}

	var payload updateSubPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		handleError(w, "decoding params", err, http.StatusBadRequest)
		return
	}
	if err := validateUpdateSubPayload(payload); err != nil {
		handleError(w, "invalid payload", err, http.StatusBadRequest)
		return
	}

	var err error
	if payload.Op == updateSubOpCancel {
		err = operations.CancelSub(payload.StripeSubcriptionID, user)
	} else if payload.Op == updateSubOpReactivate {
		err = operations.ReactivateSub(payload.StripeSubcriptionID, user)
	}

	if err != nil {
		var statusCode int
		if err == operations.ErrSubscriptionActive {
			statusCode = http.StatusBadRequest
		} else {
			statusCode = http.StatusInternalServerError
		}

		handleError(w, fmt.Sprintf("during operation %s", payload.Op), err, statusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSubResponseItem represents a subscription item in the response for get subscription
type GetSubResponseItem struct {
	PlanID    string `json:"plan_id"`
	ProductID string `json:"product_id"`
}

// GetSubResponse is a response for getSub
type GetSubResponse struct {
	SubscriptionID     string                    `json:"id"`
	Items              []GetSubResponseItem      `json:"items"`
	CurrentPeriodStart int64                     `json:"current_period_start"`
	CurrentPeriodEnd   int64                     `json:"current_period_end"`
	Status             stripe.SubscriptionStatus `json:"status"`
	CancelAtPeriodEnd  bool                      `json:"cancel_at_period_end"`
}

func respondWithEmptySub(w http.ResponseWriter) {
	emptyGetSubResponse := GetSubResponse{
		Items: []GetSubResponseItem{},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyGetSubResponse); err != nil {
		handleError(w, "encoding response", err, http.StatusInternalServerError)
		return
	}
}

func (a *App) getSub(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}
	if user.StripeCustomerID == "" {
		respondWithEmptySub(w)
		return
	}

	listParams := &stripe.SubscriptionListParams{}
	listParams.Filters.AddFilter("customer", "", user.StripeCustomerID)
	listParams.Filters.AddFilter("status", "", "active")
	i := sub.List(listParams)

	if !i.Next() {
		if err := i.Err(); err != nil {
			handleError(w, "fetching subscription", err, http.StatusInternalServerError)
			return
		}

		// If no active subscription exists, respond with an empty subscription
		respondWithEmptySub(w)
		return
	}

	s := i.Subscription()

	resp := GetSubResponse{
		SubscriptionID:     s.ID,
		CurrentPeriodStart: s.CurrentPeriodStart,
		CurrentPeriodEnd:   s.CurrentPeriodEnd,
		Status:             s.Status,
		CancelAtPeriodEnd:  s.CancelAtPeriodEnd,
	}

	for _, item := range s.Items.Data {
		i := GetSubResponseItem{
			PlanID:    item.Plan.ID,
			ProductID: item.Plan.Product.ID,
		}
		resp.Items = append(resp.Items, i)
	}

	respondJSON(w, resp)
}

// GetStripeSourceResponse is a response for getStripeToken
type GetStripeSourceResponse struct {
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	ExpMonth uint8  `json:"exp_month"`
	ExpYear  uint16 `json:"exp_year"`
}

func respondWithEmptyStripeToken(w http.ResponseWriter) {
	var resp GetStripeSourceResponse

	respondJSON(w, resp)
}

// getStripeCard retrieves card information from stripe and returns a stripe.Card
// It handles legacy 'card' resource which have 'card_' prefixes, as well as the
// more up-to-date 'source' resources which have 'src_' prefixes.
func getStripeCard(stripeCustomerID, sourceID string) (*stripe.Card, error) {
	if strings.HasPrefix(sourceID, "card_") {
		params := &stripe.CardParams{
			Customer: stripe.String(stripeCustomerID),
		}
		cd, err := card.Get(sourceID, params)
		if err != nil {
			return nil, errors.Wrap(err, "fetching card")
		}

		return cd, nil
	} else if strings.HasPrefix(sourceID, "src_") {
		src, err := source.Get(sourceID, nil)
		if err != nil {
			return nil, errors.Wrap(err, "fetching source")
		}

		brand, ok := src.TypeData["brand"].(string)
		if !ok {
			return nil, errors.New("casting brand")
		}
		last4, ok := src.TypeData["last4"].(string)
		if !ok {
			return nil, errors.New("casting last4")
		}
		expMonth, ok := src.TypeData["exp_month"].(float64)
		if !ok {
			return nil, errors.New("casting exp_month")
		}
		expYear, ok := src.TypeData["exp_year"].(float64)
		if !ok {
			return nil, errors.New("casting exp_year")
		}

		cd := &stripe.Card{
			Brand:    stripe.CardBrand(brand),
			Last4:    last4,
			ExpMonth: uint8(expMonth),
			ExpYear:  uint16(expYear),
		}

		return cd, nil
	}

	return nil, errors.Errorf("malformed sourceID %s", sourceID)
}

type updateStripeSourcePayload struct {
	Source  stripe.Source `json:"source"`
	Country string        `json:"country"`
}

func validateUpdateStripeSourcePayload(p updateStripeSourcePayload) error {
	if p.Source.ID == "" {
		return errors.New("empty source id")
	}
	if p.Country == "" {
		return errors.New("empty country")
	}

	return nil
}

func (a *App) updateStripeSource(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var payload updateStripeSourcePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		handleError(w, "decoding params", err, http.StatusBadRequest)
		return
	}
	if err := validateUpdateStripeSourcePayload(payload); err != nil {
		http.Error(w, errors.Wrap(err, "validating payload").Error(), http.StatusBadRequest)
		return
	}

	db := database.DBConn
	tx := db.Begin()

	if err := tx.Model(&user).
		Update(map[string]interface{}{
			"billing_country": payload.Country,
		}).Error; err != nil {
		tx.Rollback()
		handleError(w, "updating user", err, http.StatusInternalServerError)
		return
	}

	c, err := customer.Get(user.StripeCustomerID, nil)
	if err != nil {
		tx.Rollback()
		handleError(w, "retriving customer", err, http.StatusInternalServerError)
		return
	}

	if _, err := removeCustomerSource(user.StripeCustomerID, c.DefaultSource.ID); err != nil {
		tx.Rollback()
		handleError(w, "removing source", err, http.StatusInternalServerError)
		return
	}

	if _, err := addCustomerSource(user.StripeCustomerID, payload.Source.ID); err != nil {
		tx.Rollback()
		handleError(w, "attaching source", err, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		handleError(w, "committing transaction", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) getStripeSource(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}
	if user.StripeCustomerID == "" {
		respondWithEmptyStripeToken(w)
		return
	}

	c, err := customer.Get(user.StripeCustomerID, nil)
	if err != nil {
		handleError(w, "fetching stripe customer", err, http.StatusInternalServerError)
		return
	}

	if c.DefaultSource == nil {
		respondWithEmptyStripeToken(w)
		return
	}

	cd, err := getStripeCard(user.StripeCustomerID, c.DefaultSource.ID)
	if err != nil {
		handleError(w, "fetching stripe source", err, http.StatusInternalServerError)
		return
	}

	resp := GetStripeSourceResponse{
		Brand:    string(cd.Brand),
		Last4:    cd.Last4,
		ExpMonth: cd.ExpMonth,
		ExpYear:  cd.ExpYear,
	}

	respondJSON(w, resp)
}

func (a *App) stripeWebhook(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, "reading body", err, http.StatusServiceUnavailable)
		return
	}

	webhookSecret := os.Getenv("StripeWebhookSecret")
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		handleError(w, "verifying stripe webhook signature", err, http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "customer.subscription.deleted":
		{
			var subscription stripe.Subscription
			if json.Unmarshal(event.Data.Raw, &subscription); err != nil {
				handleError(w, "unmarshaling payload", err, http.StatusBadRequest)
				return
			}

			operations.MarkUnsubscribed(subscription.Customer.ID)
		}
	default:
		{
			msg := fmt.Sprintf("Unsupported webhook event type %s", event.Type)
			handleError(w, msg, err, http.StatusBadRequest)
			return
		}
	}

	// Return a response to acknowledge receipt of the event
	w.WriteHeader(http.StatusOK)
}
