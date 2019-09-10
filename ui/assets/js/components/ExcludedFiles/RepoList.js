import React from 'react';
import createReactClass from 'create-react-class';
import { RepoButton } from './RepoButton';

export var RepoList = createReactClass({
    render: function() {
        var repos = [],
            _this = this;
        this.props.repos.forEach(function(repo, index){
            repos.push(<RepoButton key={"repo-list-" + index} repo={repo} onRepoClick={_this.props.onRepoClick} currentRepo={_this.props.repo} />);
        });

        return (
            <div id="repolist">
                {repos}
            </div>
        );
    }
});
