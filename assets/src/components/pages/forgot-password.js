import React from 'react'
import {Alert, Button, Col, Form, Row} from "react-bootstrap"
import {Link} from "react-router-dom";
import axios from "axios";
import PulseLoader from "react-spinners/PulseLoader";

class ForgotPassword extends React.Component {
    constructor(props) {
        super (props);

        this.state = {
            email: "",
            isLoading: false,
            error: false,
            success: false,
        };

        this.forgotPassword = this.forgotPassword.bind (this)
        this.onInputChange = this.onInputChange.bind (this);
    }

    onInputChange(event) {
        this.setState ({
            [event.target.name]: event.target.value
        });
    }

    forgotPassword() {
        this.setState ({
            isLoading: true
        })
        const json = JSON.stringify ({
            email: this.state.email,
        });
        axios.post ('/api/v1/forgot-password', json,
            {
                headers: {
                    "Content-type": "application/json",
                }
            },).then (result => {
            this.setState ({
                isLoading: false,
                error: false,
                success: true,
            })
        }).catch (error => {
            console.log (error)
            this.setState ({
                isLoading: false,
                error: true,
                success: false,
            })
        });
    }


    render() {
        let buttonText = ""
        if (this.state.isLoading) {
            buttonText = (<PulseLoader color="#f7f7f7" loading={this.state.isLoading} size={10}/>)
        } else {
            buttonText = "Request Password Reset Link"
        }
        return (<Row>
            <Col md={{span: 4, offset: 4}}>
                <h2>Forgot Password</h2>
                <Form>
                    <Form.Text>An email with a reset link will be sent to you if an account with that email
                        exists.</Form.Text>
                    <Form.Group controlId="formEmail">
                        <Form.Label>Email</Form.Label>
                        <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on" onChange={this.onInputChange} value={this.state.email}/>
                        <Form.Text>Forgot your email? Please contact <a
                            href="mailto:markus@triangulate.xyz?subject=Triangulate.xyz%3A%20Forgot%20Email&content=Please%20include%20the%20time%20when%20you%20purchased%20a%20premium%20account%20and%20any%20purchase%20reference%20you%20may%20have%20received%20from%20Stripe.">support@triangulate.xyz</a></Form.Text>
                    </Form.Group>
                    <Alert show={this.state.error} variant="danger">
                        An error occurred. Please try again. If you are still having trouble please contact <a
                        href="mailto:support@triangulate.xyz?subject=Triangulate.xyz%3A%20Trouble%20recovering%20password">support@triangulate.xyz</a>.
                    </Alert>
                    <Alert show={this.state.success} variant="success">
                        Thank you! If you have an account with that email registered with us we will send you an email with instructions on how to reset your password shortly.
                    </Alert>
                    <Form.Group>
                        <Button disabled={this.state.isLoading} onClick={this.forgotPassword} variant="primary">{buttonText}</Button>
                        <Form.Text><Link to="/login">Back to login</Link></Form.Text>
                    </Form.Group>
                </Form>
            </Col>
        </Row>)
    }
}

export default ForgotPassword