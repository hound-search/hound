export function EscapeRegExp(regexp) {
    return regexp.replace(/[-[\]{}()*+!<=:?.\/\\^$|#\s,]/g, '\\$&');
}

export function ExpandVars(template, values) {
    for (var name in values) {
        template = template.replace('{' + name + '}', values[name]);
    }
    return template;
};

export function UrlParts(repo, path, line, rev) {
    var url = repo.url.replace(/\.git$/, ''),
        pattern = repo['url-pattern'],
        hostname = '',
        project = '',
        repoName = '',
        path = path || '',
        port = '',
        filename = path.substring(path.lastIndexOf('/') + 1),
        anchor = line ? ExpandVars(pattern.anchor, { line : line, filename : filename }) : '';

    // Determine if the URL passed is a GitHub wiki
    var wikiUrl = /\.wiki$/.exec(url);
    if (wikiUrl) {
        url = url.replace(/\.wiki/, '/wiki')
        path = path.replace(/\.md$/, '')
        anchor = '' // wikis do not support direct line linking
    }

    // Hacky solution to fix _some more_ of the 404's when using SSH style URLs.
    // This works for both github style URLs (git@github.com:username/Foo.git) and
    // bitbucket style URLs (ssh://hg@bitbucket.org/username/Foo).

    // Regex explained: Match either `git` or `hg` followed by an `@`.
    // Next, slurp up the hostname by reading until either a `:` or `/` is found.
    // If a port is specified, slurp that up too. Finally, grab the project and
    // repo names.
    var sshParts = /(git|hg)@(.*?)(:[0-9]+)?(:|\/)(.*)(\/)(.*)/.exec(url);
    if (sshParts) {
        hostname = '//' + sshParts[2]
        project = sshParts[5]
        repoName = sshParts[7]
        // Port is omitted in most cases. Bitbucket Server is special:
        // ssh://git@bitbucket.atlassian.com:7999/ATLASSIAN/jira.git
        if(sshParts[3]){
            port = sshParts[3]
        }
        url = hostname + port + '/' + project + '/' + repoName;
    }

    return {
        url : url,
        hostname: hostname,
        port: port,
        project: project,
        'repo': repoName,
        path: path,
        rev: rev,
        anchor: anchor
    };
}

export function UrlToRepo(repo, path, line, rev) {
    var urlParts = UrlParts(repo, path, line, rev),
        pattern = repo['url-pattern']

    // I'm sure there is a nicer React/jsx way to do this:
    return ExpandVars(pattern['base-url'], urlParts);
}
