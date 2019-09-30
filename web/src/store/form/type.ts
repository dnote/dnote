export interface FormState {
  auth: {
    email: string;
  };
}

export const UPDATE_AUTH_EMAIL = 'form/UPDATE_AUTH_EMAIL';

export interface UpdateAuthEmail {
  type: typeof UPDATE_AUTH_EMAIL;
  data: {
    email: string;
  };
}

export type FormActionType = UpdateAuthEmail;
