import React from 'react';
import { ExcludedRow } from './ExcludedRow';

export const ExcludedTable = (props) => {

    const { files, searching, repo } = props;

    if (searching) {
        return (
            <div id="no-result">
                <img src="images/busy.gif" /><div>Searching...</div>
            </div>
        );
    }

    const rows = files.map((file, index) => (
        <ExcludedRow
            key={`exclude-row-${index}`}
            file={ file }
            repo={ repo }
        />
    ));

    return (
        <table>
            <thead>
            <tr>
                <th>Filename</th>
                <th>Reason</th>
            </tr>
            </thead>
            <tbody className="list">{ rows }</tbody>
        </table>
    );
};
