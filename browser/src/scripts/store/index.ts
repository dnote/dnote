import { combineReducers, createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';

import location from './location/reducers';
import settings from './settings/reducers';
import books from './books/reducers';
import composer from './composer/reducers';

const rootReducer = combineReducers({
  location,
  settings,
  books,
  composer
});

// configuruStore returns a new store that contains the appliation state
export default function configureStore(initialState) {
  const typedWindow = window as any;

  const composeEnhancers =
    typedWindow.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

  return createStore(
    rootReducer,
    initialState,
    composeEnhancers(applyMiddleware(createLogger, thunkMiddleware))
  );
}
