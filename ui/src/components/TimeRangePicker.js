import React, { Component } from 'react';

import { Paper, Grid, Divider, Chip, Typography, TextField, Button } from '@material-ui/core';

class TimeRangePicker extends Component{
    // state = {
    //     start: '',
    //     end: '',
    // }

    handleDateChange = (event, field) => {
        // this.setState({[field]:event.target.value});
        if (field ==='start') {
            if (this.props.onStartChange) this.props.onStartChange(event);
        } else {
            if (this.props.onEndChange) this.props.onEndChange(event);
        }
    }

    // handleButtonClick = (event) => {
    //     if (this.props.onRangeSet) 
    //         this.props.onRangeSet({
    //             start: this.state.start, 
    //             end: this.state.end,
    //         });
    // }

    render() {
        const {start, end} = this.props;

        return (
            <Grid container spacing={0} alignItems="center" justify="flex-end">
                    <Grid item xs={3} style={{textAlign:'right'}}>
                        <TextField
                            label="Start"
                            type="datetime-local"
                            value={start}
                            InputLabelProps={{
                                shrink: true,
                            }}
                            onChange={event => this.handleDateChange(event, 'start')}
                        />
                    </Grid>
                    <Grid item xs={3} style={{textAlign:'right'}}>
                        <TextField
                            label="End"
                            type="datetime-local"
                            value={end}
                            InputLabelProps={{
                                shrink: true,
                            }}
                            onChange={event => this.handleDateChange(event, 'end')}
                        />
                    </Grid>
                    <Grid item xs={1} style={{textAlign:'center'}}>
                        <Button onClick={this.props.onSetRange}>filter</Button>
                    </Grid>
            </Grid>
        );
    }
}

export default (TimeRangePicker);