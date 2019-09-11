import React, { useEffect, useState } from 'react';
import reqwest from 'reqwest';
import { RepoList } from './RepoList';
import { ExcludedTable } from './ExcludedTable'

export const FilterableExcludedFiles = () => {

    const [ files, setFiles ] = useState([]);
    const [ repos, setRepos ] = useState([]);
    const [ repo, setRepo ] = useState(null);
    const [ searching, setSearching ] = useState(false);

    useEffect(() => {

        reqwest({
            url: 'api/v1/repos',
            type: 'json',
            success (data) {
                setRepos(data);
            },
            error (xhr, status, err) {
                // TODO(knorton): Fix these
                console.error(err);
            }
        });

    }, []);

    const clickOnRepo = (repo) => {

        setSearching(true);
        setRepo(repos[repo]);

        reqwest({
            url: 'api/v1/excludes',
            data: { repo: repo },
            type: 'json',
            success (data) {
                setSearching(false);
                setFiles(data);
            },
            error (xhr, status, err) {
                // TODO(knorton): Fix these
                console.error(err);
            }
        });
    };

    return (
        <div id="excluded_container">
            <a href="/">Home</a>
            <h1>Excluded Files</h1>

            <div id="excluded_files" className="table-container">
                <RepoList repos={ Object.keys(repos) } onRepoClick={ clickOnRepo } repo={ repo } />
                <ExcludedTable files={ files } searching={ searching } repo={ repo } />
            </div>
        </div>
    );
};
