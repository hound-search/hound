import React from 'react';
import { Model } from '../../helpers/Model';
import { FilesView } from './FilesView';

export const ResultView = (props) => {

    const { query, ignoreCase, results, error } = props;
    const regexp = new RegExp(query.trim(), ignoreCase.trim() === 'fosho' && 'ig' || 'g');
    const isLoading = results === null && query;
    const noResults = !!results && results.length === 0;

    if (error) {
        return (
            <div id="no-result" className="error">
                <strong>ERROR:</strong>{ error }
            </div>
        );
    }

    if (!isLoading && noResults) {
        // TODO(knorton): We need something better here. :-(
        return (
            <div id="no-result">
                &ldquo;Nothing for you, Dawg.&rdquo;<div>0 results</div>
            </div>
        );
    }

    const repos = results
        ? results.map((result, index) => (
            <div className="repo" key={`results-view-${index}`}>
                <div className="title">
                    <span className="mega-octicon octicon-repo"></span>
                    <span className="name">{ Model.NameForRepo(result.Repo) }</span>
                </div>
                <FilesView
                    matches={ result.Matches }
                    rev={ result.Rev }
                    repo={ result.Repo }
                    regexp={ regexp }
                    totalMatches={ result.FilesWithMatch }
                />
            </div>
        ))
        : '';

    return (
        <div id="result">
            <div id="no-result" className={ isLoading && 'loading' || 'hidden' }>
                <img src="images/busy.gif" /><div>Searching...</div>
            </div>
            { repos }
        </div>
    );
};
