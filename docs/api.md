API
---

## `/` serves the web app

## `/api/v1/repos` returns list of indexed repositories

## `/api/v1/search` performs searches

For example using cURL to search for `foobar` in Hound:

```console
https://hound.examplle.com/api/v1/search?repos=Hound&q=foobar
...
```

## `/api/v1/excludes` returns list of excluded files

## `/api/v1/update` and `/api/v1/github-webhook` hook for updating repositories

## `/healthz` health-check

As configured by `health-check-uri` in the configuration. Defaults to `healthz`.

## `/metrics` Prometheus metrics

Returns the default Prometheus Go-metrics and some custom values on number and
duration of search requests.

To get the number of searches as a five-minute running average:

```
rate(hound_api_search_all_duration_seconds_count[5m])
```
