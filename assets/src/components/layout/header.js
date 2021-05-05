import React from 'react'
import {Nav, Navbar, Row} from 'react-bootstrap'
import {Link} from "react-router-dom";

function Header() {
    return (<Navbar bg="dark" variant="dark" expand="lg" style={{'margin-bottom': '30px'}}>
        <Navbar.Brand href="/">Triangulate.xyz</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav"/>
        <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="mr-auto">
                <Link to="/" class="nav-link">Home</Link>
                <Link to="/premium" class="nav-link">Premium</Link>
                <Link to="/login" class="nav-link">Login</Link>
            </Nav>
        </Navbar.Collapse>
    </Navbar>)
}

export default Header