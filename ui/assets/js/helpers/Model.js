import reqwest from 'reqwest';
import merge from 'merge-anything';
import { Signal } from './Signal';
import { UrlToRepo } from './common';
import { ParamsFromUrl } from '../utils';

/**
 * The data model for the UI is responsible for conducting searches and managing
 * all results.
 */
export var Model = {
    // raised when a search begins
    willSearch: new Signal(),

    // raised when a search completes
    didSearch: new Signal(),

    willLoadMore: new Signal(),

    didLoadMore: new Signal(),

    didError: new Signal(),

    didLoadRepos : new Signal(),

    ValidRepos: function(repos) {
        var all = this.repos,
            seen = {};
        return repos.filter(function(repo) {
            var valid = all[repo] && !seen[repo];
            seen[repo] = true;
            return valid;
        });
    },

    RepoCount: function() {
        return Object.keys(this.repos).length;
    },

    Load: function() {
        var _this = this;
        var next = function() {
            var params = ParamsFromUrl();
            _this.didLoadRepos.raise(_this, _this.repos);

            if (params.q !== '') {
                _this.Search(params);
            }
        };

        if (typeof ModelData != 'undefined') {
            var data = JSON.parse(ModelData),
                repos = {};
            for (var name in data) {
                repos[name] = data[name];
            }
            this.repos = repos;
            next();
            return;
        }

        reqwest({
            url: 'api/v1/repos',
            type: 'json',
            success: function(data) {
                _this.repos = data;
                next();
            },
            error: function(xhr, status, err) {
                // TODO(knorton): Fix these
                console.error(err);
            }
        });
    },

    Search: function(params) {
        this.willSearch.raise(this, params);
        var _this = this,
            startedAt = Date.now();

        params = merge({
            stats: 'fosho',
            repos: '*',
            rng: ':20',
        }, params);

        if (params.repos === '') {
            params.repos = '*';
        }

        _this.params = params;

        // An empty query is basically useless, so rather than
        // sending it to the server and having the server do work
        // to produce an error, we simply return empty results
        // immediately in the client.
        if (params.q == '') {
            _this.results = [];
            _this.resultsByRepo = {};
            _this.didSearch.raise(_this, _this.Results);
            return;
        }

        reqwest({
            url: 'api/v1/search',
            data: params,
            type: 'json',
            success: function(data) {
                if (data.Error) {
                    _this.didError.raise(_this, data.Error);
                    return;
                }

                var matches = data.Results,
                    stats = data.Stats,
                    results = [];
                for (var repo in matches) {
                    if (!matches[repo]) {
                        continue;
                    }

                    var res = matches[repo];
                    results.push({
                        Repo: repo,
                        Rev: res.Revision,
                        Matches: res.Matches,
                        FilesWithMatch: res.FilesWithMatch,
                    });
                }

                results.sort(function(a, b) {
                    return b.Matches.length - a.Matches.length || a.Repo.localeCompare(b.Repo);
                });

                var byRepo = {};
                results.forEach(function(res) {
                    byRepo[res.Repo] = res;
                });

                _this.results = results;
                _this.resultsByRepo = byRepo;
                _this.stats = {
                    Server: stats.Duration,
                    Total: Date.now() - startedAt,
                    Files: stats.FilesOpened
                };

                _this.didSearch.raise(_this, _this.results, _this.stats);
            },
            error: function(xhr, status, err) {
                _this.didError.raise(this, "The server broke down");
            }
        });
    },

    LoadMore: function(repo) {
        var _this = this,
            results = this.resultsByRepo[repo],
            numLoaded = results.Matches.length,
            numNeeded = results.FilesWithMatch - numLoaded,
            numToLoad = Math.min(2000, numNeeded),
            endAt = numNeeded == numToLoad ? '' : '' + numToLoad;

        _this.willLoadMore.raise(this, repo, numLoaded, numNeeded, numToLoad);

        var params = merge(this.params, {
            rng: numLoaded+':'+endAt,
            repos: repo
        });

        reqwest({
            url: 'api/v1/search',
            data: params,
            type: 'json',
            success: function(data) {
                if (data.Error) {
                    _this.didError.raise(_this, data.Error);
                    return;
                }

                var result = data.Results[repo];
                results.Matches = results.Matches.concat(result.Matches);
                _this.didLoadMore.raise(_this, repo, _this.results);
            },
            error: function(xhr, status, err) {
                _this.didError.raise(this, "The server broke down");
            }
        });
    },

    NameForRepo: function(repo) {
        var info = this.repos[repo];
        if (!info) {
            return repo;
        }

        var url = info.url,
            ax = url.lastIndexOf('/');
        if (ax  < 0) {
            return repo;
        }

        var name = url.substring(ax + 1).replace(/\.git$/, '');

        var bx = url.lastIndexOf('/', ax - 1);
        if (bx < 0) {
            return name;
        }

        return url.substring(bx + 1, ax) + ' / ' + name;
    },

    UrlToRepo: function(repo, path, line, rev) {
        return UrlToRepo(this.repos[repo], path, line, rev);
    }

};
