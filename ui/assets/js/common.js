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
        anchor = line ? ExpandVars(pattern.anchor, { line : line, filename : filename, repo : repo }) : '';

    // Determine if the URL passed is a GitHub wiki
    var wikiUrl = /\.wiki$/.exec(url);
    if (wikiUrl) {
        url = url.replace(/\.wiki/, '/wiki')
        path = path.replace(/\.md$/, '')
        anchor = '' // wikis do not support direct line linking
    }

    // Check for ssh:// and hg:// protocol URLs
    // match the protocol, optionally a basic auth indicator, a
    // hostname, optionally a port, and then a path
    var ssh_protocol = /^(git|hg|ssh):\/\/([^@\/]+@)?([^:\/]+)(:[0-9]+)?\/(.*)/.exec(url);

    //// Begin EasyPost edit:  support phab links
    {
        //
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
            if (sshParts[3]) {
                port = sshParts[3]
            }
            url = hostname + port + '/' + project + '/' + repoName;
        }
    }
    //// End EasyPost edit

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

//// Begin EasyPost edit:  support phab links
function toPhabURL(parts) {
    // https://phab.easypo.net/source/easy_post/browse/master/app/helpers/sessions_helper.rb
    const rev = parts['rev'],
          path = parts['path'],
          project = parts['project'],
          repo = parts['repo'],
          anchor = parts['anchor'];
    const linenum = anchor ? anchor.replace(/^#L/, '') : undefined;
    const branch = 'master';
    // Caveat:  Assumes "master" branch is primary/exists.
    const result = ('https://phab.easypo.net/' + project
        + '/' + repo + '/browse/' + branch + '/' + path
        + ';' + rev);
    if (linenum) {
        return result + '$' + linenum;
    }
    return result;
}
//// End EasyPost edit

export function UrlToRepo(repo, path, line, rev) {
    var urlParts = UrlParts(repo, path, line, rev),
        pattern = repo['url-pattern']

    //// Begin EasyPost edit:  support phab links
    const hostname = urlParts['hostname'] || '';
    if (hostname.endsWith('phab.easypo.net')) {
        return toPhabURL(urlParts);
    }
    //// End EasyPost edit

    // I'm sure there is a nicer React/jsx way to do this:
    return ExpandVars(pattern['base-url'], urlParts);
}
