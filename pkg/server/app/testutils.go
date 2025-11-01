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
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/assets"
	"github.com/dnote/dnote/pkg/server/testutils"
)

// NewTest returns an app for a testing environment
func NewTest() App {
	return App{
		Clock:               clock.NewMock(),
		EmailBackend:        &testutils.MockEmailbackendImplementation{},
		HTTP500Page:         assets.MustGetHTTP500ErrorPage(),
		WebURL:              "http://127.0.0.0.1",
		Port:                "3000",
		DisableRegistration: false,
		DBPath:              "",
		AssetBaseURL:        "",
	}
}
