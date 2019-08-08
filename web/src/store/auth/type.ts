import { RemoteData } from '../types';

export interface UserData {
  uuid: string;
  email: string;
  emailVerified: boolean;
  pro: boolean;
}

// TODO: type
export type EmailPrefData = any;
export type SubscriptionData = any;
export type SourceData = any;

export type UserState = RemoteData<UserData>;
export type EmailPrefState = RemoteData<EmailPrefData>;
export type SubscriptionState = RemoteData<SubscriptionData>;
export type SourceState = RemoteData<SourceData>;

export interface AuthState {
  user: UserState;
  emailPreference: EmailPrefState;
  subscription: SubscriptionState;
  source: SourceState;
}

export const START_FETCHING_USER = 'auth/START_FETCHING_USER';
export const RECEIVE_USER = 'auth/RECEIVE_USER';
export const RECEIVE_USER_ERROR = 'auth/RECEIVE_USER_ERROR';
export const RECEIVE_EMAIL_PREFERENCE = 'auth/RECEIVE_EMAIL_PREFERENCE';
export const RECEIVE_EMAIL_PREFERENCE_ERROR =
  'auth/RECEIVE_EMAIL_PREFERENCE_ERROR';
export const START_FETCHING_SUBSCRIPTION = 'auth/START_FETCHING_SUBSCRIPTION';
export const RECEIVE_SUBSCRIPTION = 'auth/RECEIVE_SUBSCRIPTION';
export const CLEAR_SUBSCRIPTION = 'auth/CLEAR_SUBSCRIPTION';
export const RECEIVE_SUBSCRIPTION_ERROR = 'auth/RECEIVE_SUBSCRIPTION_ERROR';
export const START_FETCHING_SOURCE = 'auth/START_FETCHING_SOURCE';
export const RECEIVE_SOURCE = 'auth/RECEIVE_SOURCE';
export const CLEAR_SOURCE = 'auth/CLEAR_SOURCE';
export const RECEIVE_SOURCE_ERROR = 'auth/RECEIVE_SOURCE_ERROR';

export interface StartFetchingUserAction {
  type: typeof START_FETCHING_USER;
}
export interface receiveUserAction {
  type: typeof RECEIVE_USER;
  data: any;
}
export interface receiveUserErrorAction {
  type: typeof RECEIVE_USER_ERROR;
  data: any;
}
export interface ReceiveEmailPreferenceAction {
  type: typeof RECEIVE_EMAIL_PREFERENCE;
  data: any;
}
export interface receiveEmailPreferenceErrorAction {
  type: typeof RECEIVE_EMAIL_PREFERENCE_ERROR;
  data: any;
}
export interface startFetchingSubscriptionAction {
  type: typeof START_FETCHING_SUBSCRIPTION;
}
export interface receiveSubscriptionAction {
  type: typeof RECEIVE_SUBSCRIPTION;
  data: any;
}
export interface clearSubscriptionAction {
  type: typeof CLEAR_SUBSCRIPTION;
}
export interface receiveSubscriptionErrorAction {
  type: typeof RECEIVE_SUBSCRIPTION_ERROR;
  data: any;
}
export interface startFetchingSourceAction {
  type: typeof START_FETCHING_SOURCE;
}
export interface receiveSourceAction {
  type: typeof RECEIVE_SOURCE;
  data: any;
}
export interface clearSourceAction {
  type: typeof CLEAR_SOURCE;
}
export interface receiveSourceErrorAction {
  type: typeof RECEIVE_SOURCE_ERROR;
  data: any;
}

export type AuthActionType =
  | StartFetchingUserAction
  | receiveUserAction
  | receiveUserErrorAction
  | ReceiveEmailPreferenceAction
  | receiveEmailPreferenceErrorAction
  | startFetchingSubscriptionAction
  | receiveSubscriptionAction
  | clearSubscriptionAction
  | receiveSubscriptionErrorAction
  | startFetchingSourceAction
  | receiveSourceAction
  | clearSourceAction
  | receiveSourceErrorAction;
