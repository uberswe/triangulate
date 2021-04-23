import './App.css';
import React from "react";
import Image from 'react-bootstrap/Image'
import Title from "./components/Title";
import Button from "react-bootstrap/Button";
import axios from "axios"

class App extends React.Component {
    state = {
        queue: null,
        image: "",
        timer: null,
        identifier: ""
    };

    poll() {
        this.state.timer = setInterval(()=> {
            axios.get('/api/v1/generate/' + this.state.identifier).then(result => {
                console.log(result)
                this.setState({
                    queue: result.data.queue,
                    image: result.data.link
                });
                if (this.state.image !== "") {
                    clearInterval(this.state.timer)
                }
            }).catch(error => {
                console.log(error)
                clearInterval(this.state.timer)
            });
        }, 1000);
    }

    toggleButtonState = () => {
        axios.get('/api/v1/generate').then(result => {
            console.log(result)
            this.setState({
                queue: result.data.queue,
                identifier: result.data.identifier,
            });
            this.poll()
        }).catch(error => {
            console.log(error)
        });
    };

    render() {
        let queueHolder
        if (this.state.image !== "") {
            queueHolder = <Image src={this.state.image} fluid/>
        } else if (this.state.queue > 0) {
            queueHolder = <div>Rendering, you are {this.state.queue} in the queue.</div>
        }
        return (<div className="App">
                <Title/>
                {queueHolder}
                <Button onClick={this.toggleButtonState} variant="primary">Generate some art!</Button>{' '}
            </div>
        );
    }

    componentWillUnmount() {
        this.timer = null;
    }
}

export default App;
