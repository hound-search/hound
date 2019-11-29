import React from 'react';

export const RepoButton = (props) => {
    const { repo, currentRepo, onRepoClick } = props;
    return (
        <button onClick={ onRepoClick.bind(this, repo) } className={ `repo-button${ repo === currentRepo ? ' selected' : '' }` }>
            { repo }
        </button>
    );
};
