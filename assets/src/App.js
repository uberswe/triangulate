import './App.scss';
import React from "react";
import Title from "./components/Title";
import {Button, Col, Container, Form, Image, Row, Card} from "react-bootstrap";
import axios from "axios"

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
            complexityAmount: 50,
            min: 4,
            max: 4,
            imageType: "random",
            filesUpload: null,
            maxPoints: 2500,
            pointsThreshold: 20,
            sobelThreshold: 10,
            triangulateWireframe: false,
            triangulateNoise: false,
            triangulateGrayscale: false,
            randomImage: false,
            thumbnail: "",
            user_link: "",
            user_location: "",
            user_name: "",
            dots: ".",
            image_link: ""
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
            if (parseInt (value) < parseInt (this.state.min)) {
                value = this.state.min;
            }
        } else if (name === "min") {
            if (parseInt (value) > parseInt (this.state.max)) {
                value = this.state.max;
            }
        } else if (name === "width") {
            if (value > 2000) {
                value = 2000;
            } else if (value < 0) {
                value = 0;
            }
        } else if (name === "height") {
            if (value > 2000) {
                value = 2000;
            } else if (value < 0) {
                value = 0;
            }
        } else if (name === "triangulate") {
            if (!value && !this.state.shapes) {
                this.setState({
                    shapes: true
                })
            }
        } else if (name === "shapes") {
            if (!value && !this.state.triangulate) {
                this.setState({
                    triangulate: true
                })
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
                if (result.data.randomImage) {
                    this.setState({
                        randomImage: true,
                        description: result.data.description,
                        thumbnail: result.data.thumbnail,
                        user_link: result.data.user_link,
                        user_location: result.data.user_location,
                        user_name: result.data.user_name,
                        image_link: result.data.image_link
                    })
                }
                this.setState({
                    dots: this.state.dots + "."
                })
                if (this.state.image !== "") {
                    this.setState({
                        dots: "."
                    })
                    clearInterval (this.state.timer)
                    this.setState ({
                        isLoading: false
                    });
                }
            }).catch (error => {
                alert(error)
                clearInterval (this.state.timer)
                this.setState ({
                    isLoading: false,
                    dots: "."
                });
            });
        }, 1000);
    }

    toggleButtonState = () => {
        this.setState ({
            isLoading: true,
            randomImage: false,
        });
        let fdata = new FormData ();
        if (this.state.filesUpload != null) {
            this.state.filesUpload.forEach (function (d) {
                fdata.append ('fileUpload', d);
            })
        }
        let data = {
            width: this.state.width,
            height: this.state.height,
            shapes: this.state.shapes,
            type: this.state.imageType,
            shapeStroke: this.state.shapeStroke,
            triangulate: this.state.triangulate,
            triangulateBefore: this.state.triangulateBefore,
            strokeThickness: this.state.strokeThickness,
            complexityAmount: this.state.complexityAmount,
            min: this.state.min,
            max: this.state.max,
            maxPoints: this.state.maxPoints,
            pointsThreshold: this.state.pointsThreshold,
            sobelThreshold: this.state.sobelThreshold,
            triangulateWireframe: this.state.triangulateWireframe,
            triangulateNoise: this.state.triangulateNoise,
            triangulateGrayscale: this.state.triangulateGrayscale,
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
            queueHolder = <div>Rendering, your image is being generated{this.state.dots}</div>
        }
        let media = ""

        if (this.state.randomImage) {
            media = (<Card style={{ width: '100%' }}>
                <Row>
                    <Col md={6}>
                        <a href={this.state.image_link}><Card.Img style={{ height: '100%' }} src={this.state.thumbnail}/></a>
                    </Col>
                    <Col md={6}>
                        <Card.Body>
                            <Card.Title>Your random image</Card.Title>
                            <Card.Text>
                                {this.state.description}
                                <hr/>
                                Photo by <a href={this.state.user_link+"?utm_source=triangulate&utm_medium=referral"}>{this.state.user_name}</a> on <a href="https://unsplash.com/?utm_source=triangulate&utm_medium=referral">Unsplash</a>
                            </Card.Text>
                        </Card.Body>
                    </Col>
                </Row>
            </Card>)
        }

        let triangulateOptions, triangulateOptions2, triangulateOptions3, triangulateOptions4, triangulateOptions5, triangulateOptions6, triangulateOptions7, shapeOptions, shapeOptions2, shapeOptions3, shapeOptions4 = ""
        if (this.state.triangulate) {
            triangulateOptions = (<Col md={4}>
                <Form.Group controlId="sobelThreshold">
                    <Form.Label>Sobel Threshold</Form.Label>
                    <Form.Control min="5" max="20" name="sobelThreshold" value={this.state.sobelThreshold}
                                  onChange={this.handleInputChange} type="range"/>
                </Form.Group>
            </Col>);
            triangulateOptions2 = (<Col md={4}>
                <Form.Group controlId="pointsThreshold">
                    <Form.Label>Points Threshold</Form.Label>
                    <Form.Control min="10" max="30" name="pointsThreshold" value={this.state.pointsThreshold}
                                  onChange={this.handleInputChange} type="range"/>
                </Form.Group>
            </Col>);
            triangulateOptions3 = (<Col md={4}>
                <Form.Group controlId="maxPoints">
                    <Form.Label>Max Points</Form.Label>
                    <Form.Control min="500" max="5000" name="maxPoints" value={this.state.maxPoints}
                                  onChange={this.handleInputChange} type="range"/>
                </Form.Group>
            </Col>);
            triangulateOptions4 = (<Col md={4}>
                <Form.Group controlId="strokeWidth">
                    <Form.Label>Stroke Width</Form.Label>
                    <Form.Control min="1" max="10" name="strokeWidth" value={this.state.strokeWidth}
                                  onChange={this.handleInputChange} type="range"/>
                </Form.Group>
            </Col>);
            triangulateOptions5 = (<Col md={4}>
                <Form.Group>
                    <Form.Check checked={this.state.triangulateWireframe} onChange={this.handleInputChange}
                                name="triangulateWireframe"
                                label="Wireframe" type="checkbox" id={`wireframe`}/>
                </Form.Group>
            </Col>);
            triangulateOptions6 = (<Col md={4}>
                <Form.Group>
                    <Form.Check checked={this.state.triangulateNoise} onChange={this.handleInputChange}
                                name="triangulateNoise"
                                label="Noise" type="checkbox" id={`noise`}/>
                </Form.Group>
            </Col>);
            triangulateOptions7 = (<Col md={4}>
                <Form.Group>
                    <Form.Check checked={this.state.triangulateGrayscale} onChange={this.handleInputChange}
                                name="triangulateGrayscale"
                                label="Grayscale" type="checkbox" id={`grayscale`}/>
                </Form.Group>
            </Col>);
        }
        if (this.state.shapes) {
            shapeOptions = (<Col md={4}>
                        <Row>
                            <Col md={6}>
                                <Form.Group controlId="shapeVertexMin">
                                    <Form.Label>Min vertices: {this.state.min}</Form.Label>
                                    <Form.Control min="3" max="10" name="min" value={this.state.min}
                                                  onChange={this.handleInputChange} type="range"/>
                                </Form.Group>
                            </Col>
                            <Col md={6}>
                                <Form.Group controlId="shapeVertexMax">
                                    <Form.Label>Max vertices: {this.state.max}</Form.Label>
                                    <Form.Control min="3" max="10" name="max" value={this.state.max}
                                                  onChange={this.handleInputChange} type="range"/>
                                </Form.Group>
                            </Col>
                        </Row>
                    </Col>);
            shapeOptions2 = (<Col md={4}>
                        <Form.Group>
                            <Form.Check checked={this.state.shapeStroke} onChange={this.handleInputChange}
                                        name="shapeStroke"
                                        label="Add a stroke to shapes" type="checkbox" id={`Shapes`}/>
                        </Form.Group>
            </Col>);
            shapeOptions3 = (<Col md={4}>
                        <Form.Group controlId="strokeThickness">
                            <Form.Label>Stroke thickness</Form.Label>
                            <Form.Control min="1" max="10" name="strokeThickness" value={this.state.strokeThickness}
                                          onChange={this.handleInputChange} type="range"/>
                        </Form.Group>
            </Col>);
            if (this.state.triangulate) {
                shapeOptions4 = (<Form.Group>
                    <Form.Check checked={this.state.triangulateBefore} onChange={this.handleInputChange}
                                name="triangulateBefore"
                                label="Triangulate before shapes" type="checkbox" id={`TriangulateBefore`}/>
                </Form.Group>);
            }
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
            <Container>
                <Row>
                    <Col md={12}>
                        <Title/>
                    </Col>
                    <Col md={12}>
                        <Form>
                            <Row>
                                <Col md={6}>
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
                                </Col>
                                <Col md={6}>
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
                                                <Form.Control value={this.state.height}
                                                              onChange={this.handleInputChange}
                                                              name={`height`} type="Text"/>
                                            </Form.Group>
                                        </Col>
                                    </Row>
                                </Col>
                            </Row>
                            <Row>
                                <Col md={4}>
                            <Form.Group>
                                <Form.Check checked={this.state.triangulate} onChange={this.handleInputChange}
                                            name="triangulate"
                                            label="Triangulate" type="checkbox" id={`Triangulate`}/>
                            </Form.Group>
                                </Col>
                                <Col md={4}>
                            <Form.Group>
                                <Form.Check checked={this.state.shapes} onChange={this.handleInputChange} name="shapes"
                                            label="Add Shapes" type="checkbox" id={`Shapes`}/>
                            </Form.Group>
                                </Col>
                                {triangulateOptions5}
                                {triangulateOptions6}
                                {triangulateOptions7}
                                {shapeOptions4}
                                {shapeOptions2}
                                <Col md={4}>
                                    <Form.Group controlId="ComplexityAmount">
                                        <Form.Label>Complexity</Form.Label>
                                        <Form.Control min="1" max="100" name="complexityAmount" value={this.state.complexityAmount}
                                                      onChange={this.handleInputChange} type="range"/>
                                    </Form.Group>
                                </Col>
                                {shapeOptions}
                                {shapeOptions3}
                                {triangulateOptions}
                                {triangulateOptions2}
                                {triangulateOptions3}
                                {triangulateOptions4}
                            </Row>
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
                <Row style={{ 'margin-top':'20px' }}>
                    <Col md={12}>
                        {media}
                    </Col>
                </Row>
                <Row style={{ 'margin-bottom':'50px' }}>
                    <Col md={12}>
                        <hr/>
                        Created by <a href="https://github.com/uberswe">Markus Tenghamn</a> | Contact <a href="mailto:markus@triangulate.xyz">markus@triangulate.xyz</a>
                        <hr/>
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
