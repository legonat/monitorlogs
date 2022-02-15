import React from 'react'
import { Redirect } from 'react-router-dom'
import LoginComponent from "./login.component"
// import Cookies from 'universal-cookie'

// const cookies = new Cookies()

class LoginContainer extends React.Component {


  render() {
    // var cookie = cookies.get('rtoken')
    // if (cookie !== undefined) {
    //   return <Redirect to="/" />
    // }
    return <LoginComponent user={this.props.user} pass={this.props.pass} handleUsernameInput={this.props.handleUsernameInput}
      handlePassInput={this.props.handlePassInput} handleLoginSubmit={this.props.handleLoginSubmit}
      handleRememberMeCheckbox={this.props.handleRememberMeCheckbox} canRememberMe={this.props.canRememberMe}
      submitted={this.props.submitted}/>;
  }
}

export default (LoginContainer)