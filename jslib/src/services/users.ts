/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import { EmailPrefData, UserData } from '../operations/types';

export interface UpdateProfileParams {
  email: string;
  password: string;
}

export interface UpdatePasswordParams {
  oldPassword: string;
  newPassword: string;
}

export interface RegisterParams {
  email: string;
  password: string;
}

export interface SigninParams {
  email: string;
  password: string;
}

export interface SigninResponse {
  key: string;
  expires_at: number;
}

export interface GetEmailPreferenceParams {
  // if not logged in, users can optionally make an authenticated request using a token
  token?: string;
}

export interface GetEmailPreferenceResponse {
  inactive_reminder: boolean;
  product_update: boolean;
}

export interface UpdateEmailPreferenceParams {
  token?: string;
  inactiveReminder?: boolean;
  productUpdate?: boolean;
}

export interface ResetPasswordParams {
  token: string;
  password: string;
}

export interface GetMeResponse {
  user: {
    uuid: string;
    email: string;
    email_verified: boolean;
    pro: boolean;
  };
}

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    updateUser: ({ name }) => {
      const payload = { name };

      return client.patch('/account/profile', payload);
    },

    updateProfile: ({ email, password }: UpdateProfileParams) => {
      const payload = {
        email,
        password
      };

      return client.patch('/account/profile', payload);
    },

    updatePassword: ({ oldPassword, newPassword }: UpdatePasswordParams) => {
      const payload = {
        old_password: oldPassword,
        new_password: newPassword
      };

      return client.patch('/account/password', payload);
    },

    register: (params: RegisterParams) => {
      const payload = {
        email: params.email,
        password: params.password
      };

      return client.post('/v3/register', payload);
    },

    signin: (params: SigninParams) => {
      const payload = {
        email: params.email,
        password: params.password
      };

      return client.post<SigninResponse>('/v3/signin', payload).then(resp => {
        return {
          key: resp.key,
          expiresAt: resp.expires_at
        };
      });
    },

    signout: () => {
      return client.post('/v3/signout');
    },

    sendResetPasswordEmail: ({ email }) => {
      const payload = { email };

      return client.post('/reset-token', payload);
    },

    sendEmailVerificationEmail: () => {
      return client.post('/verification-token');
    },

    verifyEmail: ({ token }) => {
      const payload = { token };

      return client.patch('/verify-email', payload);
    },

    updateEmailPreference: ({
      token,
      inactiveReminder,
      productUpdate
    }: UpdateEmailPreferenceParams): Promise<EmailPrefData> => {
      const payload: any = {};

      if (inactiveReminder !== undefined) {
        payload.inactive_reminder = inactiveReminder;
      }
      if (productUpdate !== undefined) {
        payload.product_update = productUpdate;
      }

      let endpoint = '/account/email-preference';
      if (token) {
        endpoint = `${endpoint}?token=${token}`;
      }

      return client
        .patch<GetEmailPreferenceResponse>(endpoint, payload)
        .then(res => {
          return {
            inactiveReminder: res.inactive_reminder,
            productUpdate: res.product_update
          };
        });
    },

    getEmailPreference: ({
      token
    }: GetEmailPreferenceParams): Promise<EmailPrefData> => {
      let endpoint = '/account/email-preference';
      if (token) {
        endpoint = `${endpoint}?token=${token}`;
      }

      return client.get<GetEmailPreferenceResponse>(endpoint).then(res => {
        return {
          inactiveReminder: res.inactive_reminder,
          productUpdate: res.product_update
        };
      });
    },

    getMe: (): Promise<UserData> => {
      return client.get<GetMeResponse>('/me').then(res => {
        const { user } = res;

        return {
          uuid: user.uuid,
          email: user.email,
          emailVerified: user.email_verified,
          pro: user.pro
        };
      });
    },

    resetPassword: ({ token, password }: ResetPasswordParams) => {
      const payload = { token, password };

      return client.patch('/reset-password', payload);
    }
  };
}
