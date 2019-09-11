import React, { useEffect, useState, useRef } from 'react';
import { FormatNumber, ParamValueToBool } from '../../utils';
import { Model } from '../../helpers/Model';
import Select from 'react-select';

export const SearchBar = (props) => {

    const { query, ignoreCase, files, repos, allRepos, stats, onSearchRequested } = props;
    const [ showAdvanced, setShowAdvanced] = useState(false);
    const [ search, setSearch ] = useState({ query, ignoreCase, files, repos } );
    const queryInput = useRef(null);
    const fileInput = useRef(null);

    const hasAdvancedValues = () => (
        ( search.files && search.files.trim() !== '' ) ||
        ( search.ignoreCase && search.ignoreCase.trim() === 'fosho' ) ||
        ( search.repos && search.repos.length > 0 )
    );

    useEffect( () => {
        setSearch({ query, files, repos, ignoreCase });
    }, [query, files, repos, ignoreCase, allRepos, stats]);

    const repoOptions = allRepos.map(rname => ({
        value: rname,
        label: rname
    }));

    const selectedRepos = repoOptions.filter(o => search.repos.indexOf(o.value) >= 0);

    const showAdvancedCallback = () => {
        setShowAdvanced(true);
        if (search.query.trim() !== '') {
            fileInput.current.focus();
        }
    };

    const hideAdvancedCallback = () => {
        setShowAdvanced(false);
        if (queryInput.current) {
            queryInput.current.focus();
        }
    };

    const elementChanged = (prop, checkbox, evt) => {
        setSearch({
            ...search,
            [prop]: checkbox
                ? evt.currentTarget.checked && 'fosho' || 'nope'
                : evt.currentTarget.value
        });
    };

    const submitQuery = () => {
        if (search.query.trim() !== '') {
            onSearchRequested({
                q: search.query,
                i: search.ignoreCase,
                files: search.files,
                repos: Model.ValidRepos(search.repos) === Model.RepoCount()
                    ? ''
                    : search.repos.join(',')
            });
        }
    };

    const queryGotKeydown = (event) => {
        switch (event.keyCode) {
            case 40:
                showAdvancedCallback();
                fileInput.current.focus();
                break;
            case 38:
                hideAdvancedCallback();
                break;
            case 13:
                submitQuery();
                break;
        }
    };

    const queryGotFocus = () => {
        if ( !hasAdvancedValues() ) {
            hideAdvancedCallback();
        }
    };

    const filesGotKeydown = (event) => {
        switch (event.keyCode) {
            case 38:
                // if advanced is empty, close it up.
                if (search.files.trim() === '') {
                    hideAdvancedCallback();
                }
                queryInput.current.focus();
                break;
            case 13:
                submitQuery();
                break;
        }
    };

    const repoSelected = (selected) => {
        setSearch({
            ...search,
            repos: selected
                ? selected.map(item => item.value)
                : []
        });
    };

    const statsView = stats
        ? (
            <div className="stats">
                <div className="stats-left">
                    <a href="excluded_files.html"
                       className="link-gray">
                        Excluded Files
                    </a>
                </div>
                <div className="stats-right">
                    <div className="val">{ FormatNumber(stats.Total) }ms total</div> /
                    <div className="val">{ FormatNumber(stats.Server) }ms server</div> /
                    <div className="val">{ stats.Files } files</div>
                </div>
            </div>
        )
        : '';

    return (
        <div id="input">
            <div id="ina">
                <input
                    ref={ queryInput }
                    type="text"
                    placeholder="Search by Regexp"
                    autoComplete="off"
                    autoFocus
                    value={ search.query }
                    onFocus={ queryGotFocus }
                    onChange={ elementChanged.bind(this, "query", false) }
                    onKeyDown={ queryGotKeydown }
                />
                <div className="button-add-on">
                    <button id="dodat" onClick={ submitQuery }></button>
                </div>
            </div>

            <div id="inb" className={ showAdvanced ? 'opened' : 'closed' }>
                <div id="adv">
                    <span className="octicon octicon-chevron-up hide-adv" onClick={ hideAdvancedCallback }></span>
                    <div className="field">
                        <label htmlFor="files">File Path</label>
                        <div className="field-input">
                            <input
                                ref={ fileInput }
                                type="text"
                                placeholder="regexp"
                                value={ search.files }
                                onChange={ elementChanged.bind(this, "files", false) }
                                onKeyDown={ filesGotKeydown }
                            />
                        </div>
                    </div>
                    <div className="field">
                        <label htmlFor="ignore-case">Ignore Case</label>
                        <div className="field-input">
                            <input type="checkbox" onChange={ elementChanged.bind(this, "ignoreCase", true) } checked={ ParamValueToBool(search.ignoreCase) } />
                        </div>
                    </div>
                    <div className="field-repo-select">
                        <label className="multiselect_label" htmlFor="repos">Select Repo</label>
                        <div className="field-input">
                            <Select
                                options={ repoOptions }
                                onChange={ repoSelected }
                                value={ selectedRepos }
                                isMulti
                                closeMenuOnSelect={ false }
                            />
                        </div>
                    </div>
                </div>
                <div className="ban" onClick={ showAdvancedCallback }>
                    <em>Advanced:</em> ignore case, filter by path, stuff like that.
                </div>
            </div>
            { statsView }
        </div>
    );
};
