import React from 'react'

class Billing extends React.Component {

    componentDidMount() {
        fetch('/api/v1/portal', {
            method: 'POST'
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(data) {
                window.location.href = data.url;
            })
            .catch(function(error) {
                console.error('Error:', error);
            });
    }

    render() {
        return ""
    }
}

export default Billing