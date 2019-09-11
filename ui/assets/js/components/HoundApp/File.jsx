import React from 'react';
import { CoalesceMatches, ContentFor } from '../../utils';
import { Model } from '../../helpers/Model';
import { Match } from './Match';

export const File = (props) => {

    const { repo, rev, match, regexp } = props;
    const filename = match.Filename;
    const blocks = CoalesceMatches(match.Matches);

    const matches = blocks.map((block, index) => (
        <Match
            key={`match-${repo}-${index}`}
            block={ block }
            repo={ repo }
            regexp={ regexp }
            rev={ rev }
            filename={ filename }
        />
    ));

    return (
        <div className="file">
            <div className="title">
                <a href={ Model.UrlToRepo(repo, filename, null, rev) }>
                    { filename }
                </a>
            </div>
            <div className="file-body">
                { matches }
            </div>
        </div>
    );

};
