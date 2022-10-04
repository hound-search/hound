import { Signal } from "./signal";

describe("Signal", () => {
    test("Raising untapped Signal succeeds", () => {
        const sig = new Signal();
        sig.raise(420);
    });

    test("Signal taps are notified on raise", () => {
        const sig = new Signal();
        let a, b;
        sig.tap((v) => (a = v));
        sig.tap((v) => (b = v));
        sig.raise(420);
        expect(a).toBe(420);
        expect(b).toBe(420);
    });

    test("Untapping a signal stops notifications", () => {
        const sig = new Signal();

        let a;
        const af = (v) => (a = v);

        let b;
        const bf = (v) => (b = v);

        sig.tap(af);
        sig.tap(bf);
        sig.raise(420);

        expect(a).toBe(420);
        expect(b).toBe(420);

        sig.untap(bf);
        sig.raise(666);

        expect(a).toBe(666);
        expect(b).toBe(420);

        sig.untap(af);
        sig.raise(800);

        expect(a).toBe(666);
        expect(b).toBe(420);
    });

    test("Untapping during dispatch delivers", () => {
        const sig = new Signal();

        let a;
        const af = (v) => {
            a = v;
            sig.untap(af);
        };

        let b;
        const bf = (v) => {
            b = v;
            sig.untap(bf);
        };

        sig.tap(af);
        sig.tap(bf);
        sig.raise(420);

        expect(a).toBe(420);
        expect(b).toBe(420);

        sig.raise(666);

        expect(a).toBe(420);
        expect(b).toBe(420);
    });

    test("Tapping during dispatch terminates", () => {
        const sig = new Signal();

        const af = () => {
            sig.tap(af);
        };

        // this ensures that raise only dispatches on the snapshot of
        // listeners that are tapped when raise is called. If raise
        // sees listeners that are added during dispatch, it would be
        // stuck in an infinite iteration.
        sig.tap(af);
        sig.raise();
        expect(sig.listeners.length).toBe(2);
    });
});
