import React, { Component } from 'react';
import axios from 'axios';
import { ResponsiveContainer, BarChart, Bar } from 'recharts';
import { XAxis, YAxis, Tooltip, Legend } from 'recharts';
import { Typography } from '@material-ui/core';

const renderLegend = (props) => {
  const { payload } = props;
  return (
    <div>
      {payload.map((entry, index) => (
          <Typography key={index} variant="caption">{entry.value.toLowerCase()}</Typography>
      ))}
    </div>
  );
}

class BiaxialBarChart extends Component {

  state = {
    data: [],
    domainY: [],
  };

  componentDidMount() {
    axios.get(this.props.source)
    .then(resp => {
      var data = resp.data;
      for (var i = 0; i<data.length;i++) {
        data[i].Timestamp = (new Date(data[i].Timestamp/1e6)).toLocaleString();
      }

      var maxPrice = Math.max.apply(Math, data.map(function(o) { return o.Price; })),
          minPrice = Math.min.apply(Math, data.map(function(o) { return o.Price; })),
          domain = [minPrice, maxPrice];
      this.setState({data: resp.data, domainY:domain});
    });
  }

  render() {
    const { data, domainY } = this.state;

    return (
      <ResponsiveContainer width="100%" height={100}>
        <BarChart data={data}>
          <XAxis dataKey="Timestamp"/>
          <YAxis yAxisId="left" type="number" domain={domainY} orientation="left"/>
          <YAxis yAxisId="right" type="number" domain={domainY} orientation="right"/>
          <Tooltip />
          <Legend iconType="line" iconSize={12} content={renderLegend}/>
          <Bar yAxisId="left" dataKey="Price" fill="#0088FE" minPointSize={5}/>
        </BarChart>
      </ResponsiveContainer>
    );
  }
}



export default (BiaxialBarChart);
