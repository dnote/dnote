import qs from 'qs';

function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }
  return response.text().then((body) => {
    const error = new Error(body);
    error.response = response;

    throw error;
  });
}

function parseJSON(response) {
  if (response.headers.get('Content-Type') === 'application/json') {
    return response.json();
  }

  return Promise.resolve();
}

function request(url, options) {
  return fetch(url, options)
    .then(checkStatus)
    .then(parseJSON);
}

export function post(url, data, options = {}) {
  return request(url, {
    method: 'POST',
    body: JSON.stringify(data),
    ...options,
  });
}

export function get(url, options = {}) {
  let endpoint = url;

  if (options.params) {
    endpoint = `${endpoint}?${qs.stringify(options.params)}`;
  }

  return request(endpoint, {
    method: 'GET',
    ...options,
  });
}
