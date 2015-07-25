/** @jsx React.DOM */

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

var EncodeURLQuery = function(query, regex, wholeWord) {
  if (!(regex == true)) {
    // copied from stackoverflow / google's closure lib
    query = query.replace(/([-()\[\]{}+?*.$\^|,:#<!\\])/g, '\\$1');
  }
  if (wholeWord) {
    query = "\\b" + query + "\\b";
  }
  return query;
};

var DecodeURLQuery = function(query, regex, wholeWord) {
  if (wholeWord) {
    query = query.replace(/\\b/g, "");
  }
  if (!regex) {
    query = query.replace(/\\(.)/mg, "$1");
  }
  return query;
};

var EncodeURLFiles = function(files, mode) {
  if (mode === "simple") {
    files = files.replace(/\./g, "\\.");
    files = files.replace(/\*/g, ".*");
    return files.replace(/\,\s*/g, "|");
  } else {
    return files;
  }
};

var DecodeURLFiles =  function(files, mode) {
  if (mode === "simple") {
    files = files.replace(/\|/g, ", ");
    files = files.replace(/\.\*/g, "*");
    return files.replace(/\\\./g, ".");
  } else {
    return files;
  }
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
  } else {
    params["repos"] = params["repos"].toString();
  }

  return params;
};

var ParamsFromUrl = function(params) {
  params = params || {
    q: '',
    r: '',
    i: '',
    w: '',
    mode: 'regex',
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

    params = $.extend(params, { rng: ':20', stats: 'fosho' });

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
      _this.didSearch.raise(_this, _this.results);
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

var SearchPanel = React.createClass({
  componentWillMount: function() {
    this.setState({
      q: DecodeURLQuery(this.props.q, this.props.r, this.props.w),
      r: this.props.r,
      i: this.props.i,
      w: this.props.w,
      mode: this.props.mode,
      files: this.props.files,
      repos: this.props.repos
    });
  },

  getInitialState: function() {
    return {
      q: '',
      r: false,
      i: false,
      w: false,
      mode: 'regex',
      files: '',
      repos: ''
    };
  },

  submitQuery: function(event) {
    // deep-copy current state and encode query and files
    encodedState = $.extend(true, {}, this.state);
    encodedState.q = this.getEncodedQuery();
    encodedState.files = EncodeURLFiles(this.state.files, this.state.mode);
    this.props.onSearchRequested(encodedState);

    return false;
  },

  getEncodedQuery: function() {
    return EncodeURLQuery(this.state.q, this.state.r, this.state.w);
  },

  getRegExp: function() {
    var flags = (this.state.i) ? 'ig' : 'g';
    return new RegExp(this.getEncodedQuery(), flags);
  },

  updateState: function(newState) {
    this.setState(newState);
  },

  render: function() {
    var statsView = '';
    if (this.state.stats) {
      statsView = (
        <div className="stats">
          <div className="stats-left">
            <a href="/excluded_files.html" className="link-gray">Excluded Files</a>
          </div>
          <div className="stats-right">
            <div className="val">{FormatNumber(this.state.stats.Total)}ms total</div> /
            <div className="val">{FormatNumber(this.state.stats.Server)}ms server</div> /
            <div className="val">{this.state.stats.Files} files</div>
          </div>
        </div>
      );
    }

    var searchParams = { q: this.state.q,
                         i: this.state.i,
                         r: this.state.r,
                         w: this.state.w };
    var advancedParams = { mode:  this.state.mode,
                           files: this.state.files,
                           repos: this.state.repos };

    return (
      <form id="input" onSubmit={this.submitQuery}>
        <SearchBar ref="searchBar"
                   params={searchParams}
                   updatePanel={this.updateState} />
        <Advanced ref="advanced"
                  params={advancedParams}
                  updatePanel={this.updateState} />
        {statsView}
      </form>
    );
  }
});

var SearchBar = React.createClass({
  componentWillMount: function() {
    this.setState({ query: this.props.params.q });
  },

  componentDidMount: function() {
    var q = this.refs.q.getDOMNode();

    // TODO(knorton): Can't set this in jsx
    q.setAttribute('autocomplete', 'off');
    q.focus();
  },

  handleChange: function(event) {
    this.props.updatePanel({ q: event.target.value.trim() });

    this.setState({ query: event.target.value });
  },

  render: function() {
    return (
      <div id="ina">
        <div id="toggles">
          <SearchToggle id="regex-toggle"
                        src="images/regex.svg"
                        option="r"
                        title="Regex"
                        initialChecked={this.props.params.r}
                        updatePanel={this.props.updatePanel} />
          <SearchToggle id="case-toggle"
                        src="images/case.svg"
                        title="Case-sensitive"
                        option="i"
                        initialChecked={this.props.params.i}
                        updatePanel={this.props.updatePanel} />
          <SearchToggle id="word-toggle"
                        src="images/wholeWord.svg"
                        title="Whole word"
                        option="w"
                        initialChecked={this.props.params.w}
                        updatePanel={this.props.updatePanel} />
        </div>
        <input type="text"
               ref="q"
               autocomplete="off"
               placeholder="Search by regex or simple match"
               value={this.state.query}
               onChange={this.handleChange} />
        <div className="button-add-on">
          <button id="dodat" type="submit"></button>
        </div>
      </div>
    );
  }
});

var SearchToggle = React.createClass({
  getInitialState: function() {
    return {
      checked: this.props.initialChecked
    };
  },

  handleChange: function(event) {
    var newState = {};
    newState[this.props.option] = event.target.checked;
    this.props.updatePanel(newState);

    this.setState({ checked: event.target.checked });
  },

  render: function() {
    return (
      <div className="toggle">
        <input type="checkbox"
               id={this.props.id}
               checked={this.state.checked}
               onChange={this.handleChange} />
        <label htmlFor={this.props.id} title={this.props.title}>
          <img src={this.props.src} alt={this.props.id} />
        </label>
      </div>
    );
  }
});

/* The Advanced dropdown gives the users two modes of searching
 * their files: regex and simple. Regex lets the user simply enter
 * a regex of their choosing to match filePaths while Simple lets the user
 * user a sublime-familiar interface of *.js and filepaths to match
 * their selection.
 */
var Advanced = React.createClass({
  componentDidMount: function() {
    if (this.hasValues()) this.show();
  },

  getInitialState: function() {
    return {
      showAdv: false,
      showBan: true
    };
  },

  hasValues: function() {
    return this.props.params.repos != '*' || this.props.params.files != '';
  },


  toggleShown: function() {
    this.setState({ showAdv: !this.state.showAdv, showBan: !this.state.showBan });
  },

  show: function() {
    this.setState({ showAdv: true, showBan: false });
  },

  render: function() {
    var showAdv = (this.state.showAdv) ? 'is-shown' : '';
    var showBan = (this.state.showBan) ? 'ban is-shown' : 'ban';

    return (
      <div id="inb">
        <div id="adv" className={showAdv}>
          <span className="octicon octicon-chevron-up"
                onClick={this.toggleShown} />
          <div className="field">
            <label>File Path:</label>
            <RepoSelect repos={this.props.params.repos}
                        updatePanel={this.props.updatePanel}
                        showAdvanced={this.show} />
            <span className="slash-delimiter">/</span>
            <FilePath files={this.props.params.files}
                      mode={this.props.params.mode}
                      updatePanel={this.props.updatePanel}
                      showAdvanced={this.show} />
          </div>
        </div>
        <div className={showBan} onClick={this.toggleShown}>
          <span className="octicon octicon-chevron-down" />
          <em> Advanced:</em> filter by repo(s) or file-path, stuff like that.
        </div>
      </div>
    );
  }
});

var FilePath = React.createClass({
  componentWillMount: function() {
    var decodedFiles = DecodeURLFiles(this.props.files, this.props.mode);
    this.setState({ files: decodedFiles, mode: this.props.mode });
  },

  switchMode: function() {
    newMode = (this.state.mode != 'regex') ? 'regex' : 'simple';
    this.props.updatePanel({ mode: newMode });

    this.setState({ mode: newMode });
  },

  handleChange: function(event) {
    this.props.updatePanel({ files: event.target.value.trim() });
    this.setState({ files: event.target.value });
  },

  render: function() {
    var altText = (this.state.mode == 'regex') ? "htdocs\\/.*\\.(php|js)$"
                                               : "htdocs/*.php, htdocs/*.js";
    var sliderState = (this.state.mode == 'regex') ? "slider" : "slider simple";
    var showTip = (this.state.mode == 'regex') ? "simple-tip"
                                               : "simple-tip is-shown";

    return (
      <div className="field-input file-path">
        <input type="text"
               value={this.state.files}
               placeholder={altText}
               onFocus={this.props.showAdvanced}
               onChange={this.handleChange} />
        <button type="button"
                className="file-mode"
                onClick={this.switchMode}>
          <div className="regex">regex</div>
          <div className="simple">simple</div>
          <div className={sliderState}></div>
        </button>
        <div className={showTip}>Tip: you can use * as a wildcard, and
        comma-delimit multiple options, e.g. "*.js, *.py"</div>
      </div>
    );
  }
});

var RepoSelect = React.createClass({
  componentWillMount: function() {
    this.setState({ value: this.props.repos.toString() });

    var _this = this;
    Model.didLoadRepos.tap(function(model, repos) {
      _this.setState({ allRepos: Object.keys(repos) });
    });
  },

  getInitialState: function() {
    return {
      showRepos: false,
      allRepos: []
    };
  },

  toggleRepos: function() {
    this.setState({ showRepos: !this.state.showRepos });
  },

  showRepos: function() {
    this.setState({ showRepos: true });
  },

  preview: function() {
    var dir = (this.state.showRepos) ? 'up' : 'down';
    var chevron = '<div class="octicon octicon-chevron-'+dir+'"></div>';
    if (this.state.value == '' || this.state.value == '*') {
      return { __html: '<span>All Repos</span>'+chevron };
    }
    return { __html: '<span>'+this.state.value+'</span>'+chevron };
  },

  focusSelect: function(event) {
    if (event.keyCode != 9 && event.keyCode != 32) {
      this.showRepos();
      this.refs.input.getDOMNode().focus();
    }
  },

  handleChange: function(event) {
    var newState = {};
    var selected = [].map.call(event.target.selectedOptions, function(opt) {
      return opt.value;
    }).toString();

    this.props.updatePanel({ repos: selected });
    this.setState({ value: selected });
  },

  render: function() {
    var repoOptions = [];
    this.state.allRepos.forEach(function(repoName){
      repoOptions.push(<RepoOption value={repoName} initial={this.state.value} />);
    }, this);

    var showRepos = (this.state.showRepos) ? 'multiselect is-shown' : 'multiselect';
    var initialValue = (this.props.repos) ? this.props.repos.toString() : '';

    return (
      <div className="field-input repos">
        <button type="button"
                className="repo-preview"
                onClick={this.toggleRepos}
                onKeyDown={this.focusSelect}
                onFocus={this.props.showAdvanced}
                dangerouslySetInnerHTML={this.preview()} />
        <select ref="input"
                className={showRepos}
                multiple={true}
                size={Math.min(10, this.state.allRepos.length)}
                defaultValue={initialValue}
                onChange={this.handleChange}>
          {repoOptions}
        </select>
      </div>
    );
  }
});

var RepoOption = React.createClass({
  isSelected: function() {
    return this.props.initial.split(',').indexOf(this.props.value) > -1;
  },

  render: function() {
    return (
      <option value={this.props.value} selected={this.isSelected()}>{this.props.value}</option>
    )
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
      var matches = blocks.map(function(block) {
        var lines = block.map(function(line) {
          var content = ContentFor(line, regexp);
          return (
            <div className="line">
              <a href={Model.UrlToRepo(repo, filename, line.Number, rev)}
                  className="lnum"
                  target="_blank">{line.Number}</a>
              <span className="lval" dangerouslySetInnerHTML={{__html:content}} />
            </div>
          );
        });

        return (
          <div className="match">{lines}</div>
        );
      });

      return (
        <div className="file">
          <div className="title">
            <a href={Model.UrlToRepo(repo, match.Filename, null, rev)}>
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
        <div className="repo">
          <div className="title">
            <span className="mega-octicon octicon-repo"></span>
            <span className="name">{Model.NameForRepo(result.Repo)}</span>
          </div>
          <FilesView matches={result.Matches}
              rev={result.Rev}
              repo={result.Repo}
              regexp={regexp}
              totalMatches={result.FilesWithMatch} />
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
      r: ParamValueToBool(params.r),
      i: ParamValueToBool(params.i),
      w: ParamValueToBool(params.w),
      mode: params.mode,
      files: params.files,
      repos: params.repos
    });

    var _this = this;
    Model.didSearch.tap(function(model, results, stats) {
      _this.refs.searchPanel.setState({
        stats: stats
      });

      _this.refs.resultView.setState({
        results: results,
        regexp: _this.refs.searchPanel.getRegExp(),
        error: null
      });
    });

    Model.didLoadMore.tap(function(model, repo, results) {
      _this.refs.resultView.setState({
        results: results,
        regexp: _this.refs.searchPanel.getRegExp(),
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
      Model.Search(params);
    });
  },

  onSearchRequested: function(params) {
    this.updateHistory(params);
    Model.Search(params);
  },

  updateHistory: function(params) {
    var path = location.pathname +
               '?q=' + encodeURIComponent(params.q) +
               '&r=' + encodeURIComponent(params.r) +
               '&i=' + encodeURIComponent(params.i) +
               '&w=' + encodeURIComponent(params.w) +
               '&mode=' + encodeURIComponent(params.mode) +
               '&files=' + encodeURIComponent(params.files) +
               '&repos=' + params.repos;
    history.pushState({path:path}, '', path);
  },

  render: function() {
    return (
      <div>
        <SearchPanel ref="searchPanel"
                     q={this.state.q}
                     r={this.state.r}
                     i={this.state.i}
                     w={this.state.w}
                     mode={this.state.mode}
                     files={this.state.files}
                     repos={this.state.repos}
                     onSearchRequested={this.onSearchRequested} />
        <ResultView ref="resultView"
                    q={this.state.q} />
      </div>
    );
  }
});

React.renderComponent(
  <App />,
  document.getElementById('root')
);
Model.Load();
