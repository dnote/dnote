import initServices from 'jslib/services';
import config from './config';

const services = initServices({
  baseUrl: config.apiEndpoint,
  pathPrefix: ''
});

export default services;
