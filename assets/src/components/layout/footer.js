import React from 'react'
import {Col, Row} from 'react-bootstrap'
import {Link} from "react-router-dom";

function Footer() {
    return (<Row style={{'margin-bottom': '50px'}}>
        <Col md={12}>
            <hr/>
            Created by <a href="https://github.com/uberswe">Markus Tenghamn</a> | Contact <a
            href="mailto:markus@triangulate.xyz">support@triangulate.xyz</a>
            <br/>
            This project uses ideas and code from <a
            href="https://github.com/esimov/triangle">triangle</a> and <a
            href="https://github.com/preslavrachev/generative-art-in-go">generative-art-in-go</a>
            <br/>
            <Link to="/terms-of-service">Terms of service</Link> | <Link to="/privacy-policy">Privacy Policy</Link>
        </Col>
    </Row>)
}

export default Footer