import React from 'react'
import {Button, Col, Form, Row, Alert} from "react-bootstrap"
import {Link} from "react-router-dom";
import axios from "axios";
import { css } from "@emotion/core";
import PulseLoader from "react-spinners/PulseLoader";

class Login extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            email: "",
            password: "",
            isLoading: false,
            error: false,
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
        this.setState({
            isLoading: true
        })
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
            this.setState({
                isLoading: false,
                error:false,
            })
        }).catch (error => {
            console.log (error)
            this.setState({
                isLoading: false,
                error:true,
            })
        });
    }

    render() {
        let buttonText = "Login"
        if (this.state.isLoading) {
            buttonText = (<PulseLoader color="#f7f7f7" loading={this.state.isLoading} size={10}/>)
        }
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
                        <Alert show={this.state.error} variant="danger">
                            An error occurred. Please make sure that you have entered the correct password and the email that belongs to your account. If you are still having trouble please contact <a
                            href="mailto:support@triangulate.xyz?subject=Triangulate.xyz%3A%20Trouble%20logging%20in">support@triangulate.xyz</a>.
                        </Alert>
                        <Button disabled={this.state.isLoading} onClick={this.login} variant="primary">{buttonText}</Button>
                    </Form.Group>
                </Form>
            </Col>
        </Row>)
    }
}

export default Login