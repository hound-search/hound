export var Signal = function() {
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
