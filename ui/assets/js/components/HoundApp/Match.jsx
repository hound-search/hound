import React from 'react';
import { ContentFor } from '../../utils';
import { Line } from './Line';

export const Match = (props) => {

    const { block, repo, regexp, rev, filename } = props;

    const lines = block.map((line, index) => {
        const content = ContentFor(line, regexp);
        return (
            <Line
                key={`line-${index}`}
                line={ line }
                rev={ rev }
                repo={ repo }
                filename={ filename }
                content={ content }
            />
        );
    });

    return (
        <div className="match">
            { lines }
        </div>
    );

};
