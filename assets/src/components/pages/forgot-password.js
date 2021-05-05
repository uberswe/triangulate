import React from 'react'
import {Button, Col, Form, Row} from "react-bootstrap"
import {Link} from "react-router-dom";

function ForgotPassword() {
    return (<Row>
        <Col md={{span: 4, offset: 4}}>
            <h2>Forgot Password</h2>
            <Form>
                <Form.Text>An email with a reset link will be sent to you if an account with that email exists.</Form.Text>
                <Form.Group controlId="formEmail">
                    <Form.Label>Email</Form.Label>
                    <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on"/>
                    <Form.Text>Forgot your email? Please contact <a
                        href="mailto:markus@triangulate.xyz">support@triangulate.xyz</a></Form.Text>
                </Form.Group>
                <Form.Group>
                    <Button variant="primary">Request Password Reset Link</Button>
                    <Form.Text><Link to="/login">Back to login</Link></Form.Text>
                </Form.Group>
            </Form>
        </Col>
    </Row>)
}

export default ForgotPassword