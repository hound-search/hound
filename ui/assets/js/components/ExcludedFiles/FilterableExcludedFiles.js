import React from 'react';
import createReactClass from 'create-react-class';
import reqwest from 'reqwest';
import { RepoList } from './RepoList';
import { ExcludedTable } from './ExcludedTable'

export var FilterableExcludedFiles = createReactClass({
    getInitialState: function() {
        var _this = this;
        reqwest({
            url: 'api/v1/repos',
            type: 'json',
            success: function(data) {
                _this.setState({ repos: data });
            },
            error: function(xhr, status, err) {
                // TODO(knorton): Fix these
                console.error(err);
            }
        });

        return {
            files: [],
            repos: [],
            repo: null,
        };
    },

    onRepoClick: function(repo) {
        var _this = this;
        _this.setState({
            searching: true,
            repo: this.state.repos[repo],
        });
        reqwest({
            url: 'api/v1/excludes',
            data: {repo: repo},
            type: 'json',
            success: function(data) {
                _this.setState({ files: data, searching: false });
            },
            error: function(xhr, status, err) {
                // TODO(knorton): Fix these
                console.error(err);
            }
        });
    },

    render: function() {
        return (
            <div id="excluded_container">
                <a href="/">Home</a>
                <h1>Excluded Files</h1>

                <div id="excluded_files" className="table-container">
                    <RepoList repos={Object.keys(this.state.repos)} onRepoClick={this.onRepoClick} repo={this.state.repo} />
                    <ExcludedTable files={this.state.files} searching={this.state.searching} repo={this.state.repo} />
                </div>
            </div>
        );
    }
});
