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
    color: '#757575',
  }
});


function nodeFromHeaders(headers) {
    return {
        name: headers['node-name'],
        address: headers['node-addr'],
        cpu: parseInt(headers['node-cpu']),
        memory: parseInt(headers['node-memory']),
    };
}

class ContainerList extends Component {
    state = {
        node: {},
        containers: [],
    };

    componentDidMount() {   
        axios.get(this.props.source)
        .then(resp => {
          var node = nodeFromHeaders(resp.headers);
          var data = resp.data;
          for (var i = 0; i < data.length; i++) {
            data[i].PercentMemory = data[i].Memory > 0 ? ((data[i].Memory/node.memory)*100).toFixed(0) : 100;
            data[i].PercentCPU = data[i].CPUShares > 0 ? ((data[i].CPUShares/node.cpu)*100).toFixed(0) : 100;
          }
          this.setState({containers: resp.data, node: node});
        });
    }
    
    render() {
        const { classes } = this.props;
        const { containers, node } = this.state;
        return (
            <Paper className={classes.root}>
                <div className={classes.header}>
                    <Grid container spacing={0} alignItems="center">
                        <Grid item xs={9}>
                            <div>{node.name}</div>
                            <div><small className={classes.light}>{node.address}</small></div>
                        </Grid>
                        <Grid item xs={3}>
                            <Grid container spacing={0} alignItems="center">
                                <Grid item xs={6}><small className={classes.light}>CPU:</small></Grid>
                                <Grid item xs={6}>{node.cpu} <small>shares</small></Grid>
                                <Grid item xs={6}><small className={classes.light}>Memory:</small></Grid>
                                <Grid item xs={6}>{node.memory === 0 ? 0 : node.memory/(1024*1024)} <small>MB</small></Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </div>
                <Table className={classes.table}>
                    <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell align="center">CPU (shares)</TableCell>
                        <TableCell align="center">Memory (MB)</TableCell>
                        <TableCell align="center">Run Time</TableCell>
                        <TableCell align="center">Alloc Time</TableCell>
                    </TableRow>
                    </TableHead>
                    <TableBody>
                    {containers.map(row => {
                        return (
                            <TableRow key={row.ID}>
                                <TableCell component="th" scope="row">{row.Name}</TableCell>
                                <TableCell align="center">{row.CPUShares}<small className={classes.light}> ({row.PercentCPU}%)</small></TableCell>
                                <TableCell align="center">{row.Memory === 0 ? 0 : row.Memory/(1024*1024)}<small className={classes.light}> ({row.PercentMemory}%)</small></TableCell>
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