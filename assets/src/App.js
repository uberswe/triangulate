import './App.scss';
import React from "react";
import Header from "./components/layout/header";
import FourOhFour from "./components/pages/404";
import {Container, Row} from "react-bootstrap";
import Generator from "./components/generator/generator";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import Premium from "./components/pages/premium";
import Footer from "./components/layout/footer";
import Login from "./components/pages/login";
import ForgotPassword from "./components/pages/forgot-password";

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
                        <Route path="/forgot-password">
                            <ForgotPassword/>
                        </Route>
                        <Route path="/reset-password">
                            <p>Use this form to reset your password.</p>
                        </Route>
                        <Route path="/terms-of-service">
                            <p>N/A</p>
                        </Route>
                        <Route path="/privacy-policy">
                            <p>N/A</p>
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
