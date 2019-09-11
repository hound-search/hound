import React from 'react';
import {Model} from '../../helpers/Model';

export const Line = (props) => {

    const { line, rev, repo, filename, content } = props;

    return (
        <div className="line">
            <a href={ Model.UrlToRepo(repo, filename, line.Number, rev) }
               className="lnum"
               target="_blank"
            >
                { line.Number }
            </a>
            <span className="lval" dangerouslySetInnerHTML={ {__html:content} } />
        </div>
    );

};
