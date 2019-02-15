import React, { Component } from 'react';
import axios from 'axios';

import './App.css';
import ContainerList from './components/ContainerList';

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
        var data = resp.data;
        for (var i = 0; i < data.length; i++) {
          data[i].URL = `http://${data[i].Addr}:${data[i].Port}/container/`;
        }
        this.setState({nodes:data});
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
          return <ContainerList key={node.Name} source={node.URL} />
        })}
      </div>
    );
  }
}

export default App;
