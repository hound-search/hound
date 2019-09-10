import React from 'react';
import createReactClass from 'create-react-class';
import { ParamsFromUrl } from '../../utils';
import { Model } from '../../helpers/Model';
import { SearchBar } from './SearchBar';
import { ResultView } from './ResultView';

export var App = createReactClass({
    componentWillMount: function() {
        var params = ParamsFromUrl(),
            repos = (params.repos == '') ? [] : params.repos.split(',');

        this.setState({
            q: params.q,
            i: params.i,
            files: params.files,
            repos: repos
        });

        var _this = this;
        Model.didLoadRepos.tap(function(model, repos) {
            // If all repos are selected, don't show any selected.
            if (model.ValidRepos(_this.state.repos).length == model.RepoCount()) {
                _this.setState({repos: []});
            }
        });

        Model.didSearch.tap(function(model, results, stats) {
            _this.refs.searchBar.setState({
                stats: stats,
                repos: repos,
            });

            _this.refs.resultView.setState({
                results: results,
                regexp: _this.refs.searchBar.getRegExp(),
                error: null
            });
        });

        Model.didLoadMore.tap(function(model, repo, results) {
            _this.refs.resultView.setState({
                results: results,
                regexp: _this.refs.searchBar.getRegExp(),
                error: null
            });
        });

        Model.didError.tap(function(model, error) {
            _this.refs.resultView.setState({
                results: null,
                error: error
            });
        });

        window.addEventListener('popstate', function(e) {
            var params = ParamsFromUrl();
            _this.refs.searchBar.setParams(params);
            Model.Search(params);
        });
    },
    onSearchRequested: function(params) {
        this.updateHistory(params);
        Model.Search(this.refs.searchBar.getParams());
    },
    updateHistory: function(params) {
        var path = location.pathname +
            '?q=' + encodeURIComponent(params.q) +
            '&i=' + encodeURIComponent(params.i) +
            '&files=' + encodeURIComponent(params.files) +
            '&repos=' + params.repos;
        history.pushState({path:path}, '', path);
    },
    render: function() {
        return (
            <div>
                <SearchBar ref="searchBar"
                           q={this.state.q}
                           i={this.state.i}
                           files={this.state.files}
                           repos={this.state.repos}
                           onSearchRequested={this.onSearchRequested} />
                <ResultView ref="resultView" q={this.state.q} />
            </div>
        );
    }
});
