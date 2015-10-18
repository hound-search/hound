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
            anchor = line ? lib.ExpandVars(pattern.anchor, { line : line, filename : filename }) : '';

        // Hacky solution to fix _some more_ of the 404's when using SSH style URLs.
        // This works for both github style URLs (git@github.com:username/Foo.git) and
        // bitbucket style URLs (ssh://hg@bitbucket.org/username/Foo).

        // Regex explained: Match either `git` or `hg` followed by an `@`.
        // Next, slurp up the hostname by reading until either a `:` or `/` is found.
        // Finally, grab all remaining characters.
        var sshParts = /(git|hg)@(.*?)(:|\/)(.*)/.exec(url);
        if (sshParts) {
          url = '//' + sshParts[2] + '/' + sshParts[4];
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


try {
  if (localStorage["useDarkTheme"] === "true") {
    document.getElementById('theme').href="css/hound-dark.css";
  }
} catch (e) {

}
