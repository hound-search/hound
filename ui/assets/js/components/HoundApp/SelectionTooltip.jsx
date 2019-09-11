import React, { useState, useEffect } from 'react';
import { SelectionManager } from '../../helpers/SelectionManager';

export const SelectionTooltip = (props) => {

    const { delay } = props;

    const supported = SelectionManager.Supported();
    const [ active, setActive ] = useState(false);
    const [ data, setData ] = useState({
        url: '',
        top: 0,
        left: 0,
        text: ''
    });

    let timeoutDelay;

    useEffect(() => {

        if ( supported ) {

            document.addEventListener('click',  () => {

                clearTimeout(delay);

                timeoutDelay = setTimeout(() => {

                    const selection = SelectionManager.GetSelection();

                    if (selection) {
                        setData(selection);
                        setActive(true);
                    } else {
                        setActive(false);
                    }

                }, delay);

            });

        }

    }, []);

    const onClickTooltip = (e) => {
        e.stopPropagation();
        setTimeout( () => {
            SelectionManager.clearSelection();
            setActive(false);
        }, 100);
    };

    const element = supported
        ? (
            <a
                className={ `selection-tooltip octicon octicon-search${ active && ' active' || ''}` }
                href={ data.url }
                style={{ top: data.top, left: data.left }}
                onClick={ onClickTooltip }
                target='_blank'
                rel="noopener noreferrer"
            >
                <span className="selection-tooltip-text">
                    { data.text }
                </span>
            </a>
        )
        : ''

    return element;

};
