import { closestElement } from '../utils';

export const SelectionManager = {

    Supported () {
        return (
            window &&
            'getSelection' in window &&
            document &&
            document.body &&
            'getBoundingClientRect' in document.body
        );
    },

    GetSelection () {

        const selection = window.getSelection();
        const anchorNode = selection.anchorNode;
        const selectionText = selection.toString().trim();
        const newLineReg = /[\r\n]+/;
        const escapeReg = /[.?*+^$[\]\\(){}|-]/g;
        const urlReg = /([\?&])q=([^&$]+)/;

        if ( selectionText.length && !newLineReg.test(selectionText) && closestElement(anchorNode, 'lval') ) {

            const url = window.location.href;
            const escapedText = encodeURIComponent(selectionText.replace(escapeReg, '\\$&'));
            const searchURL = url.replace(urlReg, '$1q=' + escapedText);

            const selectionRange = selection.getRangeAt(0);
            const selectionRect = selectionRange.getBoundingClientRect();
            const scrollTop = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;

            return {
                text: selectionText,
                url: searchURL,
                left: selectionRect.left + selectionRect.width + 5,
                top: selectionRect.top + scrollTop + 5
            };

        }

        return null;

    },

    clearSelection () {
        const selection = window.getSelection();
        selection.removeAllRanges();
    }

};
