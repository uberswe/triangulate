import React from 'react'
import {Button, Col, Form, Row} from "react-bootstrap"

class Premium extends React.Component {
    constructor(props) {
        super (props);
        this.state = {
            price_id: this.props.price_id,
            stripe_key: this.props.stripe_key,
            stripe: null,
            email: "",
            password: "",
            stripeAdded: false
        };

        this.buttonClick = this.buttonClick.bind (this);
        this.onInputChange = this.onInputChange.bind (this);
    }

    onInputChange(event) {
        this.setState ({
            [event.target.name]: event.target.value
        });
    }

    loadStripe(stripe_key) {
        if (!this.state.stripeAdded && stripe_key !== "") {
            this.setState ({
                stripeAdded: true
            })
            const script = document.createElement ("script");
            script.src = "https://js.stripe.com/v3/";
            script.async = true;
            script.onload = () => this.stripeLoaded(stripe_key);
            document.body.appendChild (script);
        }
    }

    stripeLoaded(stripe_key) {
        this.setState ({
            stripe: window.Stripe(stripe_key)
        })
    }

    componentDidMount() {
        this.loadStripe(this.state.stripe_key)
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (this.props !== prevProps) {
            this.setState ({
                price_id: this.props.price_id,
                stripe_key: this.props.stripe_key,
            })
        }
        this.loadStripe(this.state.stripe_key)
    }

    buttonClick() {
        let createCheckoutSession = function (priceId, email, password) {
            return fetch ("/api/v1/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify ({
                    priceId: priceId,
                    email: email,
                    password: password
                })
            }).then (function (result) {
                return result.json ();
            });
        };

        let premium = this;

        createCheckoutSession (this.state.price_id, this.state.email, this.state.password).then (function (data) {
            premium.state.stripe
                .redirectToCheckout ({
                    sessionId: data.sessionId
                })
                .then (handleResult);
        });
    }

    render() {
        return (<Row>
            <Col md={6}>
                <h2>Premium: 5 EUR/mo</h2>
                <p>Hey! I'm <a href="https://github.com/uberswe">Markus</a> and I created this website.
                    For business purposes I have my own company called <a href="https://www.beubo.com">Beubo</a> which
                    is what will be shown during all transactions and on all invoices.</p>
                <p>With a premium account you help support the costs of running this website and other projects I work
                    on!</p>
                <p>You also receive several perks:</p>
                <ul>
                    <li>Premium priority queue</li>
                    <li>Image sizes up to 10000x10000</li>
                    <li>No Triangulate.xyz watermark</li>
                </ul>
                <p>Interested in something more? Contact me at <a
                    href="mailto:markus@triangulate.xyz?subject=Triangulate.xyz%3A%20More%20than%20premium">support@triangulate.xyz</a>.
                </p>
            </Col>
            <Col md={6}>
                <h3>Sign up for Premium</h3>
                <Form>
                    <Form.Group controlId="formEmail">
                        <Form.Label>Email</Form.Label>
                        <Form.Control name={`email`} type="email" id="formEmail" autocomplete="on"
                                      onChange={this.onInputChange} value={this.state.email}/>
                        <Form.Text>Your email is stored as a <a href="https://en.wikipedia.org/wiki/SHA-2"
                                                                target="_blank">SHA-256</a> hash which is only used
                            during login or password reset requests.
                            If you have any problems or questions please contact <a
                                href="mailto:markus@triangulate.xyz?subject=Triangulate.xyz%3A%20Sign%20up%20support">support@triangulate.xyz</a></Form.Text>
                    </Form.Group>
                    <Form.Group controlId="formPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control name={`password`} type="password" id="formPassword" autocomplete="on"
                                      onChange={this.onInputChange} value={this.state.password}/>
                    </Form.Group>
                    <Form.Group>
                        <Button disabled={!this.state.stripe} onClick={this.buttonClick} variant="primary">Pay 5 EUR/mo via Stripe</Button>
                    </Form.Group>
                </Form>
            </Col>
        </Row>)
    }
}

export default Premium