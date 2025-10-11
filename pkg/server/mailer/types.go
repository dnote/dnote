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

package mailer

// EmailResetPasswordTmplData is a template data for reset password emails
type EmailResetPasswordTmplData struct {
	AccountEmail string
	Token        string
	WebURL       string
}

// EmailResetPasswordAlertTmplData is a template data for reset password emails
type EmailResetPasswordAlertTmplData struct {
	AccountEmail string
	WebURL       string
}

// WelcomeTmplData is a template data for welcome emails
type WelcomeTmplData struct {
	AccountEmail string
	WebURL       string
}
