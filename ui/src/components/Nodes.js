import React, { Component } from 'react';
import { Divider, Grid, Paper, Typography, withStyles } from '@material-ui/core';
import NodeHeader from './NodeHeader';
import DonutChart from './DonutChart';

const chartData = (data, field) => {
    return data.map(item => {
        return {name: item.Name.substr(0,8), value: item[field]};
    });
}

const sumField = (data, field) => {
    var total = 0;
    for (var i=0; i<data.length;i++) {
        total += data[i][field];
    }
    return total
}

const seriesToMB = (data) => {
    return data.map(item => {
        return {name:item.name, value: Math.floor(item.value/(1024*1024))};
    });
}

const styles = (theme) => ({
    infoLabel: {
        display: 'inline-block'
    }
})

class Nodes extends Component{
    render() {
        const {classes, data} = this.props;
        // const cpuData = chartData(data, 'CPUShares');
        // const memData = seriesToMB(chartData(data, 'Memory'));
        const cpuTotal = sumField(data, 'CPUShares');
        const memTotal = sumField(data, 'Memory');
        // const nodeCount = data.length;
        return (
            <div>
                <Grid container spacing={0} alignItems="center" justify="flex-end">
                    <Grid item xs={3}>
                        <Grid container spacing={0} alignItems="center">
                            <Grid item xs={6}><Typography>Nodes:</Typography></Grid>
                            <Grid item xs={6}><Typography>{data.length}</Typography></Grid>
                            <Grid item xs={6}><Typography>CPU:</Typography></Grid>
                            <Grid item xs={6}
                                ><Typography className={classes.infoLabel}>{cpuTotal} </Typography>
                                <Typography className={classes.infoLabel} variant="caption">shares</Typography>
                            </Grid>
                            <Grid item xs={6}><Typography>Memory:</Typography></Grid>
                            <Grid item xs={6}>
                                <Typography className={classes.infoLabel}>{(memTotal/(1024*1024*1024)).toFixed(2)}</Typography>
                                <Typography className={classes.infoLabel} variant="caption">GB</Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
                <Divider/>
                {/* <Grid container spacing={0} alignItems="center">
                    <Grid item xs={5} style={{textAlign:'center'}}>
                        <Typography>{cpuTotal} shares</Typography>
                        <DonutChart height={300} width={540} 
                            innerRadius={65} outerRadius={90} 
                            colors={['#70af97']} title="CPU" data={cpuData} unit="shares"
                        />      
                    </Grid>
                    <Grid item xs={2}>
                        <Paper>
                            <Typography>{nodeCount} nodes</Typography>
                        </Paper>
                    </Grid>
                    <Grid item xs={5} style={{textAlign:'center'}}>
                    <Typography>{memTotal} MB</Typography>
                        <DonutChart height={300} width={540} 
                            innerRadius={65} outerRadius={90} 
                            colors={['#b9a4a9']}
                            title="Memory" data={memData} unit="MB"
                        />
                    </Grid>
                </Grid> */}
                {data.map(node => {
                   return (
                        <div key={node.Name}>
                            <NodeHeader node={node} />
                            <Divider/>
                        </div>
                   );
                })}
            </div>
        );
    }
}

export default withStyles(styles)(Nodes);