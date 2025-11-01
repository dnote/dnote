/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestSendWelcomeEmail(t *testing.T) {
	emailBackend := testutils.MockEmailbackendImplementation{}
	a := NewTest()
	a.EmailBackend = &emailBackend
	a.BaseURL = "http://example.com"

	if err := a.SendWelcomeEmail("alice@example.com"); err != nil {
		t.Fatal(err, "failed to perform")
	}

	assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
	assert.Equal(t, emailBackend.Emails[0].From, "noreply@example.com", "email sender mismatch")
	assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")

}

func TestSendPasswordResetEmail(t *testing.T) {
	emailBackend := testutils.MockEmailbackendImplementation{}
	a := NewTest()
	a.EmailBackend = &emailBackend
	a.BaseURL = "http://example.com"

	if err := a.SendPasswordResetEmail("alice@example.com", "mockTokenValue"); err != nil {
		t.Fatal(err, "failed to perform")
	}

	assert.Equalf(t, len(emailBackend.Emails), 1, "email queue count mismatch")
	assert.Equal(t, emailBackend.Emails[0].From, "noreply@example.com", "email sender mismatch")
	assert.DeepEqual(t, emailBackend.Emails[0].To, []string{"alice@example.com"}, "email sender mismatch")

}

func TestGetSenderEmail(t *testing.T) {
	testCases := []struct {
		baseURL        string
		expectedSender string
	}{
		{
			baseURL:        "https://www.example.com",
			expectedSender: "noreply@example.com",
		},
		{
			baseURL:        "https://www.example2.com",
			expectedSender: "alice@example2.com",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("base url %s", tc.baseURL), func(t *testing.T) {
		})
	}
}
