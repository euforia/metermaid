import React, { Component } from 'react';


 function toHHMMSS(msec_num, fix) {
    // var sec_num = parseInt(this, 10); // don't forget the second param
    const sec_num = msec_num/1000;
    var hours   = Math.floor(sec_num / 3600);
    var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
    var seconds = (sec_num - (hours * 3600) - (minutes * 60)).toFixed(fix);

    if (hours   < 10) {hours   = "0"+hours;}
    if (minutes < 10) {minutes = "0"+minutes;}
    if (seconds < 10) {seconds = "0"+seconds;}
    return hours+'h '+minutes+'m '+seconds+'s';
}


// <TimerExample start={Date.now()} />
class TimeTicker extends Component {        
    
    state = {
        elapsed: 0,
    };

    componentDidMount() {
        if (this.props.stop === 0) {
            this.timer = setInterval(() => {
                this.setState({elapsed: new Date() - new Date(this.props.start/1000000)});
            }, 1000);
        } else {
            this.setState({elapsed: (this.props.stop-this.props.start)/1000000});
        }
    }

    componentWillUnmount() {
        if (this.props.stop === 0) clearInterval(this.timer);
    }

    render() {
        // Calculate elapsed to tenth of a second:
        // const elapsed = Math.round(this.state.elapsed / 100);
        // This will give a number with one digit after the decimal dot (xx.x):
        // const seconds = (elapsed / 10).toFixed(0);    

        // Although we return an entire <p> element, react will smartly update
        // only the changed parts, which contain the seconds variable.
        // return <p>This example was started <b>{seconds} seconds</b> ago.</p>;
        return (
            <span>{toHHMMSS(this.state.elapsed, this.props.precision)}</span>
        );
    }
}

export default (TimeTicker);