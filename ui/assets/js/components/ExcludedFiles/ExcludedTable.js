import React from 'react';
import createReactClass from 'create-react-class';
import { ExcludedRow } from './ExcludedRow';

export var ExcludedTable = createReactClass({
    render: function() {
        var _this = this;
        if (this.props.searching) {
            return (<div id="no-result"><img src="images/busy.gif" /><div>Searching...</div></div>);
        }

        var rows = [];
        this.props.files.forEach(function(file, index) {
            rows.push(<ExcludedRow key={"exclude-row-" + index} file={file} repo={_this.props.repo} />);
        });

        return (
            <table>
                <thead>
                <tr>
                    <th>Filename</th>
                    <th>Reason</th>
                </tr>
                </thead>
                <tbody className="list">{rows}</tbody>
            </table>
        );
    }
});
