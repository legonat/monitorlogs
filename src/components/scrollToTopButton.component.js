import React, { Component } from "react";
import { Button } from "react-bootstrap";
import { faArrowCircleUp } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { AuthContext } from '../App'

export default class NavbarContainer extends Component {
  render() {

    const scrollToTop = () => {
        window.scrollTo({top: 0, left: 0, behavior: 'smooth'})
    }

    return (
      <AuthContext.Consumer>
        { value =>
          !value.isAuthenticated
            ? <div></div>
            : <Button variant="outline-light" onClick={scrollToTop} className="nav-to-top-btn">
            <FontAwesomeIcon icon={faArrowCircleUp} size="2x" />
          </Button>
        }
        </AuthContext.Consumer>
    )
  }
}