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
            shapeStroke: true,
            triangulate: false,
            triangulateBefore: false,
            strokeThickness: 5,
            blurAmount: 5,
            min: 4,
            max: 4,
            imageType: "random",
            filesUpload: null
        };

        this.handleInputChange = this.handleInputChange.bind (this);
        this.onFileChange = this.onFileChange.bind (this)
    }

    onFileChange = event => {
        const files = Array.from (event.target.files)
        this.setState ({filesUpload: files});
        console.log (this.state.filesUpload)
    };

    handleInputChange(event) {
        const target = event.target;
        let value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;
        if (name === "max") {
            if (value < this.state.min) {
                value = this.state.min;
            }
        } else if (name === "min") {
            if (value > this.state.max) {
                value = this.state.max;
            }
        } else if (name === "width") {
            if (value > 1200) {
                value = 1200;
            } else if (value < 0) {
                value = 0;
            }
        } else if (name === "height") {
            if (value > 1200) {
                value = 1200;
            } else if (value < 0) {
                value = 0;
            }
        }
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
        let fdata = new FormData ();
        this.state.filesUpload.forEach(function(d){
            fdata.append('fileUpload', d);
        })
        let data = {
            width: this.state.width,
            height: this.state.height,
            shapes: this.state.shapes,
            type: this.state.imageType,
            shapeStroke: this.state.shapeStroke,
            triangulate: this.state.triangulate,
            triangulateBefore: this.state.triangulateBefore,
            strokeThickness: this.state.strokeThickness,
            blurAmount: this.state.blurAmount,
            min: this.state.min,
            max: this.state.max,
        }
        for (const [key, value] of Object.entries (data)) {
            fdata.append (key, value)
        }
        axios.post ('/api/v1/generate', fdata,
            {
                headers: {
                    "Content-type": "multipart/form-data",
                }
            },).then (result => {
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
        let queueHolder = ""
        if (this.state.image !== "") {
            queueHolder = <Image src={this.state.image} fluid/>
        } else if (this.state.queue > 0) {
            queueHolder = <div>Rendering, you are {this.state.queue} in the queue.</div>
        } else if (this.state.queue === 0) {
            queueHolder = <div>Rendering, your image is being generated.</div>
        }
        let shapeOptions = ""
        if (this.state.shapes) {
            shapeOptions = (<Row>
                    <Col md={12}>
                        <Form.Group>
                            <Form.Check checked={this.state.triangulateBefore} onChange={this.handleInputChange}
                                        name="triangulateBefore"
                                        label="Triangulate before shapes" type="checkbox" id={`TriangulateBefore`}/>
                        </Form.Group>
                        <Form.Group controlId="shapeVertexMin">
                            <Form.Label>Minimum number of vertices: {this.state.min}</Form.Label>
                            <Form.Control min="3" max="10" name="min" value={this.state.min}
                                          onChange={this.handleInputChange} type="range"/>
                        </Form.Group>
                        <Form.Group controlId="shapeVertexMax">
                            <Form.Label>Maximum number of vertices: {this.state.max}</Form.Label>
                            <Form.Control min="3" max="10" name="max" value={this.state.max}
                                          onChange={this.handleInputChange} type="range"/>
                        </Form.Group>
                        <Form.Group>
                            <Form.Check checked={this.state.shapeStroke} onChange={this.handleInputChange}
                                        name="shapeStroke"
                                        label="Add a stroke to shapes" type="checkbox" id={`Shapes`}/>
                        </Form.Group>
                        <Form.Group controlId="strokeThickness">
                            <Form.Label>Stroke thickness</Form.Label>
                            <Form.Control min="1" max="10" name="strokeThickness" value={this.state.strokeThickness}
                                          onChange={this.handleInputChange} type="range"/>
                        </Form.Group>
                    </Col>
                </Row>
            );
        }
        let imageUpload = ""
        if (this.state.imageType === "upload") {
            let imageUploadLabel = "Select a file"
            if (this.state.fileUpload === null) {
                imageUploadLabel = "Selected file"
            }
            imageUpload = (<Form.Group>
                <Form.Group>
                    <Form.File id="fileUpload" name="fileUpload" label={imageUploadLabel}
                               onChange={this.onFileChange}/>
                </Form.Group>
            </Form.Group>);
        }
        // Generate an image
        // Pick colors
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
                                                checked={this.state.imageType === "upload"}
                                                onChange={this.handleInputChange}/>
                                    <Form.Check name="imageType" inline label="Use a random image" type="radio"
                                                id={`inline-radio-2`} value={'random'}
                                                checked={this.state.imageType === "random"}
                                                onChange={this.handleInputChange}/>
                                </div>
                            </Form.Group>
                            {imageUpload}
                            <Row>
                                <Col md={6}>
                                    <Form.Group controlId="widthGroup">
                                        <Form.Label>Width</Form.Label>
                                        <Form.Control value={this.state.width} onChange={this.handleInputChange}
                                                      name={`width`}
                                                      type="Text"/>

                                    </Form.Group>
                                </Col>
                                <Col md={6}>
                                    <Form.Group controlId="heightGroup">
                                        <Form.Label>Height</Form.Label>
                                        <Form.Control value={this.state.height} onChange={this.handleInputChange}
                                                      name={`height`} type="Text"/>
                                    </Form.Group>
                                </Col>
                            </Row>
                            <Form.Group>
                                <Form.Check checked={this.state.triangulate} onChange={this.handleInputChange}
                                            name="triangulate"
                                            label="Triangulate (not implemented)" type="checkbox" id={`Triangulate`}/>
                            </Form.Group>
                            <Form.Group>
                                <Form.Check checked={this.state.shapes} onChange={this.handleInputChange} name="shapes"
                                            label="Add Shapes (Always)" type="checkbox" id={`Shapes`}/>
                            </Form.Group>
                            {shapeOptions}
                            <Form.Group controlId="BlurAmount">
                                <Form.Label>Blur amount</Form.Label>
                                <Form.Control min="1" max="10" name="blurAmount" value={this.state.blurAmount}
                                              onChange={this.handleInputChange} type="range"/>
                            </Form.Group>
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
                        This project uses ideas and code from <a
                        href="https://github.com/esimov/triangle">github.com/esimov/triangle</a> and <a
                        href="https://github.com/preslavrachev/generative-art-in-go">github.com/preslavrachev/generative-art-in-go</a>
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
