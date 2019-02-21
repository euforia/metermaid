import React, { Component } from 'react';
import axios from 'axios';

import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import { Paper, Grid, Divider, Chip, Typography, TextField, Button } from '@material-ui/core';
import DonutChart from './DonutChart';
import BiaxialBarChart from './BiaxialBarChart';
import NodeHeader from './NodeHeader';
import ContainersTable from './ContainersTable';
import TimeRangePicker from './TimeRangePicker';

const styles = theme => ({
  root: {
    width: '100%',
    overflowX: 'auto',
  },
  light: {
    color: '#757575',
  },
//   tag: {
//       margin: theme.spacing.unit/4,
//       color: '#757575',
//       height: 26,
//       fontSize: 12,
//   },
  costBoard: {
      padding: theme.spacing.unit*3,
      margin: theme.spacing.unit,
      textAlign: 'center',
  }
});


const toHHMMSS = (msec_num, fix) => {
    // var sec_num = parseInt(this, 10); // don't forget the second param
    const sec_num = msec_num/1000;
    var hours   = Math.floor(sec_num / 3600);
    var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
    var seconds = (sec_num - (hours * 3600) - (minutes * 60)).toFixed(fix);

    var days = -1;
    if (hours   < 10) {
        hours   = "0"+hours;
    } else if (hours>23) {
        days = Math.floor(hours/24);
        hours = hours % 24;
    }

    if (minutes < 10) {minutes = "0"+minutes;}
    if (seconds < 10) {seconds = "0"+seconds;}

    if (days>-1) {
        return days+'d '+hours+'h '+minutes+'m '+seconds+'s';
    }
    return hours+'h '+minutes+'m '+seconds+'s';
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

const formatTime = (d) => {
    var month = d.getMonth()+1;
    if (month<10) month = "0"+month;
    var day = d.getDate()
    if (day<10) day = "0"+day;
    var hour = d.getHours()
    if (hour<10) hour = "0"+hour;
    var mins = d.getMinutes()
    if (mins<10) mins = "0"+mins;

    return d.getFullYear()+'-'+month+'-'+day+'T'+hour+':'+mins;
}

const memMB = (d) => {
    return Math.floor(d/(1024*1024))
}

// const mapToList = (m) => {
//     var out = [];
//     for (var k in m) {
//         out.push({key:k,value:m[k]});
//     }
//     return out;
// } 

class Node extends Component {
    state = {
        containers: [],
        priceHistory: [],
        currNodeCost: 0,
        timeWindow: '',
        usedMem: 0,
        freeMem: 0,
        usedCPU: 0,
        freeCPU: 0,
        endTime: formatTime(new Date()),
        startTime: '',
    };

    componentDidMount() {
        this.fetchContainers();
        // this.fetchPrices();
    }
    
    fetchContainers() {   
        const { node } = this.props;
        axios.get(`http://${node.Address}/container/`)
        .then(resp => {
          var data = resp.data;
          for (var i = 0; i < data.length; i++) {
            data[i].MemoryPercent = data[i].Memory > 0 ? ((data[i].Memory/node.Memory)*100).toFixed(0) : 100;
            data[i].CPUPercent = data[i].CPUShares > 0 ? ((data[i].CPUShares/node.CPUShares)*100).toFixed(0) : 100;
          }
          
          var usedMem = computeUsed(data,'Memory'),
              freeMem = node.Memory-usedMem;
          usedMem = memMB(usedMem);
          freeMem = memMB(freeMem);

          var usedCPU = computeUsed(data,'CPUShares'),
              freeCPU = node.CPUShares-usedCPU;
          
            this.setState({
                containers: data,
                usedMem: usedMem,
                freeMem: freeMem,
                usedCPU: usedCPU,
                freeCPU: freeCPU,
            });
            
            this.fetchPrices();
        });
    }

    fetchPrices = () => {
        const { node } = this.props;
        var {startTime, endTime} = this.state;

        var url = `http://${node.Address}/price/?end=${endTime}:00-08:00`;
        if (startTime !=='') url += `&start=${startTime}:00-08:00`;
        
        axios.get(url)
        .then(resp => {
            // Pricing
            var data = resp.data.History;
            for (var i = 0; i<data.length;i++) {
                data[i].Time = (new Date(data[i].Timestamp/1e6)).toLocaleString();
            }

            if (startTime === '') {
                startTime = formatTime(new Date(data[0].Timestamp/1e6));
            }

            this.setState({
                priceHistory: data, 
                currNodeCost:resp.data.Total,
                timeWindow: toHHMMSS((data[data.length-1].Timestamp-data[0].Timestamp)/1e6,0),
                startTime: startTime,
                // endTime: endTime,
            });
        });
    }

    handleStartDateChange = (event) => {
        this.setState({startTime:event.target.value});
    }
    handleEndDateChange = (event) => {
        this.setState({endTime:event.target.value});
    }

    render() {
        const { classes, node } = this.props;
        const { containers } = this.state;
        const { priceHistory, currNodeCost, timeWindow } = this.state;
        const { startTime, endTime } = this.state;

        // const tags = mapToList(node.Meta);
        const memData = donutChartData(this.state.usedMem,this.state.freeMem);
        const cpuData = donutChartData(this.state.usedCPU,this.state.freeCPU);

        return (
            <Paper className={classes.root}>
                <NodeHeader node={node} />
                <Divider/>
                <Grid container spacing={0} alignItems="center" alignContent="center" justify="space-evenly">
                    <Grid item xs={12} style={{paddingTop: 20}}>
                        <TimeRangePicker start={startTime} end={endTime}
                            onStartChange={this.handleStartDateChange} 
                            onEndChange={this.handleEndDateChange}
                            onSetRange={this.fetchPrices}
                        />
                    </Grid>
                    {/* <Grid item xs={5} style={{textAlign: 'center'}}>
                    </Grid>
                    <Grid item xs={3} style={{textAlign:'right'}}>
                        <TextField
                            label="Start"
                            type="datetime-local"
                            value={startTime}
                            InputLabelProps={{
                                shrink: true,
                            }}
                            onChange={event => this.handleDateChange(event, 'startTime')}
                        />
                    </Grid>
                    <Grid item xs={3} style={{textAlign:'right'}}>
                        <TextField
                            label="End"
                            type="datetime-local"
                            value={endTime}
                            InputLabelProps={{
                                shrink: true,
                            }}
                            onChange={event => this.handleDateChange(event, 'endTime')}
                        />
                    </Grid>
                    <Grid item xs={1} style={{textAlign:'center'}}>
                        <Button onClick={event => {this.fetchPrices()}}>filter</Button>
                    </Grid> */}
                    <Grid item xs={5}>
                        <DonutChart height={300} width={540} innerRadius={65} outerRadius={90} 
                            title="CPU" data={cpuData} unit="shares"/>
                    </Grid>
                    <Grid item xs={2}>
                        <Paper className={classes.costBoard}>
                            <Typography variant="h3">{currNodeCost.toFixed(2)}</Typography>
                            <Typography variant="body2">{timeWindow}</Typography>
                        </Paper>
                    </Grid>
                    <Grid item xs={5}>
                        <DonutChart height={300} width={540} innerRadius={65} outerRadius={90} 
                            title="Memory" data={memData} unit="MB"/>
                    </Grid>
                    <Grid item xs={12} style={{textAlign:'center', paddingBottom: 10}}>
                        <BiaxialBarChart data={priceHistory} keyY="Value"/>
                    </Grid>
                </Grid>
                <ContainersTable containers={containers} pricing={priceHistory}/>
            </Paper>
        );
    }
}

Node.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Node);