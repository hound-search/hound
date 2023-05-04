/**
 * Signal represents a collection of observers that can be notified with
 * parameterized event dispatch. Listeners can tap into a signal to
 * receive callbacks when a signal is raised.
 */
export class Signal {
    constructor() {
        this.listeners = [];
    }

    /**
     * Begin listening to this signal with the given callback, l.
     *
     * @param {function} l
     */
    tap(l) {
        // Make a copy of the listeners to avoid the all too common
        // subscribe-during-dispatch problem
        this.listeners = [...this.listeners, l];
    }

    /**
     * Remove the listener, l, from the active listeners of this signal. If
     * l is not currently tapped into this signal, this method will complete
     * successfully without changing the collection of listeners.
     *
     * @param {function} l
     */
    untap(l) {
        const ix = this.listeners.indexOf(l);
        if (ix == -1) {
            return;
        }

        // Make a copy of the listeners to avoid the all to common
        // unsubscribe-during-dispatch problem
        this.listeners = this.listeners.slice(0);
        this.listeners.splice(ix, 1);
    }

    /**
     * Raise this signal by calling all the tapped listeners forwarding all
     * parameters.
     *
     * @param  {...any} args
     */
    raise(...args) {
        this.listeners.forEach((l) => l.apply(this, args));
    }
}
