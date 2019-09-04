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

import { apiClient } from '../libs/http';
import { EmailPrefData } from '../store/auth';

export function updateUser({ name }) {
  const payload = { name };

  return apiClient.patch('/account/profile', payload);
}

interface UpdateEmailParams {
  email: string;
}

export function updateEmail({ email }: UpdateEmailParams) {
  const payload = {
    email
  };

  return apiClient.patch('/account/email', payload);
}

interface UpdatePasswordParams {
  oldPassword: string;
  newPassword: string;
}

export function updatePassword({
  oldPassword,
  newPassword
}: UpdatePasswordParams) {
  const payload = {
    old_password: oldPassword,
    new_password: newPassword
  };

  return apiClient.patch('/account/password', payload);
}

interface RegisterParams {
  email: string;
  password: string;
}

export function register(params: RegisterParams) {
  const payload = {
    email: params.email,
    password: params.password
  };

  return apiClient.post('/v3/register', payload);
}

interface SigninParams {
  email: string;
  password: string;
}

export function signin(params: SigninParams) {
  const payload = {
    email: params.email,
    password: params.password
  };

  return apiClient.post('/v3/signin', payload);
}

export function signout() {
  return apiClient.post('/v3/signout');
}

export function sendResetPasswordEmail({ email }) {
  const payload = { email };

  return apiClient.post('/reset-token', payload);
}

export function sendEmailVerificationEmail() {
  return apiClient.post('/verification-token');
}

export function verifyEmail({ token }) {
  const payload = { token };

  return apiClient.patch('/verify-email', payload);
}

export function updateEmailPreference({ token, digestFrequency }) {
  const payload = { digest_weekly: digestFrequency === 'weekly' };

  let endpoint = '/account/email-preference';
  if (token) {
    endpoint = `${endpoint}?token=${token}`;
  }
  return apiClient.patch(endpoint, payload);
}

interface GetEmailPreferenceParams {
  // if not logged in, users can optionally make an authenticated request using a token
  token?: string;
}

export function getEmailPreference({
  token
}: GetEmailPreferenceParams): Promise<EmailPrefData> {
  let endpoint = '/account/email-preference';
  if (token) {
    endpoint = `${endpoint}?token=${token}`;
  }

  return apiClient.get<EmailPrefData>(endpoint);
}

export function getMe() {
  return apiClient.get('/me').then(res => {
    return res.user;
  });
}

interface ResetPasswordParams {
  token: string;
  password: string;
}

export function resetPassword({ token, password }: ResetPasswordParams) {
  const payload = { token, password };

  return apiClient.patch('/reset-password', payload);
}

// classic
export function classicPresignin({ email }) {
  return apiClient.get(`/classic/presignin?email=${email}`);
}

interface classicPresigninPayload {
  key: string;
  expiresAt: number;
  cipherKeyEnc: string;
}

export function classicSignin({
  email,
  authKey
}): Promise<classicPresigninPayload> {
  const payload = { email, auth_key: authKey };

  return apiClient.post<any>('/classic/signin', payload).then(resp => {
    return {
      key: resp.key,
      expiresAt: resp.expires_at,
      cipherKeyEnc: resp.cipher_key_enc
    };
  });
}

interface classicSetPasswordPayload {
  password: string;
}

export function classicSetPassword({ password }: classicSetPasswordPayload) {
  const payload = {
    password
  };

  return apiClient.patch<any>('/classic/set-password', payload);
}

export function classicCompleteMigrate() {
  return apiClient.patch('/classic/migrate', '');
}
