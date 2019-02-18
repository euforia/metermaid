import React, { Component } from 'react';
import axios from 'axios';

import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import Paper from '@material-ui/core/Paper';
import TimeTicker from './TimeTicker';
import { Grid, Typography, Divider, Chip } from '@material-ui/core';
import DonutChart from './DonutChart';
import BiaxialBarChart from './BiaxialBarChart';

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
  },
  tag: {
      margin: theme.spacing.unit/2,
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

function desc(a, b, orderBy) {
    if (b[orderBy] < a[orderBy]) {
      return -1;
    }
    if (b[orderBy] > a[orderBy]) {
      return 1;
    }
    return 0;
}
  
function getSorting(order, orderBy) {
    return order === 'desc' ? (a, b) => desc(a, b, orderBy) : (a, b) => -desc(a, b, orderBy);
}

const donutChartData = (usedMem, freeMem) => {
    if (freeMem>=0) 
        return [{name: 'Used',value: usedMem},{name: 'Free',value: freeMem}];
    return [{name: 'Used',value: usedMem, color:'#ff4500'},{name: 'Over',value: freeMem, color: '#cc0000'}];
}

const computeUsed = (data, property) => {
    var used = 0;
    for(var i = 0; i<data.length; i++) {
        used += data[i][property];
    }
    return used;
}

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

class ContainerList extends Component {
    state = {
        containers: [],
        labels: [],
        order: 'asc',
        orderBy: 'Name',
        usedMem: 0,
        freeMem: 0,
        usedCPU: 0,
        freeCPU: 0,
    };

    componentDidMount() {   
        var {node} = this.props;
        axios.get('http://'+node.Address+'/container/')
        .then(resp => {
          var data = resp.data;
          for (var i = 0; i < data.length; i++) {
            data[i].MemoryPercent = data[i].Memory > 0 ? ((data[i].Memory/node.Memory)*100).toFixed(0) : 100;
            data[i].CPUPercent = data[i].CPUShares > 0 ? ((data[i].CPUShares/node.CPUShares)*100).toFixed(0) : 100;
          }
          var labels = getLabels(data);
          
          var usedMem = computeUsed(data,'Memory'),
              freeMem = node.Memory-usedMem;
          usedMem = memMB(usedMem);
          freeMem = memMB(freeMem);

          var usedCPU = computeUsed(data,'CPUShares'),
              freeCPU = node.CPUShares-usedCPU;
          this.setState({
              containers: resp.data,
              labels: labels,
              usedMem: usedMem,
              freeMem: freeMem,
              usedCPU: usedCPU,
              freeCPU: freeCPU,
            });
        });
    }

    handleRequestSort = (event, property) => {
        const orderBy = property;
        let order = 'desc';
    
        if (this.state.orderBy === property && this.state.order === 'desc') {
          order = 'asc';
        }
    
        this.setState({ order:order, orderBy:orderBy });
    }

    render() {
        const { classes, node } = this.props;
        const tags = mapToList(node.Meta);
        const { containers, labels, orderBy, order } = this.state;
        
        const memData = donutChartData(this.state.usedMem,this.state.freeMem);
        const cpuData = donutChartData(this.state.usedCPU,this.state.freeCPU);

        return (
            <Paper className={classes.root}>
                <div className={classes.header}>
                    <Grid container spacing={0} alignItems="center">
                        <Grid item xs={5}>
                            <Typography variant="subtitle1">{node.Name}</Typography>
                            <div><small className={classes.light}>{node.Address}</small></div>
                        </Grid>
                        <Grid item xs={4}></Grid>
                        <Grid item xs={3}>
                            <Grid container spacing={0} alignItems="center">
                                <Grid item xs={6}><small className={classes.light}>Platform:</small></Grid>
                                <Grid item xs={6}>{node.Platform.Name} <small>{node.Platform.Version}</small></Grid>
                                <Grid item xs={6}><small className={classes.light}>CPU:</small></Grid>
                                <Grid item xs={6}>{node.CPUShares} <small>shares</small></Grid>
                                <Grid item xs={6}><small className={classes.light}>Memory:</small></Grid>
                                <Grid item xs={6}>{node.Memory === 0 ? 0 : memMB(node.Memory)} <small>MB</small></Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </div>
                <Divider/>
                <Grid container spacing={0} alignItems="center" alignContent="center">
                    <Grid item xs={12} style={{textAlign:'center',paddingTop:10}}>
                    {tags.map(item => {
                        return (
                            <Chip label={item.key+': '+item.value} variant="outlined"
                                key={item.key} className={classes.tag}/>
                        );
                    })}
                    </Grid>
                    <Grid item xs={6}>
                        <DonutChart height={300} width={540} innerRadius={65} outerRadius={90} 
                            title="CPU" data={cpuData} unit="shares"/>
                    </Grid>
                    <Grid item xs={6}>
                        <DonutChart height={300} width={540} innerRadius={65} outerRadius={90} 
                            title="Memory" data={memData} unit="MB"/>
                    </Grid>
                    <Grid item xs={12} style={{textAlign:'center', paddingBottom: 10}}>
                        <BiaxialBarChart source={`http://${node.Address}/price/`} />
                    </Grid>
                </Grid>
                <Divider/>
                <Table className={classes.table}>
                    <TableHead>
                      <TableRow>
                        {/* <TableCell className={classes.tableCell}
                            sortDirection={orderBy === 'Name' ? order : false}>
                            <TableSortLabel
                                active={orderBy === 'Name'}
                                direction={order}
                                onClick={event => this.handleRequestSort(event, 'Name')}
                            >Name</TableSortLabel>
                        </TableCell> */}
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