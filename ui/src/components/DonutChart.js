import React, { Component } from 'react';

import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { PieChart, Pie, Sector, Cell } from 'recharts';

const styles = theme => ({
  chartContainer: {
      margin: 'auto'
  }  
});
//'#FFBB28',
const COLORS = ['#0088FE','#32cd32','#00C49F','#FF8042'];

class DonutChart extends Component{
    state = {
        activeIndex: 0,
    };

    renderActiveShape = (props) => {
        const RADIAN = Math.PI / 180;
        const {
          cx, cy, 
          innerRadius, outerRadius,
          midAngle, startAngle, endAngle,
          fill, payload, percent, value,
        } = props;
        const {title, unit} = this.props; 

        const sin = Math.sin(-RADIAN * midAngle);
        const cos = Math.cos(-RADIAN * midAngle);
        const sx = cx + (outerRadius + 10) * cos;
        const sy = cy + (outerRadius + 10) * sin;
        const mx = cx + (outerRadius + 30) * cos;
        const my = cy + (outerRadius + 30) * sin;
        const ex = mx + (cos >= 0 ? 1 : -1) * 22;
        const ey = my;
        const textAnchor = cos >= 0 ? 'start' : 'end';
    
        return (
          <g>
            <text x={cx} y={cy} dy={8} textAnchor="middle" fill="#444">{title}</text>
            <Sector
              cx={cx}
              cy={cy}
              innerRadius={innerRadius}
              outerRadius={outerRadius}
              startAngle={startAngle}
              endAngle={endAngle}
              fill={fill}
            />
            <Sector
              cx={cx}
              cy={cy}
              startAngle={startAngle}
              endAngle={endAngle}
              innerRadius={outerRadius + 6}
              outerRadius={outerRadius + 10}
              fill={fill}
            />
            <path d={`M${sx},${sy}L${mx},${my}L${ex},${ey}`} stroke={fill} fill="none" />
            <circle cx={ex} cy={ey} r={2} fill={fill} stroke="none" />
            <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} textAnchor={textAnchor} fill={fill} fontWeight="bold" fontSize={12}>{payload.name}</text>
            <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} dy={18} textAnchor={textAnchor} fill={fill} fontSize={12}>{`${value} ${unit}`}</text>
            <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} dy={36} textAnchor={textAnchor} fill="#999" fontSize={12}>
              {`(${(percent * 100).toFixed(2)}%)`}
            </text>
          </g>
        );
    };

    onPieEnter = (data, index) => {
        this.setState({
          activeIndex: index,
        });
    };

    render() {
        const {height, width} = this.props;
        const {innerRadius, outerRadius} = this.props;
        const {classes, data} = this.props;
        const {activeIndex} = this.state;
        const colors = this.props.colors ? this.props.colors : COLORS;

        return (
         <PieChart height={height} width={width} className={classes.chartContainer}>
            <Pie 
                dataKey="value"
                activeIndex={activeIndex}
                activeShape={this.renderActiveShape} 
                data={data} 
                cx={(width/2)} 
                cy={(height/2)}
                innerRadius={innerRadius}
                outerRadius={outerRadius} 
                // fill="#8884d8"
                onMouseEnter={this.onPieEnter}
            >
            {
                data.map((entry, index) =>
                    <Cell key={index} fill={entry.color ? entry.color : colors[index % colors.length]} />
                )
            }
            </Pie>
        </PieChart>
        );
    }
}

DonutChart.propTypes = {
    data: PropTypes.array.isRequired,
};

export default withStyles(styles)(DonutChart);