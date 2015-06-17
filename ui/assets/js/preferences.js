/** @jsx React.DOM */

var PreferenceLine = React.createClass({
  handleClick: function(preferenceName) {
    this.props.onPreferenceClick(preferenceName);
  },
  render: function() {
    var classname = 'field';
  }
});

var ParamValueToBool = function(v) {
  if(v == null) {
    return false;
  }

  v = v.toLowerCase();
  return v == 'fosho' || v == 'true' || v == '1';
};

var PreferencesView = React.createClass({
  ignoreCaseChanged: function(event) {
    var isChecked = this.refs.icase.getDOMNode().checked;
    localStorage.setItem('ignoreCase', isChecked);
    this.setState({
      i: isChecked
    });
  },
  autoHideAdvanced: function(event) {
    var isChecked = this.refs.autoHideAdv.getDOMNode().checked;
    localStorage.setItem('autoHideAdvanced', isChecked);
    this.setState({
      autoHideAdv: isChecked
    });
  },
  initPreferences: function() {
    // Should already be set from the main page

    var ignoreCase = localStorage.getItem('ignoreCase');
    var hideAdvanced = localStorage.getItem('autoHideAdvanced');

    if(ignoreCase == null) {
      localStorage.setItem('ignoreCase', isChecked);
    }
    if(hideAdvanced == null) {
      localStorage.setItem('autoHideAdvanced', isChecked);
    }
  },
  render: function() {
    var ignoreCase = ParamValueToBool(localStorage.getItem('ignoreCase'));
    var hideAdvanced = ParamValueToBool(localStorage.getItem('autoHideAdvanced'));

    return (
      <div id="preferences_container">
        <a href="/">Home</a>
        <h1>Preferences</h1>
        <div id="inb">
          <div id="preferences">
            <div className="field">
              <label>Ignore Case</label>
              <div className="field-input">
                <input type="checkbox" ref="icase" checked={ignoreCase} onClick={this.ignoreCaseChanged} />
              </div>
            </div>
            <div className="field">
              <label>Hide Advanced on Search</label>
              <div className="field-input">
                <input type="checkbox" ref="autoHideAdv" checked={hideAdvanced} onClick={this.autoHideAdvanced} />
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
});

React.renderComponent(
  <PreferencesView />,
  document.getElementById('root')
);
