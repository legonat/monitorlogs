import React, { Component } from 'react';
import { Route, Redirect } from 'react-router-dom'
import { AuthContext } from '../App'


const PrivateRoute = ({ component: ComposedComponent, data, ...rest }) => {

  class Authentication extends Component {

    handleRender = props => {
      let isLoggedIn = localStorage.getItem("isAuthenticated")
    if (isLoggedIn === null) {
      isLoggedIn = sessionStorage.getItem("isAuthenticated")
    } 
        if (!isLoggedIn) {
          console.log("Redirecing from Private Route")
        return <Redirect to="/login" />
      } else {
        return <ComposedComponent isAuthenticated={this.props.isAuthenticated} {...data} {...props} />
      }
    }

    render() {
      return (
        <Route {...rest} render={this.handleRender} />
      );
    }
  }

  return (
    <AuthContext.Consumer>
      {
        (value) => <Authentication isAuthenticated={value.isAuthenticated} />
      }
    </AuthContext.Consumer>
  )
};

export default PrivateRoute