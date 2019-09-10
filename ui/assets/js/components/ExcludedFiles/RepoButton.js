import React from 'react';
import createReactClass from 'create-react-class';

export var RepoButton = createReactClass({
    handleClick: function(repoName) {
        this.props.onRepoClick(repoName);
    },
    render: function() {
        var className = 'repo-button';
        if (this.props.currentRepo === this.props.repo) {
            className += ' selected';
        }

        return (
            <button onClick={this.handleClick.bind(this, this.props.repo)} className={className}>
                {this.props.repo}
            </button>
        );
    }
});
