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
import { EmailPrefData, UserData } from '../operations/types';

export interface UpdateProfileParams {
  email: string;
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

export interface classicPresigninPayload {
  key: string;
  expiresAt: number;
  cipherKeyEnc: string;
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
    classic: boolean;
  };
}

export interface classicSetPasswordPayload {
  password: string;
}

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    updateUser: ({ name }) => {
      const payload = { name };

      return client.patch('/account/profile', payload);
    },

    updateProfile: ({ email }: UpdateProfileParams) => {
      const payload = {
        email
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

    updateEmailPreference: ({ token, digestFrequency }) => {
      const payload = { digest_weekly: digestFrequency === 'weekly' };

      let endpoint = '/account/email-preference';
      if (token) {
        endpoint = `${endpoint}?token=${token}`;
      }
      return client.patch(endpoint, payload);
    },

    getEmailPreference: ({
      token
    }: GetEmailPreferenceParams): Promise<EmailPrefData> => {
      let endpoint = '/account/email-preference';
      if (token) {
        endpoint = `${endpoint}?token=${token}`;
      }

      return client.get<EmailPrefData>(endpoint);
    },

    getMe: (): Promise<UserData> => {
      return client.get<GetMeResponse>('/me').then(res => {
        const { user } = res;

        return {
          uuid: user.uuid,
          email: user.email,
          emailVerified: user.email_verified,
          pro: user.pro,
          classic: user.classic
        };
      });
    },

    resetPassword: ({ token, password }: ResetPasswordParams) => {
      const payload = { token, password };

      return client.patch('/reset-password', payload);
    },

    // classic
    classicPresignin: ({ email }) => {
      return client.get(`/classic/presignin?email=${email}`);
    },

    classicSignin: ({ email, authKey }): Promise<classicPresigninPayload> => {
      const payload = { email, auth_key: authKey };

      return client.post<any>('/classic/signin', payload).then(resp => {
        return {
          key: resp.key,
          expiresAt: resp.expires_at,
          cipherKeyEnc: resp.cipher_key_enc
        };
      });
    },

    classicSetPassword: ({ password }: classicSetPasswordPayload) => {
      const payload = {
        password
      };

      return client.patch<any>('/classic/set-password', payload);
    },

    classicCompleteMigrate: () => {
      return client.patch('/classic/migrate', '');
    }
  };
}
