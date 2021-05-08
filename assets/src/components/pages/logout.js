import React from 'react'
import axios from "axios";

class Logout extends React.Component {
    componentDidMount() {
        axios.post ('/api/v1/logout').then (result => {
            window.location = result.request.responseURL
        }).catch (error => {
            console.log (error)
        });
    }

    render() {
        return ""
    }
}

export default Logout