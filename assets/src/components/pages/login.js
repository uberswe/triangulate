import React from 'react'
import {Button, Col, Form, Row} from "react-bootstrap"
import {Link} from "react-router-dom";
import axios from "axios";

class Login extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            email: "",
            password: "",
        };

        this.login = this.login.bind(this)
        this.onInputChange = this.onInputChange.bind (this);
    }

    onInputChange(event) {
        this.setState ({
            [event.target.name]: event.target.value
        });
    }

    login() {
        const json = JSON.stringify({
            email: this.state.email,
            password: this.state.password,
        });
        axios.post ('/api/v1/login', json,
            {
                headers: {
                    "Content-type": "application/json",
                }
            },).then (result => {
                window.location = result.request.responseURL
        }).catch (error => {
            console.log (error)
        });
    }

    render() {

        return (<Row>
            <Col md={{span: 4, offset: 4}}>
                <h2>Login</h2>
                <Form>
                    <Form.Group controlId="formEmail">
                        <Form.Label>Email</Form.Label>
                        <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on" onChange={this.onInputChange} value={this.state.email}/>
                        <Form.Text>Forgot your email? Please contact <a
                            href="mailto:support@triangulate.xyz">support@triangulate.xyz</a></Form.Text>
                    </Form.Group>
                    <Form.Group controlId="formPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control name={`password`} type="password" id="formPassword" autocomplete="on" onChange={this.onInputChange} value={this.state.password}/>
                        <Form.Text><Link to="/forgot-password">Forgot your password?</Link></Form.Text>
                    </Form.Group>
                    <Form.Group>
                        <Button onClick={this.login} variant="primary">Login</Button>
                    </Form.Group>
                </Form>
            </Col>
        </Row>)
    }
}

export default Login