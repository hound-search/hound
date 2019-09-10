import React from 'react';
import createReactClass from 'create-react-class';
import { CoalesceMatches, ContentFor } from '../../utils';
import { Model } from '../../helpers/Model';

export var FilesView = createReactClass({
    onLoadMore: function(event) {
        Model.LoadMore(this.props.repo);
    },

    render: function() {
        var rev = this.props.rev,
            repo = this.props.repo,
            regexp = this.props.regexp,
            matches = this.props.matches,
            totalMatches = this.props.totalMatches;
        var files = matches.map(function(match, index) {
            var filename = match.Filename,
                blocks = CoalesceMatches(match.Matches);
            var matches = blocks.map(function(block, mindex) {
                var lines = block.map(function(line, lindex) {
                    var content = ContentFor(line, regexp);
                    return (
                        <div className="line" key={repo + "-" + lindex + "-" + mindex + "-" + index}>
                            <a href={Model.UrlToRepo(repo, filename, line.Number, rev)}
                               className="lnum"
                               target="_blank">{line.Number}</a>
                            <span className="lval" dangerouslySetInnerHTML={{__html:content}} />
                        </div>
                    );
                });

                return (
                    <div className="match" key={repo + "-lines-" + mindex + "-" + index}>{lines}</div>
                );
            });

            return (
                <div className="file" key={repo + "-file-" + index}>
                    <div className="title">
                        <a href={Model.UrlToRepo(repo, match.Filename, null, rev)}>
                            {match.Filename}
                        </a>
                    </div>
                    <div className="file-body">
                        {matches}
                    </div>
                </div>
            );
        });

        var more = '';
        if (matches.length < totalMatches) {
            more = (<button className="moar" onClick={this.onLoadMore}>Load all {totalMatches} matches in {Model.NameForRepo(repo)}</button>);
        }

        return (
            <div className="files">
                {files}
                {more}
            </div>
        );
    }
});
