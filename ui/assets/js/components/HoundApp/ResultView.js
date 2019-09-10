import React from 'react';
import createReactClass from 'create-react-class';
import { Model } from '../../helpers/Model';
import { FilesView } from './FilesView';

export var ResultView = createReactClass({
    componentWillMount: function() {
        var _this = this;
        Model.willSearch.tap(function(model, params) {
            _this.setState({
                results: null,
                query: params.q
            });
        });
    },
    getInitialState: function() {
        return { results: null };
    },
    render: function() {
        if (this.state.error) {
            return (
                <div id="no-result" className="error">
                    <strong>ERROR:</strong>{this.state.error}
                </div>
            );
        }

        if (this.state.results !== null && this.state.results.length === 0) {
            // TODO(knorton): We need something better here. :-(
            return (
                <div id="no-result">&ldquo;Nothing for you, Dawg.&rdquo;<div>0 results</div></div>
            );
        }

        if (this.state.results === null && this.state.query) {
            return (
                <div id="no-result"><img src="images/busy.gif" /><div>Searching...</div></div>
            );
        }

        var regexp = this.state.regexp,
            results = this.state.results || [];
        var repos = results.map(function(result, index) {
            return (
                <div className="repo" key={"results-view-" + index}>
                    <div className="title">
                        <span className="mega-octicon octicon-repo"></span>
                        <span className="name">{Model.NameForRepo(result.Repo)}</span>
                    </div>
                    <FilesView matches={result.Matches}
                               rev={result.Rev}
                               repo={result.Repo}
                               regexp={regexp}
                               totalMatches={result.FilesWithMatch} />
                </div>
            );
        });
        return (
            <div id="result">{repos}</div>
        );
    }
});
