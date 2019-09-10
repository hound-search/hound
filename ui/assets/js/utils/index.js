export var css = function(el, n, v) {
    el.style.setProperty(n, v, '');
};

export var FormatNumber = function(t) {
    var s = '' + (t|0),
        b = [];
    while (s.length > 0) {
        b.unshift(s.substring(s.length - 3, s.length));
        s = s.substring(0, s.length - 3);
    }
    return b.join(',');
};

export var ParamsFromQueryString = function(qs, params) {
    params = params || {};

    if (!qs) {
        return params;
    }

    qs.substring(1).split('&').forEach(function(v) {
        var pair = v.split('=');
        if (pair.length != 2) {
            return;
        }

        // Handle classic '+' representation of spaces, such as is used
        // when Hound is set up in Chrome's Search Engine Manager settings
        pair[1] = pair[1].replace(/\+/g, ' ');

        params[decodeURIComponent(pair[0])] = decodeURIComponent(pair[1]);
    });

    return params;
};

export var ParamsFromUrl = function(params) {
    params = params || {
        q: '',
        i: 'nope',
        files: '',
        repos: '*'
    };
    return ParamsFromQueryString(location.search, params);
};

export var ParamValueToBool = function(v) {
    v = v.toLowerCase();
    return v == 'fosho' || v == 'true' || v == '1';
};

export /**
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
export var CoalesceMatches = function(matches) {
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
export var EscapeHtml = function(text) {
    var e = EscapeHtml.e;
    e.textContent = text;
    return e.innerHTML;
};
EscapeHtml.e = document.createElement('div');

/**
 * Produce html for a line using the regexp to highlight matches.
 */
export var ContentFor = function(line, regexp) {
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
