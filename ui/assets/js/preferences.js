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
    docCookies.setItem('ignoreCase', isChecked);
    this.setState({
      i: isChecked
    });
  },
  autoHideAdvanced: function(event) {
    var isChecked = this.refs.autoHideAdv.getDOMNode().checked;
    docCookies.setItem('autoHideAdvanced', isChecked);
    this.setState({
      autoHideAdv: isChecked
    });
  },
  initPreferences: function() {
    // Should already be set from the main page
    var ignoreCase = docCookies.getItem('ignoreCase');
    var hideAdvanced = docCookies.getItem('autoHideAdvanced');

    if(ignoreCase == null) {
      docCookies.setItem('ignoreCase', false);
    }
    if(hideAdvanced == null) {
      docCookies.setItem('autoHideAdvanced', false);
    }
  },
  render: function() {
    var autoHideAdvanced = document.cookie.autoHideAdvanced;
    var ignoreCase = document.cookie.ignoreCase;

    return (
      <div id="preferences_container">
        <a href="/">Home</a>
        <h1>Preferences</h1>
        <div id="inb">
          <div id="preferences">
            <div className="field">
              <label>Ignore Case</label>
              <div className="field-input">
                <input type="checkbox" ref="icase" checked={ParamValueToBool(docCookies.getItem('ignoreCase'))} onClick={this.ignoreCaseChanged} />
              </div>
            </div>
            <div className="field">
              <label>Hide Advanced on Search</label>
              <div className="field-input">
                <input type="checkbox" ref="autoHideAdv" checked={ParamValueToBool(docCookies.getItem('autoHideAdvanced'))} onClick={this.autoHideAdvanced} />
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
