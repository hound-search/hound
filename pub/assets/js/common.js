// TODO(knorton): Use something to bundle this more intelligently and get this
// out of the global scope.

var ExpandVars = function(template, values) {
  for (var name in values) {
    template = template.replace('{' + name + '}', values[name]);
  }
  return template;
};

var UrlToRepo = function(repo, path, line) {
    var url = repo.url.replace(/\.git$/, ''),
        pattern = repo['url-pattern'],
        anchor = line ? ExpandVars(pattern.anchor, { line : line }) : '';

    // Hacky solution to fix _some more_ of the 404's when using SSH style URLs
    var sshParts = /git@(.*):(.*)/i.exec(url);
    if (sshParts) {
      url = '//' + sshParts[1] + '/' + sshParts[2];
    }

    // I'm sure there is a nicer React/jsx way to do this:
    return ExpandVars(pattern['base-url'], {
      url : url,
      path: path,
      anchor: anchor
    });
};
