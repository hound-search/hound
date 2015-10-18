var Signal = function() {
};

Signal.prototype = {
  listeners : [],

  tap: function(l) {
    // Make a copy of the listeners to avoid the all too common
    // subscribe-during-dispatch problem
    this.listeners = this.listeners.slice(0);
    this.listeners.push(l);
  },

  untap: function(l) {
    var ix = this.listeners.indexOf(l);
    if (ix == -1) {
      return;
    }

    // Make a copy of the listeners to avoid the all to common
    // unsubscribe-during-dispatch problem
    this.listeners = this.listeners.slice(0);
    this.listeners.splice(ix, 1);
  },

  raise: function() {
    var args = Array.prototype.slice.call(arguments, 0);
    this.listeners.forEach(function(l) {
      l.apply(this, args);
    });
  }
};

var css = function(el, n, v) {
  el.style.setProperty(n, v, '');
};

var FormatNumber = function(t) {
  var s = '' + (t|0),
      b = [];
  while (s.length > 0) {
    b.unshift(s.substring(s.length - 3, s.length));
    s = s.substring(0, s.length - 3);
  }
  return b.join(',');
};

var ParamsFromQueryString = function(qs, params) {
  params = params || {};

  if (!qs) {
    return params;
  }

  qs.substring(1).split('&').forEach(function(v) {
    var pair = v.split('=');
    if (pair.length != 2) {
      return;
    }
    if (pair[1].indexOf(',') >= 0) {
      params[decodeURIComponent(pair[0])] = pair[1].split(',');
    } else {
      params[decodeURIComponent(pair[0])] = decodeURIComponent(pair[1]);
    }
  });

  if (params["repos"] === '') {
    params["repos"] = '*';
  }

  return params;
};

var ParamsFromUrl = function(params) {
  params = params || {
    q: '',
    i: 'nope',
    files: '',
    repos: '*'
  };
  return ParamsFromQueryString(location.search, params);
};

var ParamValueToBool = function(v) {
  v = v.toLowerCase();
  return v == 'fosho' || v == 'true' || v == '1';
};

/**
 * The data model for the UI is responsible for conducting searches and managing
 * all results.
 */
var Model = {
  // raised when a search begins
  willSearch: new Signal(),

  // raised when a search completes
  didSearch: new Signal(),

  willLoadMore: new Signal(),

  didLoadMore: new Signal(),

  didError: new Signal(),

  didLoadRepos : new Signal(),

  Load: function() {
    var _this = this;
    var next = function() {
      var params = ParamsFromUrl();
      _this.didLoadRepos.raise(_this, _this.repos);

      if (params.q !== '') {
        _this.Search(params);
      }
    };

    if (typeof ModelData != 'undefined') {
      var data = JSON.parse(ModelData),
          repos = {};
      for (var name in data) {
        repos[name.toLowerCase()] = data[name];
      }
      this.repos = repos;
      next();
      return;
    }

    $.ajax({
      url: '/api/v1/repos',
      dataType: 'json',
      success: function(data) {
        _this.repos = data;
        next();
      },
      error: function(xhr, status, err) {
        // TODO(knorton): Fix these
        console.error(err);
      }
    });
  },

  Search: function(params) {
    this.willSearch.raise(this, params);
    var _this = this,
        startedAt = Date.now();

    params = $.extend({
      stats: 'fosho',
      repos: '*',
      rng: ':20',
    }, params);

    if (params.repos === '') {
      params.repos = '*';
    }

    _this.params = params;

    // An empty query is basically useless, so rather than
    // sending it to the server and having the server do work
    // to produce an error, we simply return empty results
    // immediately in the client.
    if (params.q == '') {
      _this.results = [];
      _this.resultsByRepo = {};
      _this.didSearch.raise(_this, _this.Results);
      return;
    }

    $.ajax({
      url: '/api/v1/search',
      data: params,
      type: 'GET',
      dataType: 'json',
      success: function(data) {
        if (data.Error) {
          _this.didError.raise(_this, data.Error);
          return;
        }

        var matches = data.Results,
            stats = data.Stats,
            results = [];
        for (var repo in matches) {
          if (!matches[repo]) {
            continue;
          }

          var res = matches[repo];
          results.push({
            Repo: repo,
            Rev: res.Revision,
            Matches: res.Matches,
            FilesWithMatch: res.FilesWithMatch,
          });
        }

        results.sort(function(a, b) {
          return b.Matches.length - a.Matches.length;
        });

        var byRepo = {};
        results.forEach(function(res) {
          byRepo[res.Repo] = res;
        });

        _this.results = results;
        _this.resultsByRepo = byRepo;
        _this.stats = {
          Server: stats.Duration,
          Total: Date.now() - startedAt,
          Files: stats.FilesOpened
        };

        _this.didSearch.raise(_this, _this.results, _this.stats);
      },
      error: function(xhr, status, err) {
        _this.didError.raise(this, "The server broke down");
      }
    });
  },

  LoadMore: function(repo) {
    var _this = this,
        results = this.resultsByRepo[repo],
        numLoaded = results.Matches.length,
        numNeeded = results.FilesWithMatch - numLoaded,
        numToLoad = Math.min(2000, numNeeded),
        endAt = numNeeded == numToLoad ? '' : '' + numToLoad;

    _this.willLoadMore.raise(this, repo, numLoaded, numNeeded, numToLoad);

    var params = $.extend(this.params, {
      rng: numLoaded+':'+endAt,
      repos: repo
    });

    $.ajax({
      url: '/api/v1/search',
      data: params,
      type: 'GET',
      dataType: 'json',
      success: function(data) {
        if (data.Error) {
          _this.didError.raise(_this, data.Error);
          return;
        }

        var result = data.Results[repo];
        results.Matches = results.Matches.concat(result.Matches);
        _this.didLoadMore.raise(_this, repo, _this.results);
      },
      error: function(xhr, status, err) {
        _this.didError.raise(this, "The server broke down");
      }
    });
  },

  NameForRepo: function(repo) {
    var info = this.repos[repo];
    if (!info) {
      return repo;
    }

    var url = info.url,
        ax = url.lastIndexOf('/');
    if (ax  < 0) {
      return repo;
    }

    var name = url.substring(ax + 1).replace(/\.git$/, '');

    var bx = url.lastIndexOf('/', ax - 1);
    if (bx < 0) {
      return name;
    }

    return url.substring(bx + 1, ax) + ' / ' + name;
  },

  UrlToRepo: function(repo, path, line, rev) {
    return lib.UrlToRepo(this.repos[repo], path, line, rev);
  }

};

var RepoOption = React.createClass({
  render: function() {
    return (
      <option value={this.props.value}>{this.props.value}</option>
    )
  }
});

var SearchBar = React.createClass({
  componentWillMount: function() {
    var _this = this;
    Model.didLoadRepos.tap(function(model, repos) {
      _this.setState({ allRepos: Object.keys(repos) });
    });
  },

  componentDidMount: function() {
    this.setParams(this.props);

    if (this.hasAdvancedValues()) {
      this.showAdvanced();
    }
  },

  getInitialState: function() {
    var useDarkTheme;

    try {
      useDarkTheme = localStorage["useDarkTheme"];
      if (!useDarkTheme || useDarkTheme === "false") {
        useDarkTheme = false;
      } else {
        useDarkTheme = true;
      }
    } catch (e) {
      useDarkTheme = false;
    }


    return {
      state: null,
      allRepos: [],
      useDarkTheme: useDarkTheme
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

  getSelectedRepos() {
    var repos = "";
    if (this.refs.repos.selectedOptions.length) {
      for (var i = this.refs.repos.selectedOptions.length - 1; i >= 0; i--) {
        repos += this.refs.repos.selectedOptions[i].value + ",";
      }
    }
    return repos;
  },

  getParams: function() {
    return {
      q : this.refs.q.value.trim(),
      files : this.refs.files.value.trim(),
      repos : this.getSelectedRepos(),
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
    return this.refs.files.value.trim() !== '' || this.refs.icase.checked || this.refs.repos.value !== '';
  },

  showAdvanced: function() {
    var adv = this.refs.adv,
        ban = this.refs.ban,
        q = this.refs.q,
        files = this.refs.files;

    css(adv, 'height', 'auto');
    css(adv, 'padding', '10px 0');

    css(ban, 'max-height', '0');
    css(ban, 'opacity', '0');

    if (q.value.trim() !== '') {
      files.focus();
    }
  },

  hideAdvanced: function(evt) {
    var adv = this.refs.adv,
        ban = this.refs.ban,
        q = this.refs.q,
        files = this.refs.files,
        icase = this.refs.icase,
        repos = this.refs.repos;

    if(evt && (ReactDOM.findDOMNode(files).contains(evt.target) ||
      ReactDOM.findDOMNode(icase).contains(evt.target) ||
      ReactDOM.findDOMNode(repos).contains(evt.target)) ) {
        // Probably clicking into an input don't close advanced options
        return;
    }

    css(adv, 'height', '0');
    css(adv, 'padding', '0');

    css(ban, 'max-height', '100px');
    css(ban, 'opacity', '1');

    q.focus();
  },

  changeTheme: function(evt) {
    var useDarkTheme = evt.target.checked;
    if (useDarkTheme) {
      document.getElementById('theme').href="css/hound-dark.css";
    } else {
      document.getElementById('theme').href="css/hound-light.css";
    }
    try {
      localStorage["useDarkTheme"] = useDarkTheme;
    } catch (e) {
      console.error("Cannot set theme to localStorage!")
    }
    this.setState({useDarkTheme: useDarkTheme})
  },

  render: function() {
    var repoCount = this.state.allRepos.length,
        repoOptions = [];

    this.state.allRepos.forEach(function(repoName, index){
      repoOptions.push(<RepoOption key={repoName+"-"+index} value={repoName} />);
    });

    return (
      <div id="input">
        <div id="ina">
          <input
            id="q"
            type="text"
            placeholder="Search by Regexp"
            ref="q"
            autoComplete="off"
            onKeyDown={this.queryGotKeydown}
            onFocus={this.queryGotFocus}
            autoFocus
          />
          <div className="button-add-on">
            <button id="dodat" onClick={this.submitQuery}></button>
          </div>
        </div>

        <div id="inb">
          <div id="adv" ref="adv" onClick={this.hideAdvanced}>
            <span className="octicon octicon-chevron-up hide-adv" onClick={this.hideAdvanced}></span>
            <div className="field">
              <label>File Path</label>
              <div className="field-input">
                <input
                  type="text"
                  id="files"
                  placeholder="/regexp/"
                  ref="files"
                  onKeyDown={this.filesGotKeydown}
                  onFocus={this.filesGotFocus}
                />
              </div>
            </div>
            <div className="field">
              <label>Ignore Case</label>
              <div className="field-input">
                <input type="checkbox" ref="icase" />
              </div>
            </div>
            <div className="field">
              <label className="multiselect_label">Select Repo</label>
              <div className="field-input">
                <select id="repos" className="form-control multiselect" multiple={true} size={Math.min(16, repoCount)} ref="repos">
                  {repoOptions}
                </select>
              </div>
            </div>
          </div>
          <div className="ban" ref="ban" onClick={this.showAdvanced}>
            <em>Advanced:</em> ignore case, filter by path, stuff like that.
          </div>
        </div>
        <div className="stats">
          <div className="stats-left">
            <a href="/excluded_files.html"
              className="link-gray">
                Excluded Files
            </a>
            <label className="change-theme">
              <input type="checkbox" onChange={this.changeTheme} checked={this.state.useDarkTheme} /> Use Dark Theme
            </label>
          </div>
          {this.renderStats()}
        </div>

      </div>
    );
  },

  renderStats() {
    var stats = this.state.stats;

    if (!stats) {
      return null;
    }

    return (
      <div className="stats-right">
        <div className="val">{FormatNumber(stats.Total)}ms total</div> /
        <div className="val">{FormatNumber(stats.Server)}ms server</div> /
        <div className="val">{stats.Files} files</div>
      </div>
    );
  }
});

/**
 * Take a list of matches and turn it into a simple list of lines.
 */
var MatchToLines = function(match) {
  var lines = [],
      base = match.LineNumber,
      nBefore = match.Before.length,
      nAfter = match.After.length;
  match.Before.forEach(function(line, index) {
    lines.push({
      Number : base - nBefore + index,
      Content: line,
      Match: false
    });
  });

  lines.push({
    Number: base,
    Content: match.Line,
    Match: true
  });

  match.After.forEach(function(line, index) {
    lines.push({
      Number: base + index + 1,
      Content: line,
      Match: false
    });
  });

  return lines;
};

/**
 * Take several lists of lines each representing a matching block and merge overlapping
 * blocks together. A good example of this is when you have a match on two consecutive
 * lines. We will merge those into a singular block.
 *
 * TODO(knorton): This code is a bit skanky. I wrote it while sleepy. It can surely be
 * made simpler.
 */
var CoalesceMatches = function(matches) {
  var blocks = matches.map(MatchToLines),
      res = [],
      current;
  // go through each block of lines and see if it overlaps
  // with the previous.
  for (var i = 0, n = blocks.length; i < n; i++) {
    var block = blocks[i],
        max = current ? current[current.length - 1].Number : -1;
    // if the first line in the block is before the last line in
    // current, we'll be merging.
    if (block[0].Number <= max) {
      block.forEach(function(line) {
        if (line.Number > max) {
          current.push(line);
        } else if (current && line.Match) {
          // we have to go back into current and make sure that matches
          // are properly marked.
          current[current.length - 1 - (max - line.Number)].Match = true;
        }
      });
    } else {
      if (current) {
        res.push(current);
      }
      current = block;
    }
  }

  if (current) {
    res.push(current);
  }

  return res;
};

/**
 * Use the DOM to safely htmlify some text.
 */
var EscapeHtml = function(text) {
  var e = EscapeHtml.e;
  e.textContent = text;
  return e.innerHTML;
};
EscapeHtml.e = document.createElement('div');

/**
 * Produce html for a line using the regexp to highlight matches.
 */
var ContentFor = function(line, regexp) {
  if (!line.Match) {
    return EscapeHtml(line.Content);
  }
  var content = line.Content,
      buffer = [];

  while (true) {
    regexp.lastIndex = 0;
    var m = regexp.exec(content);
    if (!m) {
      buffer.push(EscapeHtml(content));
      break;
    }

    buffer.push(EscapeHtml(content.substring(0, regexp.lastIndex - m[0].length)));
    buffer.push( '<em>' + EscapeHtml(m[0]) + '</em>');
    content = content.substring(regexp.lastIndex);
  }
  return buffer.join('');
};

var FilesView = React.createClass({
  onLoadMore: function(event) {
    Model.LoadMore(this.props.repo);
  },

  render: function() {
    var rev = this.props.rev,
        repo = this.props.repo,
        regexp = this.props.regexp,
        matches = this.props.matches,
        totalMatches = this.props.totalMatches;

    var files = matches.map(function(match, index) {
      var filename = match.Filename,
          blocks = CoalesceMatches(match.Matches);

      var matches = blocks.map(function(block, index) {

        var lines = block.map(function(line) {
          var content = ContentFor(line, regexp);
          return (
            <div key={match.Filename+"-line-"+line.Number} className="line">
              <a
                href={Model.UrlToRepo(repo, filename, line.Number, rev)}
                className="lnum"
                target="_blank"
              >
                {line.Number}
              </a>
              <span className="lval" dangerouslySetInnerHTML={{__html:content}} />
            </div>
          );
        });

        return (
          <div key={match.Filename+"-block-"+index} className="match">{lines}</div>
        );
      });

      return (
        <div key={match.Filename} className="file">
          <div className="title">
            <a
              href={Model.UrlToRepo(repo, match.Filename, null, rev)}
              target="_blank"
            >
              {match.Filename}
            </a>
          </div>
          <div className="file-body">
            {matches}
          </div>
        </div>
      );
    });

    var more = '';
    if (matches.length < totalMatches) {
      more = (<button className="moar" onClick={this.onLoadMore}>Load all {totalMatches} matches in {Model.NameForRepo(repo)}</button>);
    }

    return (
      <div className="files">
      {files}
      {more}
      </div>
    );
  }
});

var ResultView = React.createClass({
  componentWillMount: function() {
    var _this = this;
    Model.willSearch.tap(function(model, params) {
      _this.setState({
        results: null,
        query: params.q
      });
    });
  },
  getInitialState: function() {
    return { results: null };
  },
  render: function() {
    if (this.state.error) {
      return (
        <div id="no-result" className="error">
          <strong>ERROR:</strong>{this.state.error}
        </div>
      );
    }

    if (this.state.results !== null && this.state.results.length === 0) {
      // TODO(knorton): We need something better here. :-(
      return (
        <div id="no-result">&ldquo;Nothing for you, Dawg.&rdquo;<div>0 results</div></div>
      );
    }

    if (this.state.results === null && this.state.query) {
      return (
        <div id="no-result"><img src="images/busy.gif" /><div>Searching...</div></div>
      );
    }

    var regexp = this.state.regexp,
        results = this.state.results || [];

    var repos = results.map(function(result, index) {
      return (
        <div key={result.repo+"-"+index} className="repo">
          <div className="title">
            <span className="mega-octicon octicon-repo"></span>
            <span className="name">{Model.NameForRepo(result.Repo)}</span>
          </div>
          <FilesView
            matches={result.Matches}
            rev={result.Rev}
            repo={result.Repo}
            regexp={regexp}
            totalMatches={result.FilesWithMatch}
          />
        </div>
      );
    });
    return (
      <div id="result">{repos}</div>
    );
  }
});

var App = React.createClass({
  componentWillMount: function() {
    var params = ParamsFromUrl();
    this.setState({
      q: params.q,
      i: params.i,
      files: params.files,
      repos: params.repos
    });

    var _this = this;
    Model.didSearch.tap(function(model, results, stats) {
      _this.refs.searchBar.setState({
        stats: stats
      });

      _this.refs.resultView.setState({
        results: results,
        regexp: _this.refs.searchBar.getRegExp(),
        error: null
      });
    });

    Model.didLoadMore.tap(function(model, repo, results) {
      _this.refs.resultView.setState({
        results: results,
        regexp: _this.refs.searchBar.getRegExp(),
        error: null
      });
    });

    Model.didError.tap(function(model, error) {
      _this.refs.resultView.setState({
        results: null,
        error: error
      });
    });

    window.addEventListener('popstate', function(e) {
      var params = ParamsFromUrl();
      _this.refs.searchBar.setParams(params);
      Model.Search(params);
    });
  },
  onSearchRequested: function(params) {
    this.updateHistory(params);
    Model.Search(this.refs.searchBar.getParams());
  },
  updateHistory: function(params) {
    var path = location.pathname +
      '?q=' + encodeURIComponent(params.q) +
      '&i=' + encodeURIComponent(params.i) +
      '&files=' + encodeURIComponent(params.files) +
      '&repos=' + params.repos;
    history.pushState({path:path}, '', path);
  },
  render: function() {
    return (
      <div>
        <SearchBar
          ref="searchBar"
          q={this.state.q}
          i={this.state.i}
          files={this.state.files}
          repos={this.state.repos}
          onSearchRequested={this.onSearchRequested}
        />
        <ResultView ref="resultView" q={this.state.q} />
      </div>
    );
  }
});

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
Model.Load();
