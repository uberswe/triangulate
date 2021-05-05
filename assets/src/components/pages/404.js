import React from 'react'
import {Row, Col, Badge} from "react-bootstrap"

function FourOhFour() {
    return (<Row>
        <Col md={12}>
            <h2><Badge variant="secondary">404</Badge> Page Not Found</h2>
            <p>The page you are looking for does not exist.</p>
        </Col>
    </Row>)
}

export default FourOhFour