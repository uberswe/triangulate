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
import Billing from "./components/pages/billing";
import axios from "axios";
import ResetPassword from "./components/pages/reset-password";

class App extends React.Component {
    constructor(props) {
        super (props);

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

        let owa_baseUrl = 'https://a.beubo.com/';
        const owa_cmds = owa_cmds || [];
        owa_cmds.push (['setSiteId', '808e82efc935cd8ac08ccfc5aac963d2']);
        owa_cmds.push (['trackPageView']);
        owa_cmds.push (['trackClicks']);

        (function () {
            const _owa = document.createElement ('script');
            _owa.type = 'text/javascript';
            _owa.async = true;
            owa_baseUrl = ('https:' == document.location.protocol ? window.owa_baseSecUrl || owa_baseUrl.replace (/http:/, 'https:') : owa_baseUrl);
            _owa.src = owa_baseUrl + 'modules/base/js/owa.tracker-combined-min.js';
            const _owa_s = document.getElementsByTagName ('script')[0];
            _owa_s.parentNode.insertBefore (_owa, _owa_s);
        } ());

    }

    setAuthenticated = (authenticated) => {
        this.setState (state => ({isAuthenticated: authenticated}));
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
                            <Login setAuthenticated={this.setAuthenticated}
                                   isAuthenticated={this.state.isAuthenticated}/>
                        </Route>
                        <Route path="/forgot-password">
                            <ForgotPassword/>
                        </Route>
                        <Route path="/reset-password/:code/" component={ResetPassword}/>
                        <Route path="/terms-of-service">
                            <TermsOfService/>
                        </Route>
                        <Route path="/privacy-policy">
                            <PrivacyPolicy/>
                        </Route>
                        <Route path="/logout">
                            <Logout/>
                        </Route>
                        <Route path="/billing">
                            <Billing/>
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
