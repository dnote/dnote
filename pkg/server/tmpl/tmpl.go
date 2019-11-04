package tmpl

var noteMetaTags = `
<meta name="description" content="{{ .Description }}" />
<meta name="twitter:card" content="summary" />
<meta name="twitter:title" content="{{ .Title }}" />
<meta name="twitter:description" content="{{ .Description }}" />
<meta name="twitter:image" content="https://dnote-asset.s3.amazonaws.com/images/logo-text-vertical.png" />
<meta name="og:image" content="https://dnote-asset.s3.amazonaws.com/images/logo-text-vertical.png" />
<meta name="og:title" content="{{ .Title }}" />
<meta name="og:description" content="{{ .Description }}" />`
