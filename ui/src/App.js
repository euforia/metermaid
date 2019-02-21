import React, { Component } from 'react';
import axios from 'axios';

import './App.css';
import Node from './components/Node';
import Nodes from './components/Nodes';

const API_HOST = process.env.REACT_APP_METERMAID_HOST || '';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
    };
  }

  componentDidMount() {
    axios.get(API_HOST+'/node/')
    .then(resp => {
        this.setState({nodes:resp.data});
    })
    .catch(err => {
      console.log(err);
    });
  }

  render() {
    const {nodes} = this.state;
    return (
      <div className="App">
        {nodes.map(node => {
          return <Node key={node.Name} node={node} />
        })}
      </div>
    );
  }
}

export default App;
