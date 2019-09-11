import React from 'react';
import { UrlToRepo } from '../../helpers/common';

export const ExcludedRow = (props) => {

    const { file: { Filename, Reason }, repo } = props;

    const url = UrlToRepo(repo, Filename);

    return (
        <tr>
            <td className="name">
                <a href={url}>{ Filename }</a>
            </td>
            <td className="reason">{ Reason }</td>
        </tr>
    );
};
