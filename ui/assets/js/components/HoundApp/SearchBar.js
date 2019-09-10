import React from 'react';
import createReactClass from 'create-react-class';
import { Model } from '../../helpers/Model';
import { css, FormatNumber, ParamValueToBool } from '../../utils';
import Select from 'react-select';

export var SearchBar = createReactClass({
    componentWillMount: function() {
        var _this = this;
        Model.didLoadRepos.tap(function(model, repos) {
            _this.setState({
                allRepos: Object.keys(repos),
                repos: _this.props.repos
            });
        });
    },

    componentDidMount: function() {
        var q = this.refs.q;

        // TODO(knorton): Can't set this in jsx
        q.setAttribute('autocomplete', 'off');

        this.setParams(this.props);

        if (this.hasAdvancedValues()) {
            this.showAdvanced();
        }

        q.focus();
    },
    getInitialState: function() {
        return {
            state: null,
            allRepos: [],
            repos: []
        };
    },
    queryGotKeydown: function(event) {
        switch (event.keyCode) {
            case 40:
                // this will cause advanced to expand if it is not expanded.
                this.refs.files.focus();
                break;
            case 38:
                this.hideAdvanced();
                break;
            case 13:
                this.submitQuery();
                break;
        }
    },
    queryGotFocus: function(event) {
        if (!this.hasAdvancedValues()) {
            this.hideAdvanced();
        }
    },
    filesGotKeydown: function(event) {
        switch (event.keyCode) {
            case 38:
                // if advanced is empty, close it up.
                if (this.refs.files.value.trim() === '') {
                    this.hideAdvanced();
                }
                this.refs.q.focus();
                break;
            case 13:
                this.submitQuery();
                break;
        }
    },
    filesGotFocus: function(event) {
        this.showAdvanced();
    },
    submitQuery: function() {
        this.props.onSearchRequested(this.getParams());
    },
    getRegExp : function() {
        return new RegExp(
            this.refs.q.value.trim(),
            this.refs.icase.checked ? 'ig' : 'g');
    },
    getParams: function() {
        // selecting all repos is the same as not selecting any, so normalize the url
        // to have none.
        var repos = Model.ValidRepos(this.state.repos);

        if (repos.length == Model.RepoCount()) {
            repos = [];
        }

        return {
            q : this.refs.q.value.trim(),
            files : this.refs.files.value.trim(),
            repos : repos.join(','),
            i: this.refs.icase.checked ? 'fosho' : 'nope'
        };
    },
    setParams: function(params) {
        var q = this.refs.q,
            i = this.refs.icase,
            files = this.refs.files;

        q.value = params.q;
        i.checked = ParamValueToBool(params.i);
        files.value = params.files;
    },
    hasAdvancedValues: function() {
        return this.refs.files.value.trim() !== '' || this.refs.icase.checked || this.state.repos.length > 0;
    },
    showAdvanced: function() {
        var adv = this.refs.adv,
            ban = this.refs.ban,
            q = this.refs.q,
            files = this.refs.files;

        css(adv, 'height', 'auto');
        css(adv, 'padding', '10px 0');
        css(adv, 'overflow', 'visible');

        css(ban, 'max-height', '0');
        css(ban, 'opacity', '0');

        if (q.value.trim() !== '') {
            files.focus();
        }
    },
    hideAdvanced: function() {
        var adv = this.refs.adv,
            ban = this.refs.ban,
            q = this.refs.q;

        css(adv, 'height', '0');
        css(adv, 'padding', '0');
        css(adv, 'overflow', 'hidden');

        css(ban, 'max-height', '100px');
        css(ban, 'opacity', '1');

        q.focus();
    },
    repoSelected: function (selected) {
        this.setState({
            repos: selected
                ? selected.map(function (item) {
                    return item.value
                })
                : []
        });
    },
    render: function() {

        var _this = this;

        var repoOptions = this.state.allRepos.map(function (repoName) {
            return {
                value: repoName,
                label: repoName
            };
        });

        var selectedRepos = repoOptions.filter(function (option) {
            return _this.state.repos.indexOf(option.value) >= 0;
        });

        var stats = this.state.stats;
        var statsView = '';
        if (stats) {
            statsView = (
                <div className="stats">
                    <div className="stats-left">
                        <a href="excluded_files.html"
                           className="link-gray">
                            Excluded Files
                        </a>
                    </div>
                    <div className="stats-right">
                        <div className="val">{FormatNumber(stats.Total)}ms total</div> /
                        <div className="val">{FormatNumber(stats.Server)}ms server</div> /
                        <div className="val">{stats.Files} files</div>
                    </div>
                </div>
            );
        }

        return (
            <div id="input">
                <div id="ina">
                    <input id="q"
                           type="text"
                           placeholder="Search by Regexp"
                           ref="q"
                           autoComplete="off"
                           onKeyDown={this.queryGotKeydown}
                           onFocus={this.queryGotFocus}/>
                    <div className="button-add-on">
                        <button id="dodat" onClick={this.submitQuery}></button>
                    </div>
                </div>

                <div id="inb">
                    <div id="adv" ref="adv">
                        <span className="octicon octicon-chevron-up hide-adv" onClick={this.hideAdvanced}></span>
                        <div className="field">
                            <label htmlFor="files">File Path</label>
                            <div className="field-input">
                                <input type="text"
                                       id="files"
                                       placeholder="regexp"
                                       ref="files"
                                       onKeyDown={this.filesGotKeydown}
                                       onFocus={this.filesGotFocus} />
                            </div>
                        </div>
                        <div className="field">
                            <label htmlFor="ignore-case">Ignore Case</label>
                            <div className="field-input">
                                <input id="ignore-case" type="checkbox" ref="icase" />
                            </div>
                        </div>
                        <div className="field-repo-select">
                            <label className="multiselect_label" htmlFor="repos">Select Repo</label>
                            <div className="field-input">
                                <Select
                                    options={ repoOptions }
                                    onChange={ this.repoSelected }
                                    value={ selectedRepos }
                                    isMulti
                                    closeMenuOnSelect={false}
                                />
                            </div>
                        </div>
                    </div>
                    <div className="ban" ref="ban" onClick={this.showAdvanced}>
                        <em>Advanced:</em> ignore case, filter by path, stuff like that.
                    </div>
                </div>
                {statsView}
            </div>
        );
    }
});
