import React from 'react'
import {Nav, Navbar} from 'react-bootstrap'
import {Link} from "react-router-dom";

class Header extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        let nav = ""
        if (this.props.isAuthenticated) {
            nav = (<Nav className="mr-auto">
                <Link to="/" class="nav-link">Home</Link>
                <Link to="/billing" class="nav-link">Manage Billing</Link>
                <Link to="/logout" class="nav-link">Logout</Link>
            </Nav>)
        } else {
            nav = (<Nav className="mr-auto">
                <Link to="/" class="nav-link">Home</Link>
                <Link to="/premium" class="nav-link">Premium</Link>
                <Link to="/login" class="nav-link">Login</Link>
            </Nav>)
        }

        return (<Navbar bg="dark" variant="dark" expand="lg" style={{'margin-bottom': '30px'}}>
            <Navbar.Brand href="/">Triangulate.xyz</Navbar.Brand>
            <Navbar.Toggle aria-controls="basic-navbar-nav"/>
            <Navbar.Collapse id="basic-navbar-nav">
                {nav}
            </Navbar.Collapse>
        </Navbar>)
    }
}

export default Header