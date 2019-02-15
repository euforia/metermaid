import React, { Component } from 'react';
// import axios from 'axios';
import './App.css';
import ContainerList from './components/ContainerList';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      containers: [],
    };
  }

  // componentDidMount() {
  //   axios.get('http://localhost:8080/container/')
  //   .then(resp => {
  //     var data = resp.data;
  //     for (var i = 0; i<data.length; i++) {
  //       // if (data[i].Stop === 0) data[i].Stop = data[i].Start;
  //       // if (data[i].Destroy === 0) data[i].Destroy = data[i].Create;
  //     }
  //     this.setState({containers: data});
  //   });
  // }

  render() {
    return (
      <div className="App">
        {/* <header className="App-header">
        </header> */}
        <ContainerList node={{name: 'localhost', cpu: 200, memory: 256, address: 'http://localhost:8080/container/'}} />
      </div>
    );
  }
}

export default App;
