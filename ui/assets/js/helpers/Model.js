import reqwest from 'reqwest';
import { Signal } from './Signal';
import { UrlToRepo } from './common';
import { ParamsFromUrl } from '../utils';

/**
 * The data model for the UI is responsible for conducting searches and managing
 * all results.
 */
export const Model = {
    // raised when a search begins
    willSearch: new Signal(),

    // raised when a search completes
    didSearch: new Signal(),

    willLoadMore: new Signal(),

    didLoadMore: new Signal(),

    didError: new Signal(),

    didLoadRepos : new Signal(),

    ValidRepos (repos) {
        const all = this.repos;
        const seen = {};
        return repos.filter((repo) => {
            const valid = all[repo] && !seen[repo];
            seen[repo] = true;
            return valid;
        });
    },

    RepoCount () {
        return Object.keys(this.repos).length;
    },

    Load () {

        const _this = this;

        const next = () => {
            const params = ParamsFromUrl();
            this.didLoadRepos.raise(this, this.repos);
            if (params.q !== '') {
                this.Search(params);
            }
        };

        reqwest({
            url: 'api/v1/repos',
            type: 'json',
            success (data) {
                _this.repos = data;
                next();
            },
            error (xhr, status, err) {
                // TODO(knorton): Fix these
                console.error(err);
            }
        });
    },

    Search (params) {

        const _this = this;
        const startedAt = Date.now();

        this.willSearch.raise(this, params);

        params = {
            stats: 'fosho',
            repos: '*',
            rng: ':20',
            ...params
        };

        if (params.repos === '') {
            params.repos = '*';
        }

        this.params = params;

        // An empty query is basically useless, so rather than
        // sending it to the server and having the server do work
        // to produce an error, we simply return empty results
        // immediately in the client.
        if (params.q === '') {
            this.results = [];
            this.resultsByRepo = {};
            this.didSearch.raise(this, this.Results);
            return;
        }

        reqwest({
            url: 'api/v1/search',
            data: params,
            type: 'json',
            success (data) {
                if (data.Error) {
                    _this.didError.raise(_this, data.Error);
                    return;
                }

                const matches = data.Results;
                const stats = data.Stats;
                const results = [];

                for (let repo in matches) {
                    if (!matches[repo]) {
                        continue;
                    }

                    const res = matches[repo];
                    results.push({
                        Repo: repo,
                        Rev: res.Revision,
                        Matches: res.Matches,
                        FilesWithMatch: res.FilesWithMatch,
                    });
                }

                results.sort((a, b) => b.Matches.length - a.Matches.length || a.Repo.localeCompare(b.Repo));

                const byRepo = results.reduce((obj, res) => (obj[res.Repo] = res, obj), {});

                _this.results = results;
                _this.resultsByRepo = byRepo;
                _this.stats = {
                    Server: stats.Duration,
                    Total: Date.now() - startedAt,
                    Files: stats.FilesOpened
                };

                _this.didSearch.raise(_this, _this.results, _this.stats);
            },
            error (xhr, status, err) {
                _this.didError.raise(this, "The server broke down");
            }
        });
    },

    LoadMore (repo) {
        const _this = this;
        const results = this.resultsByRepo[repo];
        const numLoaded = results.Matches.length;
        const numNeeded = results.FilesWithMatch - numLoaded;
        const numToLoad = Math.min(2000, numNeeded);
        const endAt = numNeeded == numToLoad ? '' : '' + numToLoad;

        this.willLoadMore.raise(this, repo, numLoaded, numNeeded, numToLoad);

        const params = {...this.params,
            rng: numLoaded+':'+endAt,
            repos: repo
        };

        reqwest({
            url: 'api/v1/search',
            data: params,
            type: 'json',
            success (data) {
                if (data.Error) {
                    _this.didError.raise(_this, data.Error);
                    return;
                }

                const result = data.Results[repo];
                results.Matches = results.Matches.concat(result.Matches);
                _this.didLoadMore.raise(_this, repo, _this.results);
            },
            error (xhr, status, err) {
                _this.didError.raise(this, "The server broke down");
            }
        });
    },

    NameForRepo (repo) {
        const info = this.repos[repo];
        if (!info) {
            return repo;
        }

        const url = info.url;
        const ax = url.lastIndexOf('/');
        if (ax  < 0) {
            return repo;
        }

        const name = url.substring(ax + 1).replace(/\.git$/, '');

        const bx = url.lastIndexOf('/', ax - 1);

        if (bx < 0) {
            return name;
        }

        return url.substring(bx + 1, ax) + ' / ' + name;
    },

    UrlToRepo (repo, path, line, rev) {
        return UrlToRepo(this.repos[repo], path, line, rev);
    }

};
