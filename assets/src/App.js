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
import ForgotPassword from "./components/pages/forgot-password";
import PrivacyPolicy from "./components/pages/privacy-policy";
import TermsOfService from "./components/pages/terms-of-service";

class App extends React.Component {

    render() {
        return (
            <Router>
                <Header/>
                <Container style={{'min-height': '40rem'}}>
                    <Switch>
                        <Route exact path="/">
                            <Generator/>
                        </Route>
                        <Route path="/premium">
                            <Premium/>
                        </Route>
                        <Route path="/login">
                            <Login/>
                        </Route>
                        <Route path="/members">
                            <p>Members area</p>
                            <Generator/>
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
