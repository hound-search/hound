// TODO(knorton): Use something to bundle this more intelligently and get this
// out of the global scope.

var lib = {

    ExpandVars: function(template, values) {
      for (var name in values) {
        template = template.replace('{' + name + '}', values[name]);
      }
      return template;
    },

    UrlToRepo: function(repo, path, line, rev) {
        var url = repo.url.replace(/\.git$/, ''),
            pattern = repo['url-pattern'],
            filename = path.substring(path.lastIndexOf('/') + 1),
            anchor = line ? lib.ExpandVars(pattern.anchor, { line : line, filename : filename, repo: repo }) : '';

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

        // Check for bare git+ssh URIs (e.g., user@hostname:path
        var bare_ssh = /^([^@]+)@([^:]+):(.*)/.exec(url);

        if (ssh_protocol) {
          url = '//' + ssh_protocol[3] + '/' + ssh_protocol[5];
        } else if (bare_ssh) {
          url = '//' + bare_ssh[2] + '/' + bare_ssh[4];
        }

        // I'm sure there is a nicer React/jsx way to do this:
        return lib.ExpandVars(pattern['base-url'], {
          url : url,
          path: path,
          rev: rev,
          anchor: anchor
        });
    }
};
