import React from 'react';
import createReactClass from 'create-react-class';

export var RepoOption = createReactClass({
    render: function() {
        return (
            <option value={this.props.value} selected={this.props.selected}>{this.props.value}</option>
        )
    }
});
