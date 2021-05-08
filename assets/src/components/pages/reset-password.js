import React from 'react'
import {Alert, Button, Col, Form, Row} from "react-bootstrap"
import {Link} from "react-router-dom";
import axios from "axios";
import PulseLoader from "react-spinners/PulseLoader";

class ResetPassword extends React.Component {
    constructor(props) {
        super (props);

        this.state = {
            code: "",
            email: "",
            password: "",
            isLoading: false,
            error: false,
            success: false,
        };

        this.resetPassword = this.resetPassword.bind (this)
        this.onInputChange = this.onInputChange.bind (this);
    }


    componentDidMount() {
        const { code } = this.props.match.params
        this.setState({
            code: code
        })
    }

    onInputChange(event) {
        this.setState ({
            [event.target.name]: event.target.value
        });
    }

    resetPassword() {
        this.setState ({
            isLoading: true
        })
        const json = JSON.stringify ({
            code: this.state.code,
            email: this.state.email,
            password: this.state.password,
        });
        axios.post ('/api/v1/reset-password', json,
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
            buttonText = "Reset password"
        }
        return (<Row>
            <Col md={{span: 4, offset: 4}}>
                <h2>Forgot Password</h2>
                <Form>
                    <Form.Text>Please enter your email and your new password.</Form.Text>
                    <Form.Group controlId="formEmail">
                        <Form.Label>Email</Form.Label>
                        <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on"  onChange={this.onInputChange} value={this.state.email}/>
                        <Form.Text>Forgot your email? Please contact <a
                            href="mailto:markus@triangulate.xyz?subject=Triangulate.xyz%3A%20Forgot%20Email&content=Please%20include%20the%20time%20when%20you%20purchased%20a%20premium%20account%20and%20any%20purchase%20reference%20you%20may%20have%20received%20from%20Stripe.">support@triangulate.xyz</a></Form.Text>
                    </Form.Group>
                    <Form.Group controlId="formPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control name={`password`} type="password" id="formPassword" autocomplete="on" onChange={this.onInputChange} value={this.state.password}/>
                    </Form.Group>
                    <Alert show={this.state.error} variant="danger">
                        An error occurred. Please make sure the email is correctly entered and that your password is at least 8 characters in length.
                    </Alert>
                    <Alert show={this.state.success} variant="success">
                        Your password has now been reset.
                    </Alert>
                    <Form.Group>
                        <Button disabled={this.state.isLoading} onClick={this.resetPassword} variant="primary">{buttonText}</Button>
                        <Form.Text><Link to="/login">Back to login</Link></Form.Text>
                    </Form.Group>
                </Form>
            </Col>
        </Row>)
    }
}

export default ResetPassword