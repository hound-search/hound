export const Signal = function () {
};

Signal.prototype = {
    listeners : [],

    tap (l) {
        // Make a copy of the listeners to avoid the all too common
        // subscribe-during-dispatch problem
        this.listeners = this.listeners.slice(0);
        this.listeners.push(l);
    },

    untap (l) {
        const ix = this.listeners.indexOf(l);
        if (ix == -1) {
            return;
        }

        // Make a copy of the listeners to avoid the all to common
        // unsubscribe-during-dispatch problem
        this.listeners = this.listeners.slice(0);
        this.listeners.splice(ix, 1);
    },

    raise () {
        const args = Array.prototype.slice.call(arguments, 0);
        this.listeners.forEach((l) => {
            l.apply(this, args);
        });
    }
};
