import React, { Component } from 'react';
import axios from 'axios';

import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import TimeTicker from './TimeTicker';
import { Grid } from '@material-ui/core';

const styles = theme => ({
  root: {
    width: '100%',
    overflowX: 'auto',
  },
  table: {
    minWidth: 700,
  },
  header: {
    padding: theme.spacing.unit*3,
  },
  light: {
    color: '#737373',
  }
});

class ContainerList extends Component {
    state = {
        containers: [],
    };

    componentDidMount() {   
        axios.get(this.props.node.address)
        .then(resp => {
          this.setState({containers: resp.data});
        });
    }
    
    render() {
        const { classes, node } = this.props;
        const { containers } = this.state;
        return (
            <Paper className={classes.root}>
                <div className={classes.header}>
                    <Grid container spacing={0} alignItems="center">
                        <Grid item xs={10}><b>{node.name}</b></Grid>
                        <Grid item xs={2}>
                            <Grid container spacing={0} alignItems="center">
                                <Grid item xs={6}><small className={classes.light}>CPU:</small></Grid>
                                <Grid item xs={6}>{node.cpu} <small>MHz</small></Grid>
                                <Grid item xs={6}><small className={classes.light}>Memory:</small></Grid>
                                <Grid item xs={6}>{node.memory} <small>MB</small></Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </div>
                <Table className={classes.table}>
                    <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell align="center">CPU</TableCell>
                        <TableCell align="center">Memory</TableCell>
                        <TableCell align="center">Run Time</TableCell>
                        <TableCell align="center">Alloc Time</TableCell>
                    </TableRow>
                    </TableHead>
                    <TableBody>
                    {containers.map(row => {
                        return (
                            <TableRow key={row.ID}>
                                <TableCell component="th" scope="row">{row.Name}</TableCell>
                                <TableCell align="center">{row.CPUShares}</TableCell>
                                <TableCell align="center">{row.Memory}</TableCell>
                                <TableCell align="center"><TimeTicker start={row.Start} stop={row.Stop} /></TableCell>
                                <TableCell align="center"><TimeTicker start={row.Create} stop={row.Destroy} /></TableCell>
                            </TableRow>
                        );
                    })}
                    </TableBody>
                </Table>
            </Paper>
        );
    }
}

ContainerList.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ContainerList);