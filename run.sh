xgo -go go-1.17.x -targets=windows/amd64 -ldflags '-X main.apiEndpoint=https://api.getdnote.com -X main.versionTag=foo' -buildmode=exe -tags fts5 -pkg pkg/cli -x -v .

