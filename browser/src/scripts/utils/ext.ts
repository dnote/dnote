// module ext provides a cross-browser interface to access extension APIs
// by using WebExtensions API if available, and using Chrome as a fallback.
let ext: any = {};

const apis = ['tabs', 'storage', 'runtime'];

for (let i = 0; i < apis.length; i++) {
  const api = apis[i];

  try {
    if (browser[api]) {
      ext[api] = browser[api];
    }
  } catch (e) {}

  try {
    if (chrome[api] && !ext[api]) {
      ext[api] = chrome[api];

      // Standardize the signature to conform to WebExtensions API
      if (api === 'tabs') {
        const fn = ext[api].create;

        // Promisify chrome.tabs.create
        ext[api].create = function(obj) {
          return new Promise(resolve => {
            fn(obj, function(tab) {
              resolve(tab);
            });
          });
        };
      }
    }
  } catch (e) {}
}

export default ext;
