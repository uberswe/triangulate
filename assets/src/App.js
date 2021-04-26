import './App.css';
import React from "react";
import Title from "./components/Title";
import {Button, Col, Container, Form, Image, Row} from "react-bootstrap";
import axios from "axios"
import 'bootstrap/dist/css/bootstrap.min.css';

class App extends React.Component {
    constructor(props) {
        super (props);

        this.state = {
            queue: null,
            image: "",
            timer: null,
            identifier: "",
            isLoading: false,
            width: 1200,
            height: 675,
            shapes: true,
            imageType: "upload"
        };

        this.handleInputChange = this.handleInputChange.bind (this);
    }

    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;
        this.setState ({
            [name]: value
        });
    }

    poll() {
        this.state.timer = setInterval (() => {
            axios.get ('/api/v1/generate/' + this.state.identifier).then (result => {
                this.setState ({
                    queue: result.data.queue,
                    image: result.data.link
                });
                if (this.state.image !== "") {
                    clearInterval (this.state.timer)
                    this.setState ({
                        isLoading: false
                    });
                }
            }).catch (error => {
                console.log (error)
                clearInterval (this.state.timer)
                this.setState ({
                    isLoading: false
                });
            });
        }, 1000);
    }

    toggleButtonState = () => {
        this.setState ({
            isLoading: true
        });
        axios.post ('/api/v1/generate', {
            width: this.state.width,
            height: this.state.height,
            shapes: this.state.shapes,
            type: this.state.imageType
        }).then (result => {
            this.setState ({
                queue: result.data.queue,
                identifier: result.data.identifier,
            });
            this.poll ()
        }).catch (error => {
            console.log (error)
        });
    };

    render() {
        let queueHolder
        if (this.state.image !== "") {
            queueHolder = <Image src={this.state.image} fluid/>
        } else if (this.state.queue > 0) {
            queueHolder = <div>Rendering, you are {this.state.queue} in the queue.</div>
        } else if (this.state.queue === 0) {
            queueHolder = <div>Rendering, your image is being generated.</div>
        }
        return (
            <Container fluid>
                <Row>
                    <Col md={12}>
                        <Title/>
                    </Col>
                    <Col md={6}>
                        <Form>
                            <Form.Group>
                                <Form.Label column mb={12}>
                                    An image is used as a starting point
                                </Form.Label>
                                <div key={`inline-radio`} className="mb-12">
                                    <Form.Check name="imageType" inline label="Upload an image" type="radio"
                                                id={`inline-radio-1`} value={`upload`}
                                                onChange={this.handleInputChange}/>
                                    <Form.Check name="imageType" inline label="Use a random image" type="radio"
                                                id={`inline-radio-2`} value={'random'}
                                                onChange={this.handleInputChange}/>
                                </div>
                            </Form.Group>
                            <Form.Group>
                                <Form.Check checked={this.state.shapes} onChange={this.handleInputChange} name="shapes"
                                            label="Add Shapes" type="checkbox" id={`Shapes`}/>
                            </Form.Group>
                            <Form.Group controlId="sizeGroup">
                                <Form.Label>Width</Form.Label>
                                <Form.Control value={this.state.width} onChange={this.handleInputChange} name={`width`}
                                              type="Text"/>
                                <Form.Label>Height</Form.Label>
                                <Form.Control value={this.state.height} onChange={this.handleInputChange}
                                              name={`height`} type="Text"/>
                            </Form.Group>
                            // Edge Count of shapes
                            // Stroke shapes true/false
                            // Stroke size
                            // Blur amount
                            // Triangulate bool
                            // Shapes before or after triangulating?

                            // Generate an image
                            // Pick colors

                            <Form.Group>
                                <Button
                                    disabled={this.state.isLoading}
                                    onClick={this.toggleButtonState}
                                    variant="primary">Generate!</Button>
                            </Form.Group>
                        </Form>
                    </Col>
                </Row>
                <Row>
                    <Col md={12}>
                        {queueHolder}
                    </Col>
                </Row>
                <Row>
                    <Col md={12}>
                        This project uses ideas and code from <a href="https://github.com/esimov/triangle">github.com/esimov/triangle</a> and <a href="https://github.com/preslavrachev/generative-art-in-go">github.com/preslavrachev/generative-art-in-go</a>
                    </Col>
                </Row>
            </Container>
        );
    }

    componentWillUnmount() {
        this.timer = null;
    }
}

export default App;
