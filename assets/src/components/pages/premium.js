import React from 'react'
import {Alert, Button, Col, Form, Row} from "react-bootstrap"
import { css } from "@emotion/core";
import PulseLoader from "react-spinners/PulseLoader";
import {Link} from "react-router-dom";

class Premium extends React.Component {
    constructor(props) {
        super (props);
        this.state = {
            price_id: this.props.price_id,
            stripe_key: this.props.stripe_key,
            stripe: null,
            email: "",
            password: "",
            stripeAdded: false,
            isLoading: false,
            error: false,
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
        this.setState({
            isLoading: true
        })
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
                if (result.status !== 200) {
                    this.setState ({
                        isLoading: false,
                        error: true,
                    })
                }
                return result.json ();
            });
        };

        let premium = this;

        createCheckoutSession (this.state.price_id, this.state.email, this.state.password).then (function (data) {
            premium.state.stripe
                .redirectToCheckout ({
                    sessionId: data.sessionId
                })
                .then (window.handleResult).catch(function (error) {
                premium.setState({
                    isLoading: false,
                    error: true
                })
            });
            premium.setState({
                isLoading: false,
                error:false,
            })
        }).catch(function (error) {
            premium.setState({
                isLoading: false,
                error: true
            })
        });
    }

    render() {
        let buttonText = "Pay 5 EUR/mo via Stripe"
        if (this.state.isLoading) {
            buttonText = (<PulseLoader color="#f7f7f7" loading={this.state.isLoading} size={10}/>)
        }
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
                <p>Here are some of the costs of running this website</p>
                <ul>
                    <li>40 EUR/mo - Dedicated server</li>
                    <li>10 EUR/mo - Daily backups</li>
                    <li>A few hours every month goes towards maintaining and updating this website</li>
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
                        <Form.Text>Emails are stored as a <a href="https://en.wikipedia.org/wiki/SHA-2" target="_blank">SHA-256</a> hash.</Form.Text>
                    </Form.Group>
                    <Form.Group controlId="formPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control name={`password`} type="password" id="formPassword" autocomplete="on"
                                      onChange={this.onInputChange} value={this.state.password}/>
                        <Form.Text>Passwords are stored as a <a href="https://en.wikipedia.org/wiki/Bcrypt" target="_blank">Bcrypt</a> hash with a minimum cost of 10.</Form.Text>
                    </Form.Group>
                    <Alert show={this.state.error} variant="danger">
                        An error occurred. Please make sure that you have entered a password at least 8 characters long and a valid email. If you are still having trouble please contact <a
                        href="mailto:support@triangulate.xyz?subject=Triangulate.xyz%3A%20Trouble%20registering%20in">support@triangulate.xyz</a>.
                    </Alert>
                    <Form.Group>
                        <Button disabled={!this.state.stripe || this.state.isLoading} onClick={this.buttonClick} variant="primary">{buttonText}</Button>
                        <Form.Text>By registering you agree to the <Link to="/terms-of-service">Terms of service</Link> and <Link to="/privacy-policy">Privacy Policy</Link></Form.Text>
                    </Form.Group>
                </Form>
            </Col>
        </Row>)
    }
}

export default Premium