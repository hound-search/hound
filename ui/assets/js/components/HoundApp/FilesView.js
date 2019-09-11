import React from 'react';
import { Model } from '../../helpers/Model';
import { File } from './File';

export const FilesView = (props) => {

    const { matches, rev, repo, regexp, totalMatches } = props;

    const onLoadMore = () => Model.LoadMore(repo);

    const files = matches.map((match, index) => (
        <File
            key={`${repo}-file-${index}`}
            repo={ repo }
            rev={ rev }
            match={ match }
            regexp={ regexp }
        />
    ));

    const more = (matches.length < totalMatches)
        ? (
            <button className="moar" onClick={ onLoadMore }>
                Load all {totalMatches} matches in { Model.NameForRepo(repo) }
            </button>
        )
        : '';

    return (
        <div className="files">
            { files }
            { more }
        </div>
    );

};
