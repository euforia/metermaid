import React, { Component } from 'react';
import { ResponsiveContainer, BarChart, Bar } from 'recharts';
import PropTypes from 'prop-types';
import { XAxis, YAxis, Tooltip, Legend } from 'recharts';
import { Typography } from '@material-ui/core';

class BiaxialBarChart extends Component {

  state = {
    domainY: [],
  };

  componentDidMount() {
    const { data } = this.props;
    var maxPrice = Math.max.apply(Math, data.map(function(o) { return o.Price; })),
        minPrice = Math.min.apply(Math, data.map(function(o) { return o.Price; }));
    this.setState({domainY: [minPrice, maxPrice]});
  }

  renderLegend = (props) => {
    const { payload } = props;
    return (
      <div>
        {payload.map((entry, index) => (
            <Typography key={index} variant="caption">{entry.value.toLowerCase()}</Typography>
        ))}
      </div>
    );
  }

  render() {
    const { data, keyY } = this.props;
    const { domainY } = this.state;

    return (
      <ResponsiveContainer width="100%" height={100}>
        <BarChart data={data}>
          <XAxis dataKey="Time"/>
          <YAxis yAxisId="left" type="number" domain={domainY} orientation="left"/>
          <YAxis yAxisId="right" type="number" domain={domainY} orientation="right"/>
          <Tooltip />
          <Legend iconType="line" iconSize={12} content={this.renderLegend}/>
          <Bar yAxisId="left" dataKey={keyY} fill="#0088FE" minPointSize={5}/>
        </BarChart>
      </ResponsiveContainer>
    );
  }
}

BiaxialBarChart.propTypes = {
  keyY: PropTypes.string.isRequired,
};

export default (BiaxialBarChart);
