import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import TimeTicker from './TimeTicker';

const styles = theme => ({
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
  light: {
    color: '#757575',
  }
});

const getLabels = (list) => {
    var labels = {};
    for (var i = 0; i<list.length; i++) {
        for (var k in list[i].Labels) {
            labels[k] = '';
        }
    }
    delete labels.name;

    var llist = [];
    for (var l in labels) {
        llist.push(l);
    }
    llist.sort();
    return llist;
}

function stableSort(array, cmp) {
    const stabilizedThis = array.map((el, index) => [el, index]);
    stabilizedThis.sort((a, b) => {
      const order = cmp(a[0], b[0]);
      if (order !== 0) return order;
      return a[1] - b[1];
    });
    return stabilizedThis.map(el => el[0]);
}

const desc = (a, b, orderBy) => {
    if (b[orderBy] < a[orderBy]) {
      return -1;
    }
    if (b[orderBy] > a[orderBy]) {
      return 1;
    }
    return 0;
}
  
const getSorting = (order, orderBy) => {
    return order === 'desc' ? (a, b) => desc(a, b, orderBy) : (a, b) => -desc(a, b, orderBy);
}

const memMB = (d) => {
    return Math.floor(d/(1024*1024))
}

class ContainersTable extends Component {
    state = {
        order: 'desc',
        orderBy: 'UnitsBurned',
    };

    handleRequestSort = (event, property) => {
        const orderBy = property;
        let order = 'desc';
    
        if (this.state.orderBy === property && this.state.order === 'desc') {
          order = 'asc';
        }
    
        this.setState({ order:order, orderBy:orderBy });
    }

    render() {
        const { classes, containers } = this.props;
        const { orderBy, order } = this.state;
        const labels = getLabels(containers);
        return (
            <Table className={classes.table}>
                <TableHead>
                    <TableRow>
                    {labels.map(label => {
                        return (
                            <TableCell key={label} className={classes.tableCellNoWrap} 
                                sortDirection={orderBy === label ? order : false}>
                                <TableSortLabel
                                    active={orderBy === label}
                                    direction={order}
                                    onClick={event => this.handleRequestSort(event, label)}
                                >
                                    {label}
                                </TableSortLabel>
                            </TableCell>
                        );
                    })}
                    <TableCell className={classes.tableCell} align="center" 
                        sortDirection={orderBy === 'CPUShares' ? order : false}>
                        <TableSortLabel
                            active={orderBy === 'CPUShares'}
                            direction={order}
                            onClick={event => this.handleRequestSort(event, 'CPUShares')}
                        >CPU (shares)</TableSortLabel>
                    </TableCell>
                    <TableCell className={classes.tableCell} align="center"
                        sortDirection={orderBy === 'Memory' ? order : false}>
                        <TableSortLabel
                            active={orderBy === 'Memory'}
                            direction={order}
                            onClick={event => this.handleRequestSort(event, 'Memory')}
                        >Memory (MB)</TableSortLabel>
                    </TableCell>
                    <TableCell className={classes.tableCellNoWrap} align="right"
                        sortDirection={orderBy === 'runtime' ? order : false}>
                        {/* <TableSortLabel
                            active={orderBy === 'runtime'}
                            direction={order}
                            onClick={event => this.handleRequestSort(event, 'runtime')}
                        > */}
                            Run Time
                        {/* </TableSortLabel> */}
                    </TableCell>
                    <TableCell className={classes.tableCellNoWrap} align="right"
                        sortDirection={orderBy === 'alloctime' ? order : false}>
                        {/* <TableSortLabel
                            active={orderBy === 'alloctime'}
                            direction={order}
                            onClick={event => this.handleRequestSort(event, 'alloctime')}
                        > */}
                            Alloc Time
                        {/* </TableSortLabel> */}
                    </TableCell>
                    <TableCell className={classes.tableCellNoWrap}
                        sortDirection={orderBy === 'UnitsBurned' ? order : false}>
                         <TableSortLabel
                            active={orderBy === 'UnitsBurned'}
                            direction={order}
                            onClick={event => this.handleRequestSort(event, 'UnitsBurned')}
                        >
                            Price ($)
                        </TableSortLabel>
                    </TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                {stableSort(containers, getSorting(order, orderBy)).map(row => {
                    return (
                        <TableRow key={row.ID}>
                            {/* <TableCell className={classes.tableCell} component="th" scope="row">{row.Name}</TableCell> */}
                            {labels.map(label => {
                                return <TableCell key={label} className={classes.tableCellNoWrap}>{row.Labels[label]}</TableCell>
                            })}
                            <TableCell className={classes.tableCell} align="center">
                                {row.CPUShares}<small className={classes.light}> ({row.CPUPercent}%)</small>
                            </TableCell>
                            <TableCell className={classes.tableCell} align="center">
                                {row.Memory === 0 ? 0 : memMB(row.Memory)}<small className={classes.light}> ({row.MemoryPercent}%)</small>
                            </TableCell>
                            <TableCell className={classes.tableCellNoWrap} align="right">
                                <TimeTicker start={row.Start} stop={row.Stop} />
                            </TableCell>
                            <TableCell className={classes.tableCellNoWrap} align="right">
                                <TimeTicker start={row.Create} stop={row.Destroy} />
                            </TableCell>
                            <TableCell className={classes.tableCellNoWrap}>{row.UnitsBurned}</TableCell>
                        </TableRow>
                    );
                })}
                </TableBody>
            </Table>
        );
    }
}

ContainersTable.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ContainersTable);