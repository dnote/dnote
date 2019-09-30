import { NAVIGATE, NavigateAction } from './types';

export function navigate(path: string, state?): NavigateAction {
  return {
    type: NAVIGATE,
    data: { path, state }
  };
}
