import './App.scss';
import React from "react";
import Header from "./components/layout/header";
import FourOhFour from "./components/pages/404";
import {Container} from "react-bootstrap";
import Generator from "./components/generator/generator";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import Premium from "./components/pages/premium";
import Footer from "./components/layout/footer";
import Login from "./components/pages/login";
import Logout from "./components/pages/logout";
import ForgotPassword from "./components/pages/forgot-password";
import PrivacyPolicy from "./components/pages/privacy-policy";
import TermsOfService from "./components/pages/terms-of-service";
import axios from "axios";

class App extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            isAuthenticated: false,
            price_id: "",
            stripe_key: "",
        };
    }

    componentDidMount() {
        axios.get ('/api/v1/settings').then (result => {
            this.setState ({
                price_id: result.data.price_id,
                stripe_key: result.data.stripe_key,
                isAuthenticated: result.data.logged_in
            });
        }).catch (error => {
            console.log (error)
        });
    }

    setAuthenticated = (authenticated) => {
        this.setState(state => ({ isAuthenticated: authenticated }));
    };

    render() {
        return (
            <Router>
                <Header isAuthenticated={this.state.isAuthenticated}/>
                <Container style={{'min-height': '40rem'}}>
                    <Switch>
                        <Route exact path="/">
                            <Generator isAuthenticated={this.state.isAuthenticated}/>
                        </Route>
                        <Route path="/premium">
                            <Premium price_id={this.state.price_id} stripe_key={this.state.stripe_key}/>
                        </Route>
                        <Route path="/login">
                            <Login setAuthenticated={this.setAuthenticated} isAuthenticated={this.state.isAuthenticated}/>
                        </Route>
                        <Route path="/forgot-password">
                            <ForgotPassword/>
                        </Route>
                        <Route path="/reset-password">
                            <p>Use this form to reset your password.</p>
                        </Route>
                        <Route path="/terms-of-service">
                            <TermsOfService/>
                        </Route>
                        <Route path="/privacy-policy">
                            <PrivacyPolicy/>
                        </Route>
                        <Route path="/logout">
                            <Logout/>
                        </Route>
                        <Route path="*">
                            <FourOhFour/>
                        </Route>
                    </Switch>
                </Container>
                <Container fluid>
                    <Footer/>
                </Container>
            </Router>
        );
    }
}

export default App;
