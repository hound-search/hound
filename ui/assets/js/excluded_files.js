import {UrlToRepo} from './common';

var ExcludedRow = React.createClass({
  render: function() {
    var url = UrlToRepo(this.props.repo, this.props.file.Filename, this.props.rev);
    return (
      <tr>
        <td className="name">
          <a href={url}>{this.props.file.Filename}</a>
        </td>
        <td className="reason">{this.props.file.Reason}</td>
      </tr>
    );
  }
});

var ExcludedTable = React.createClass({
  render: function() {
    var _this = this;
    if (this.props.searching) {
      return (<div id="no-result"><img src="images/busy.gif" /><div>Searching...</div></div>);
    }

    var rows = [];
    this.props.files.forEach(function(file) {
      rows.push(<ExcludedRow file={file} repo={_this.props.repo} />);
    });

    return (
      <table>
          <thead>
              <tr>
                  <th>Filename</th>
                  <th>Reason</th>
              </tr>
          </thead>
          <tbody className="list">{rows}</tbody>
      </table>
    );
  }
});

var RepoButton = React.createClass({
  handleClick: function(repoName) {
    this.props.onRepoClick(repoName);
  },
  render: function() {
    var className = 'repo-button';
    if (this.props.currentRepo === this.props.repo) {
      className += ' selected';
    }

    return (
      <button onClick={this.handleClick.bind(this, this.props.repo)} className={className}>
        {this.props.repo}
      </button>
    );
  }
});

var RepoList = React.createClass({
  render: function() {
    var repos = [],
        _this = this;
    this.props.repos.forEach(function(repo){
      repos.push(<RepoButton repo={repo} onRepoClick={_this.props.onRepoClick} currentRepo={_this.props.repo} />);
    });

    return (
      <div id="repolist">
        {repos}
      </div>
    );
  }
});

var FilterableExcludedFiles = React.createClass({
  getInitialState: function() {
    var _this = this;
    $.ajax({
      url: 'api/v1/repos',
      dataType: 'json',
      success: function(data) {
        _this.setState({ repos: data });
      },
      error: function(xhr, status, err) {
        // TODO(knorton): Fix these
        console.error(err);
      }
    });

    return {
      files: [],
      repos: [],
      repo: null,
    };
  },

  onRepoClick: function(repo) {
    var _this = this;
    _this.setState({
      searching: true,
      repo: this.state.repos[repo],
    });
    $.ajax({
      url: 'api/v1/excludes',
      data: {repo: repo},
      type: 'GET',
      dataType: 'json',
      success: function(data) {
        _this.setState({ files: data, searching: false });
      },
      error: function(xhr, status, err) {
        // TODO(knorton): Fix these
        console.error(err);
      }
    });
  },

  render: function() {
    return (
      <div id="excluded_container">
        <a href="/">Home</a>
        <h1>Excluded Files</h1>

        <div id="excluded_files" className="table-container">
          <RepoList repos={Object.keys(this.state.repos)} onRepoClick={this.onRepoClick} repo={this.state.repo} />
          <ExcludedTable files={this.state.files} searching={this.state.searching} repo={this.state.repo} />
        </div>
      </div>
    );
  }
});

React.renderComponent(
  <FilterableExcludedFiles />,
  document.getElementById('root')
);
