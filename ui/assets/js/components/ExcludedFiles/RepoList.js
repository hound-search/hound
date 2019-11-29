import React from 'react';
import { RepoButton } from './RepoButton';

export const RepoList = (props) => {

    const { repo, repos, onRepoClick } = props;

    const reposBlock = repos.map((item, index) => (
        <RepoButton
            key={`repobutton-${index}`}
            repo={ item }
            onRepoClick={ onRepoClick }
            currentRepo={ repo }
        />
    ));

    return (
        <div id="repolist">
            { reposBlock }
        </div>
    );
};
