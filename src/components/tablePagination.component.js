import React, { Component } from "react";
import {Pagination} from 'react-bootstrap'
import PageItem from 'react-bootstrap/PageItem'
import { loggers } from "winston";



export default class TablePagination extends Component {


     handlePageClick = event => {
        if (event.target.nodeName !== 'SPAN') { 
           var pageNumber= parseInt(event.target.text) 
           this.props.handlePageChange(pageNumber)
        }
      }
  
       handleFirstPageClick = () => {
        this.props.handlePageChange(1)
      }
  
       handleLastPageClick = () => {
        this.props.handlePageChange(this.props.pageCount )
      }
  
       handleForwardClick = () => {
        if (this.props.pageCount > this.props.currentPage) {
          var pageNumber= this.props.currentPage + 1 
          this.props.handlePageChange(pageNumber)
          
        }
      }
  
       handleBackwardsClick = () => {
        if (this.props.currentPage > 1) {
          var pageNumber= this.props.currentPage - 1
          this.props.handlePageChange(pageNumber)
        }
      }
  
       handleFastForwardClick = () => {
        if (this.props.pageCount - this.props.currentPage > 10) {
          var pageNumber= this.props.currentPage + 10
          this.props.handlePageChange(pageNumber)
        }
      }
  
       handleFastBackwardsClick = () => {
        if (this.props.currentPage > 11) {
          var pageNumber= this.props.currentPage - 11 
          this.props.handlePageChange(pageNumber)
        }
      }



      render() {
        const pageCount = this.props.pageCount
        const currentPage = this.props.currentPage
        
        var elements = 0
        var offset = 0
        switch (true) {
            case (pageCount < 25): 
                elements = 25
                offset = 12
                break
            case (pageCount > 25 && currentPage < 90): 
                elements = 18
                offset = 9
                break
            case (pageCount > 25 && currentPage < 990): 
                elements = 14
                offset = 7
                break
            case (currentPage > 1000): 
                elements = 12
                offset = 6
                break
            default: {
              elements = 12
              offset = 6
            }
        }   
        
        var number = 1
        let intMultiplier = Math.trunc(currentPage/elements)
        var sliceNum = currentPage - (elements * intMultiplier)
        let items = [];

        if (currentPage >= elements * intMultiplier && intMultiplier > 0) {
          number = intMultiplier*elements
        }
        for (var i = 1; i <= elements; i++) {
            if (currentPage > elements - 1 && i === 1 && pageCount > elements) {
                items.push(<Pagination.Ellipsis onClick = {this.handleFastBackwardsClick}/>)
                //number = currentPage - offset + 1
            }
            items.push(
            <Pagination.Item key={number} active={number === currentPage} onClick={this.handlePageClick}>
            {number}
            </Pagination.Item>,
            );
            number++


            if (pageCount - currentPage > offset && i === elements) {
                items.push(<Pagination.Ellipsis onClick = {this.handleFastForwardClick}/>)
            }

            if (number > pageCount) {break}             
        }

        // var number = 1
        // let items = [];
        // for (var i = 1; i <= elements; i++) {

        //     if (currentPage > offset && i === 1 && pageCount > 25) {
        //         items.push(<Pagination.Ellipsis onClick = {this.handleFastBackwardsClick}/>)
        //         number = currentPage - offset + 1
        //     }
        //     items.push(
        //     <Pagination.Item key={number} active={number === currentPage} onClick={this.handlePageClick}>
        //     {number}
        //     </Pagination.Item>,
        //     );
        //     number++


        //     if (pageCount - currentPage > offset && i === elements) {
        //         items.push(<Pagination.Ellipsis onClick = {this.handleFastForwardClick}/>)
        //     }

        //     if (number > pageCount) {break}             
        // }


        return (
        <div>
            <Pagination 
            // current={this.props.current}
             // currentPage={currentPage}
              //currentPage={currentPage}
              total={this.props.pageCount}
            // showLessItems
            // showTitle={false}
            >
                <Pagination.First onClick = {this.handleFirstPageClick}/>
                <Pagination.Prev onClick = {this.handleBackwardsClick}  />
                {items}
                <Pagination.Next onClick = {this.handleForwardClick} />
                <Pagination.Last onClick = {this.handleLastPageClick} />
            </Pagination>  
        </div>
        )
      }

}




//   Backup if needed to rise CurrentPage state into App
            // const handlePageClick = event => {
            //     if (event.target.nodeName !== 'SPAN') { 
            //        setCurrentPage(parseInt(event.target.text))
            //     }
            //   }
          
            //   const handleFirstPageClick = () => {
            //     setCurrentPage(1)
            //   }
          
            //   const handleLastPageClick = () => {
            //     setCurrentPage(logsArrayLenght)
            //   }
          
            //   const handleForwardClick = () => {
            //     if (logsArrayLenght > currentPage) {
            //       setCurrentPage(currentPage + 1)
            //     }
            //   }
          
            //   const handleBackwardsClick = () => {
            //     if (currentPage > 1) {
            //       setCurrentPage(currentPage - 1)
            //     }
            //   }
          
            //   const handleFastForwardClick = () => {
            //     if (logsArrayLenght - currentPage > 5) {
            //       setCurrentPage(currentPage + 5)
            //     }
            //   }
          
            //   const handleFastBackwardsClick = () => {
            //     if (currentPage > 6) {
            //       setCurrentPage(currentPage - 5)
            //     }
            //   }
