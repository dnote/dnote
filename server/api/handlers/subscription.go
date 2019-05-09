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
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dnote/dnote/server/api/helpers"
	"github.com/dnote/dnote/server/api/logger"
	"github.com/dnote/dnote/server/api/operations"
	"github.com/dnote/dnote/server/database"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/webhook"
)

type stripeToken struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

var planID = "plan_EpgsEvY27pajfo"

func init() {
}

// createSub creates a subscription for a the current user
func (a *App) createSub(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}
	if user.StripeCustomerID != "" {
		http.Error(w, "Customer already exists", http.StatusForbidden)
		return
	}

	var tok stripeToken
	if err := json.NewDecoder(r.Body).Decode(&tok); err != nil {
		http.Error(w, errors.Wrap(err, "decoding params").Error(), http.StatusInternalServerError)
		return
	}

	customerParams := &stripe.CustomerParams{
		Plan:  &planID,
		Email: &tok.Email,
	}
	err := customerParams.SetSource(tok.ID)
	if err != nil {
		http.Error(w, errors.Wrap(err, "setting source").Error(), http.StatusInternalServerError)
		return
	}

	//TODO: if customer exists, update not create
	c, err := customer.New(customerParams)
	if err != nil {
		http.Error(w, errors.Wrap(err, "creating customer").Error(), http.StatusInternalServerError)
		return
	}

	user.StripeCustomerID = c.ID
	user.Cloud = true
	if err := db.Save(&user).Error; err != nil {
		http.Error(w, errors.Wrap(err, "updating user").Error(), http.StatusInternalServerError)
		return
	}
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
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}
	if user.StripeCustomerID == "" {
		http.Error(w, "Customer does not exist", http.StatusForbidden)
		return
	}

	var payload updateSubPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, errors.Wrap(err, "decoding params").Error(), http.StatusInternalServerError)
		return
	}

	if err := validateUpdateSubPayload(payload); err != nil {
		http.Error(w, errors.Wrap(err, "invalid payload").Error(), http.StatusBadRequest)
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

		http.Error(w, errors.Wrapf(err, "during operation %s", payload.Op).Error(), statusCode)
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
	emptyGetSubREsponse := GetSubResponse{
		Items: []GetSubResponseItem{},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyGetSubREsponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getSub(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
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
			http.Error(w, errors.Wrap(err, "fetching subscription").Error(), http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getStripeSource(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}
	if user.StripeCustomerID == "" {
		respondWithEmptyStripeToken(w)
		return
	}

	c, err := customer.Get(user.StripeCustomerID, nil)
	if err != nil {
		http.Error(w, errors.Wrap(err, "fetching stripe customer").Error(), http.StatusInternalServerError)
		return
	}
	if c.DefaultSource == nil {
		respondWithEmptyStripeToken(w)
		return
	}

	params := &stripe.CardParams{
		Customer: stripe.String(user.StripeCustomerID),
	}
	cd, err := card.Get(c.DefaultSource.ID, params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "fetching stripe card").Error(), http.StatusInternalServerError)
		return
	}

	resp := GetStripeSourceResponse{
		Brand:    string(cd.Brand),
		Last4:    cd.Last4,
		ExpMonth: cd.ExpMonth,
		ExpYear:  cd.ExpYear,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) stripeWebhook(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Err("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	webhookSecret := os.Getenv("StripeWebhookSecret")
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		logger.Err("Error verifying the signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "customer.subscription.deleted":
		{
			var subscription stripe.Subscription
			if json.Unmarshal(event.Data.Raw, &subscription); err != nil {
				logger.Err(errors.Wrap(err, "unmarshaling").Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			operations.MarkUnsubscribed(subscription.Customer.ID)
		}
	default:
		{
			logger.Err("Unsupported webhook event type %s", event.Type)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Return a response to acknowledge receipt of the event
	w.WriteHeader(http.StatusOK)
}
