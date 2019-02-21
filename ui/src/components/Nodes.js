import React, { Component } from 'react';
import { Grid } from '@material-ui/core';
import NodeHeader from './NodeHeader';

class Nodes extends Component{
    render() {
        const {data} = this.props;
        return (
            <div>
                {data.map(node => {
                   return <NodeHeader node={node} />
                })}
            </div>
        );
    }
}

export default (Nodes);