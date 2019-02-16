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
  tableCellNoWrap: {
    whiteSpace: 'nowrap',
    padding: theme.spacing.unit*2,
  },
  tableCell: {
    padding: theme.spacing.unit*2,
  },
  header: {
    padding: theme.spacing.unit*3,
  },
  light: {
    color: '#757575',
  }
});

function getLabels(list) {
    var labels = {};
    for (var i = 0; i<list.length; i++) {
        for (var k in list[i].Labels) {
            labels[k] = '';
        }
    }
    var llist = [];
    for (var l in labels) {
        llist.push(l);
    }
    llist.sort();
    return llist;
}

class ContainerList extends Component {
    state = {
        containers: [],
        labels: [],
    };

    componentDidMount() {   
        var {node} = this.props;
        axios.get('http://'+node.Address+'/container/')
        .then(resp => {
          var data = resp.data;
          for (var i = 0; i < data.length; i++) {
            data[i].PercentMemory = data[i].Memory > 0 ? ((data[i].Memory/node.Memory)*100).toFixed(0) : 100;
            data[i].PercentCPU = data[i].CPUShares > 0 ? ((data[i].CPUShares/node.CPUShares)*100).toFixed(0) : 100;
          }
          var labels = getLabels(data);
          this.setState({
              containers: resp.data,
              labels: labels,
            });
        });
    }
    
    render() {
        const { classes, node } = this.props;
        const { containers, labels } = this.state;

        return (
            <Paper className={classes.root}>
                <div className={classes.header}>
                    <Grid container spacing={0} alignItems="center">
                        <Grid item xs={9}>
                            <div>{node.Name}</div>
                            <div><small className={classes.light}>{node.Address}</small></div>
                        </Grid>
                        <Grid item xs={3}>
                            <Grid container spacing={0} alignItems="center">
                                <Grid item xs={6}><small className={classes.light}>CPU:</small></Grid>
                                <Grid item xs={6}>{node.CPUShares} <small>shares</small></Grid>
                                <Grid item xs={6}><small className={classes.light}>Memory:</small></Grid>
                                <Grid item xs={6}>{node.Memory === 0 ? 0 : Math.floor(node.Memory/(1024*1024))} <small>MB</small></Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </div>
                <Table className={classes.table}>
                    <TableHead>
                    <TableRow>
                        <TableCell className={classes.tableCell}>Name</TableCell>
                        <TableCell className={classes.tableCell} align="center">CPU (shares)</TableCell>
                        <TableCell className={classes.tableCell} align="center">Memory (MB)</TableCell>
                        <TableCell className={classes.tableCellNoWrap} align="center">Run Time</TableCell>
                        <TableCell className={classes.tableCellNoWrap} align="center">Alloc Time</TableCell>
                        {labels.map(label => {
                            return <TableCell key={label} className={classes.tableCellNoWrap}>{label}</TableCell>
                        })}
                    </TableRow>
                    </TableHead>
                    <TableBody>
                    {containers.map(row => {
                        return (
                            <TableRow key={row.ID}>
                                <TableCell className={classes.tableCell} component="th" scope="row">{row.Name}</TableCell>
                                <TableCell className={classes.tableCell} align="center">{row.CPUShares}<small className={classes.light}> ({row.PercentCPU}%)</small></TableCell>
                                <TableCell className={classes.tableCell} align="center">{row.Memory === 0 ? 0 : row.Memory/(1024*1024)}<small className={classes.light}> ({row.PercentMemory}%)</small></TableCell>
                                <TableCell className={classes.tableCellNoWrap} align="center"><TimeTicker start={row.Start} stop={row.Stop} /></TableCell>
                                <TableCell className={classes.tableCellNoWrap} align="center"><TimeTicker start={row.Create} stop={row.Destroy} /></TableCell>
                                {labels.map(label => {
                                    return <TableCell key={label} className={classes.tableCellNoWrap}>{row.Labels[label]}</TableCell>
                                })}
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