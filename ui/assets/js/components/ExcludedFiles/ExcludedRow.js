import React from 'react';
import createReactClass from 'create-react-class';
import { UrlToRepo } from '../../helpers/common';

export var ExcludedRow = createReactClass({
    render: function() {
        var url = UrlToRepo(this.props.repo, this.props.file.Filename, this.props.rev);
        return (
            <tr>
                <td className="name">
                    <a href={url}>{this.props.file.Filename}</a>
                </td>
                <td className="reason">{this.props.file.Reason}</td>
            </tr>
        );
    }
});
