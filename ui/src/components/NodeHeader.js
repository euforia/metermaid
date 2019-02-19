import React, { Component } from 'react';
import { withStyles } from '@material-ui/core';
import { Grid, Typography, } from '@material-ui/core';

const styles = theme => ({
    header: {
        padding: theme.spacing.unit*3,
    },
});

const memMB = (d) => {
    return Math.floor(d/(1024*1024))
}

class NodeHeader extends Component {
    render() {
        const { classes, node } = this.props;
        return (
            <div className={classes.header}>
                <Grid container spacing={0} alignItems="center">
                    <Grid item xs={5}>
                        <Typography variant="subtitle1">{node.Name}</Typography>
                        <Typography variant="caption">{node.Address}</Typography>
                    </Grid>
                    <Grid item xs={4}></Grid>
                    <Grid item xs={3}>
                        <Grid container spacing={0} alignItems="center">
                            <Grid item xs={6}><Typography variant="caption">Platform:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" style={{display:'inline'}}>{node.Platform.Name} </Typography>
                                <Typography variant="caption" style={{display:'inline'}}>{node.Platform.Version}</Typography>
                            </Grid>
                            <Grid item xs={6}><Typography variant="caption">CPU:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" style={{display:'inline'}}>{node.CPUShares} </Typography>
                                <Typography variant="caption" style={{display:'inline'}}>shares</Typography>
                            </Grid>
                            <Grid item xs={6}><Typography variant="caption">Memory:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" style={{display:'inline'}}>{node.Memory === 0 ? 0 : memMB(node.Memory)} </Typography>
                                <Typography variant="caption" style={{display:'inline'}}>MB</Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </div>
        );
    }
}

export default withStyles(styles)(NodeHeader);