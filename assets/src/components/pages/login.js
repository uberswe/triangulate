import React from 'react'
import {Button, Col, Form, Row} from "react-bootstrap"
import {Link} from "react-router-dom";

function Login() {
    return (<Row>
        <Col md={{span: 4, offset: 4}}>
            <h2>Login</h2>
            <Form>
                <Form.Group controlId="formEmail">
                    <Form.Label>Email</Form.Label>
                    <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on"/>
                    <Form.Text>Forgot your email? Please contact <a
                        href="mailto:markus@triangulate.xyz">support@triangulate.xyz</a></Form.Text>
                </Form.Group>
                <Form.Group controlId="formPassword">
                    <Form.Label>Password</Form.Label>
                    <Form.Control name={`password`} type="password" id="formPassword" autocomplete="on"/>
                    <Form.Text><Link to="/forgot-password">Forgot your password?</Link></Form.Text>
                </Form.Group>
                <Form.Group>
                    <Button variant="primary">Login</Button>
                </Form.Group>
            </Form>
        </Col>
    </Row>)
}

export default Login