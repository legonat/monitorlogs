import React from 'react'

export default function DotsLoader(props){

    // var element = (
    //     <span>
    //         <div className="bounce1"></div>
    //         <div className="bounce2"></div>
    //         <div className="bounce3"></div>
    //     </span>    
    // )

    if (props.color === "white" ) {
        return(
        <div className="spinner-white">
            <div className="bounce1"></div>
            <div className="bounce2"></div>
            <div className="bounce3"></div>
        </div>
        )
    }

    if (props.color === "blue" ) {
        return(
        <div className="spinner-blue">
            <div className="bounce1"></div>
            <div className="bounce2"></div>
            <div className="bounce3"></div>
        </div>
        )
    }



}