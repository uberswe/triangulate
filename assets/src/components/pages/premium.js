import React from 'react'
import {Row, Col, Form, Button} from "react-bootstrap"

function Premium() {
    return (<Row>
        <Col md={6}>
            <h2>Premium: 5 EUR/Month</h2>
            <p>Hey! I'm <a href="https://github.com/uberswe">Markus</a> and I created this website.
                For business purposes I have my own company called <a href="https://www.beubo.com">Beubo</a> which is what will be shown during all transactions and on all invoices.</p>
            <p>With a premium account you help support the costs of running this website and other projects I work on!</p>
            <p>You also receive several perks:</p>
            <ul>
                <li>Premium priority queue</li>
                <li>Image sizes up to 10000x10000</li>
                <li>No Triangulate.xyz watermark</li>
            </ul>
        </Col>
        <Col md={6}>
            <h3>Sign up for Premium</h3>
            <Form>
                <Form.Group controlId="formEmail">
                    <Form.Label>Email</Form.Label>
                    <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on"/>
                    <Form.Text>Your email is stored as a <a href="https://en.wikipedia.org/wiki/SHA-2" target="_blank">SHA-256</a> hash which is only used during login or password reset requests.
                        If you have any problems or questions please contact <a href="mailto:markus@triangulate.xyz">support@triangulate.xyz</a></Form.Text>
                </Form.Group>
                <Form.Group controlId="formPassword">
                    <Form.Label>Password</Form.Label>
                    <Form.Control name={`password`} type="password" id="formPassword" autocomplete="on"/>
                </Form.Group>
                <Form.Group>
                    <Button variant="primary">Proceed to Payment via Stripe</Button>
                </Form.Group>
        </Form>
        </Col>
    </Row>)
}

export default Premium