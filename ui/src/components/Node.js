import React, { Component } from 'react';
import axios from 'axios';

import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import { Paper, Grid, Divider, Chip } from '@material-ui/core';
import DonutChart from './DonutChart';
import BiaxialBarChart from './BiaxialBarChart';
import NodeHeader from './NodeHeader';
import ContainersTable from './ContainersTable';

const styles = theme => ({
  root: {
    width: '100%',
    overflowX: 'auto',
  },
  light: {
    color: '#757575',
  },
  tag: {
      margin: theme.spacing.unit/2,
      color: '#757575',
  },
});

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

const computePrice = (ph, start) => {
    var total = 0;
    for (var i = 0; i < ph.length; i++) {
        if (ph[i].Timestamp > start) {
            for (var j = i; j < ph.length; j++) {
                total += ph[j].Price;
            }
            return total;
        }
    }
    return total;
}

class Node extends Component {
    state = {
        containers: [],
        priceHistory: [],
        usedMem: 0,
        freeMem: 0,
        usedCPU: 0,
        freeCPU: 0,
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
                usedMem: usedMem,
                freeMem: freeMem,
                usedCPU: usedCPU,
                freeCPU: freeCPU,
            });
            
            this.fetchPrices(data);
        });
    }

    fetchPrices = (containers) => {
        const { node } = this.props;
        axios.get(`http://${node.Address}/price/`)
        .then(resp => {
            // Pricing
            var data = resp.data;
            for (var i = 0; i<data.length;i++) {
                data[i].Time = (new Date(data[i].Timestamp/1e6)).toLocaleString();
            }
            // Set containers price
            for (i=0;i<containers.length;i++) {
                containers[i].Price = computePrice(data, containers[i].Start);
            }

            this.setState({priceHistory: data, containers: containers});
        });
    }

    render() {
        const { classes, node } = this.props;
        const { containers } = this.state;
        const { priceHistory } = this.state;

        const tags = mapToList(node.Meta);
        const memData = donutChartData(this.state.usedMem,this.state.freeMem);
        const cpuData = donutChartData(this.state.usedCPU,this.state.freeCPU);

        return (
            <Paper className={classes.root}>
                <NodeHeader node={node} />
                <Divider/>
                <Grid container spacing={0} alignItems="center" alignContent="center" justify="space-evenly">
                    <Grid item xs={12} style={{textAlign:'center',paddingTop:10}}>
                    {tags.map(item => {
                        return (
                            <Chip label={item.key+': '+item.value} variant="outlined"
                                key={item.key} className={classes.tag}/>
                        );
                    })}
                    </Grid>
                    <Grid item xs={5}>
                        <DonutChart height={300} width={540} innerRadius={65} outerRadius={90} 
                            title="CPU" data={cpuData} unit="shares"/>
                    </Grid>
                    <Grid item xs={5}>
                        <DonutChart height={300} width={540} innerRadius={65} outerRadius={90} 
                            title="Memory" data={memData} unit="MB"/>
                    </Grid>
                    <Grid item xs={12} style={{textAlign:'center', paddingBottom: 10}}>
                        <BiaxialBarChart data={priceHistory} />
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