import React, { Component } from 'react';
import { withStyles } from '@material-ui/core';
import { Grid, Typography, Chip} from '@material-ui/core';

const styles = theme => ({
    header: {
        padding: theme.spacing.unit*3,
    },
    tag: {
        margin: theme.spacing.unit/4,
        color: '#757575',
        height: 22,
        fontSize: 11,
    },
    tagList: {
        paddingTop: theme.spacing.unit*2,
    }
});

const memMB = (d) => {
    return Math.floor(d/(1024*1024))
}

const mapToList = (m) => {
    var out = [];
    for (var k in m) {
        out.push({key:k,value:m[k]});
    }
    return out;
}

class NodeHeader extends Component {
    render() {
        const { classes, node } = this.props;
        const tags = mapToList(node.Meta);

        return (
            <div className={classes.header}>
                <Grid container spacing={0} alignItems="center">
                    <Grid item xs={6}>
                        <Typography variant="title">{node.Name}</Typography>
                        {/* <Typography variant="caption">{node.Address}</Typography> */}
                    </Grid>
                    <Grid item xs={3}>
                        <Grid container spacing={0} alignItems="center">
                            {/* <Grid item xs={6}><Typography variant="caption">Platform:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" style={{display:'inline'}}>{node.Platform.Name} </Typography>
                                <Typography variant="caption" style={{display:'inline'}}>{node.Platform.Version}</Typography>
                            </Grid> */}
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
                    <Grid item xs={3}>
                        <Grid container spacing={0} alignItems="center">
                            <Grid item xs={6}><Typography variant="caption">Platform:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" style={{display:'inline'}}>{node.Platform.Name} </Typography>
                                <Typography variant="caption" style={{display:'inline'}}>{node.Platform.Version}</Typography>
                            </Grid>
                            <Grid item xs={6}><Typography variant="caption">Address:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2">{node.Address} </Typography>
                                {/* <Typography variant="caption" style={{display:'inline'}}>shares</Typography> */}
                            </Grid>
                            {/* <Grid item xs={6}><Typography variant="caption">Memory:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography variant="body2" style={{display:'inline'}}>{node.Memory === 0 ? 0 : memMB(node.Memory)} </Typography>
                                <Typography variant="caption" style={{display:'inline'}}>MB</Typography>
                            </Grid> */}
                        </Grid>
                    </Grid>
                    <Grid item xs={12} className={classes.tagList}>
                        {tags.map(item => {
                            return (
                                <Chip label={item.key+': '+item.value} variant="outlined"
                                    key={item.key} className={classes.tag}/>
                            );
                        })}
                    </Grid>
                </Grid>
            </div>
        );
    }
}

export default withStyles(styles)(NodeHeader);