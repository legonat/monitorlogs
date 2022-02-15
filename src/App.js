import React, { useState, useEffect } from 'react';
import '../node_modules/bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import { HashRouter, Switch, Route, useHistory, useLocation, matchPath, useRouteMatch } from "react-router-dom";
import Cookies from "universal-cookie"

import LoginContainer from "./components/LoginContainer";
import Logout from "./components/logout.component"
import Register from "./components/register.component";
import Board from "./components/board.component"
import Navbar from "./components/navbar.component"
import PrivateRoute from "./components/PrivateRoute"
import ScrollButton from "./components/scrollToTopButton.component"
import Logger from "./utils/logger"
import AuthModule from './utils/authModule';

import DatePicker from "react-datepicker";
//import addDays from 'date-fns/addDays';
import "react-datepicker/dist/react-datepicker.css";
//import { render } from '@testing-library/react';
import { Integrations } from '@sentry/tracing';
import { loggers } from 'winston';
import { useGestures } from 'framer-motion';




const serverUrl = "http://localhost:5000"
const cookies = new Cookies()
const auth = new AuthModule()
const logger = new Logger()
export const AuthContext = React.createContext();



export default function App() {



  const history = useHistory()
  const location = useLocation()
  const match = useRouteMatch("/")
  const ItemsPerPage = 200
  const PreloadPageCount = 10

  const [user, setUser] = useState("")
  const [pass, setPass] = useState("")
  // const [logs, setLogs] = useState([])
  // const [errors, setErrors] = useState([])
  // const [searchResults, setSearchResults] = useState([])
  const [confirmPass, setConfirmPass] = useState("")
  const [passwordMatch, setPasswordMatch] = useState(true)
  const [searchText, setSearchText] = useState("")
  const [completeSearchRequest, setCompleteSearchRequest] = useState("")
  const [submitted, setSubmitted] = useState(false)
  const [canRememberMe, setCanRememberMe] = useState(false)
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isAuthorized, setIsAuthorized] = useState(false)
  //const [authToken, setAuthToken] = useState("") //TODO keep authToken hashed
  const [selectedSessionId, setSelectedSessionId] = useState(0)
  const [selectedSession, setSelectedSession] = useState("")
  const [sessions, setSessions] = useState([])
  const [services, setServices] = useState([])
  const [filenames, setFilenames] = useState([])
  const [selectedFilename, setSelectedFilename] = useState("")
  const [selectedService, setSelectedService] = useState("")
  const [filterText, setFilterText] = useState("")
  const [filteredLogs, setFilteredLogs] = useState([])
  const [filteredErrors, setFilteredErrors] = useState([])
  const [filteredSearch, setFilteredSearch] = useState([])
  const [loading, setLoading] = useState(false)
  const [currentTab, setCurrentTab] = useState(0)
  const [logsPageCount, setLogsPageCount] = useState(0)
  const [errorsPageCount, setErrorsPageCount] = useState(0)
  const [searchPageCount, setSearchPageCount] = useState(0)
  const [selectedLogId, setSelectedLogId] = useState(0)
  const [searchedPage, setSearchedPage] = useState(0)
  const [searchedSessionId, setSearchedSessionId] = useState(0) //TODO check if needed (store in session)
  const [searchedFilename, setSearchedFilename] = useState("")
  const [currentLogsPage, setCurrentLogsPage] = useState(1)
  const [currentErrorsPage, setCurrentErrorsPage] = useState(1)
  const [currentSearchPage, setCurrentSearchPage] = useState(1)
  const [isFingerprintReady, setIsFingerprintReady] = useState(false)
  const [isAuthNeeded, setIsAuthNeeded] = useState(false)
  const [dateRange, setDateRange] = useState([null, null]);
  const [startDate, endDate] = dateRange;
  const [isOpen, setIsOpen] = useState(false);
  const [needScrollDown, setNeedScrollDown] = useState(true)
  const [resourses, setResources] = useState([])
  const [options, setOptions] = useState([])
  const [offset, setOffset] = useState(0)
  const [fetchedLogsLength, setFetchedLogsLength] = useState(0)
  const [fetchedErrorsLength, setFetchedErrorsLength] = useState(0)
  const [fetchedSearchLength, setFetchedSearchLength] = useState(0)
  const [currentSlicePage, setCurrentSlicePage] = useState(0)
  const [isGoingForward, setIsGoingForward] = useState(false)
  const [isGoingBackwards, setIsGoingBackwards] = useState(false)
  const [forwardPreloadLogs, setForwardPreloadLogs] = useState([])
  const [backwardPreloadLogs, setBackwardPreloadLogs] = useState([])
  const [currentLogSlice, setCurrentLogSlice] = useState([])
  const [isLoadingSilent, setIsLoadingSilent] = useState(false)
  const [isFetchSuccessful, setIsFetchSuccessful] = useState(true)
  const [fetchedLogsSlice, setFetchedLogsSlice] = useState([])
  const [fetchedErrorsSlice, setFetchedErrorsSlice] = useState([])
  const [fetchedSearchSlice, setFetchedSearchSlice] = useState([])
  const [fetchedFirstLogRowId, setFetchedFirstLogRowId] = useState(0)
  const [isLogsUpdateNeeded, setIsLogsUpdateNeeded] = useState(false)
  const [isErrorsUpdateNeeded, setIsErrorsUpdateNeeded] = useState(false)
  const [logsOffset, setLogsOffset] = useState(0)
  const [errorsOffset, setErrorsOffset] = useState(0)



  // const onChange = (dates) => {
  //   const [start, end] = dates;
  //   setStartDate(start);
  //   setEndDate(end);
  // }

  //<DatePicker selected={startDate} onChange={(date) => setStartDate(date)} />

  const ExampleCustomInput = React.forwardRef(({ value, onClick }, ref) => (
    <button className="example-custom-input" onClick={onClick} ref={ref}>
      {value}
    </button>
  ));

  const setValue = (value) => {
    if (value !== null) {
      return (React.forwardRef(({ value, onClick }, ref) => (
        <button className="example-custom-input" onClick={onClick} ref={ref}>
          {value}
        </button>
      ))
      )
    } else {
      return (React.forwardRef(({ value, onClick }, ref) => (
        <button className="example-custom-input" onClick={onClick} ref={ref}>
          {value}
        </button>
      ))
      )
    }
  }

  var datePicker = <div> <DatePicker
    className="nav-form col-auto form-control"
    onInputClick={() => { setIsOpen(true); setNeedScrollDown(false) }}
    onClickOutside={() => { setIsOpen(false); setNeedScrollDown(true) }}
    open={isOpen}
    selectsRange={true}
    startDate={startDate}
    endDate={endDate}
    onChange={(update) => {
      setDateRange(update);
      if (update[1]) {
        logger.console("Date Picker Closed")
        setIsOpen(false);
        setNeedScrollDown(true);
        setSelectedSession('')
        setSelectedSessionId(0)
        setIsLogsUpdateNeeded(true);
        setIsErrorsUpdateNeeded(true)
      }
    }}
    isClearable={true}
    placeholderText="Select a date"
    maxDate={new Date()}
    dateFormat="yyyy-MM-dd"
    showYearDropdown
    dateFormatCalendar="MMMM"
    yearDropdownItemNumber={15}
    scrollableYearDropdown
    //monthsShown={2}
    todayButton="Today"
    strictParsing
  />
  </div>
  // <DatePicker
  //     selected={startDate}
  //     onChange={onChange}
  //     startDate={startDate}
  //     endDate={endDate}
  //     selectsRange
  //     inline
  //   />



  // TODO Add aToken validity Check

  // logger.log("info", "TESSTING")

  // const pathname = UsePathname();

  // const checkLocation = () => {


  //   logger.console(pathname)

  //   if (pathname === "Login" || pathname === "Registration" || pathname === "Logout") {
  //     setIsAuthNeeded(false)
  //     return
  //   }
  //   if (patname === "Root") {
  //     setIsAuthNeeded(true)
  //   }
  // }

  // checkLocation()

  const handleUsernameInput = event => {
    setUser(event.target.value);
  }

  const handleFilterTextInput = event => {
    setFilterText(event.target.value);
    if (event.target.value === "") {
      setLogsPageCount(getTotalLength(fetchedLogsLength))
      setErrorsPageCount(getTotalLength(fetchedErrorsLength))
      setSearchPageCount(getTotalLength(fetchedSearchLength))
    }
  }

  const handlePassInput = event => {
    let passV = event.target.value;
    setPass(passV);
    if (passV !== confirmPass) {
      setPasswordMatch(false)
    }
    else {
      setPasswordMatch(true)
    }
  }
  const handleConfirmPassInput = event => {
    let confirmPassV = event.target.value;
    setConfirmPass(confirmPassV);
    if (confirmPassV !== pass) {
      setPasswordMatch(false)
    }
    else {
      setPasswordMatch(true)
    }
  }

  const handleRememberMeCheckbox = event => {
    let isChecked = event.target.checked;
    setCanRememberMe(isChecked)
  }



  const handleSearchInput = event => {

    if (event.target.nodeName !== 'INPUT') {
      setSelectedLogId(0)
      // setSearchResults([])
      setSearchedSessionId(0)
      setSearchedFilename("")
      setSearchText("")
      setCompleteSearchRequest('')
      setFetchedSearchSlice([])
      return
    }

    setSearchText(event.target.value)
  }

  function handleErrorSilent(response) {
    if (!response.ok) {
      logger.err(response.status + ": " + response.statusText);
      return null
    }
    return response.json();
  }

  function handleErrorAlert(response) {

    if (response.ok === undefined) {
      throw Error("Request error")
    }
    if (response.status === 401) {
      alert("401: Unauthorized")
      goToHomepage()
      return
    }
    if (!response.ok) {
      logger.err(response)
      throw Error(response.statusText);
    }
    if (response.status !== 200) {
      logger.err(response)
      throw Error(response.statusText);
    }
    logger.console(response)
    return response.json()
  }


  const handleSelectSessionChange = event => {
    logger.console('selected session ' + event)
    var id = parseInt(event)
    // const selectedIndex = event.target.options.selectedIndex;
    // var id = parseInt(event.target.options[selectedIndex].getAttribute('data-key'))
    // setselectedSessionId(id);
    // setSelectedSession(event.target.value)

    //setLogs([])
    //resetPagination()
    setSelectedSessionId(id);
    //setSelectedSession(sessions[id - 1].value)
    //handleGetLogsBySession(selectedSessionId)
  }

  const handleSelectFilenameChange = (event) => {
    logger.console("selected logfileId" + event)
    let id = parseInt(event)
    let filename = filenames[id - 1].value
    console.info("Selected Filename" + filename)
    setSessions([])
    setSelectedLogId(0)
    setSelectedFilename(filename)
    //resetPagination()
    //setSelectedSession("")
    setSelectedSessionId(0)
    setSelectedSession('')
    setIsLoadingSilent(false)
    // setIsErrorsUpdateNeeded(true)
    // setIsLogsUpdateNeeded(true)
  }

  const handleSelectServiceChange = event => {
    if (event.target.nodeName !== 'SELECT') {
      setSelectedService("")
      setLogsPageCount(getTotalLength(fetchedLogsLength))
      setErrorsPageCount(getTotalLength(fetchedErrorsLength))
      setSearchPageCount(getTotalLength(fetchedSearchLength))
      logger.console("Selected node " + event.target.nodeName)
      return
    }
    logger.console("Selected " + event.target)
    setSelectedService(event.target.value)
  }



  function handleTabClick(i) {
    setCurrentTab(i)
  }

  function resetPagination() {
    setCurrentLogsPage(1)
    setCurrentErrorsPage(1)
    setCurrentSearchPage(1)
    setSearchedPage(0)
    setOffset(0)
    logger.console("REsetPagination Complete")
  }

  function goToHomepage() {
    history.replace('/')
    history.go(0)
  }

  function reloadPage() {
    window.location.reload()
  }

  function goToLoginPage() {
    localStorage.clear()
    sessionStorage.clear()
    history.push('#/login')
    //history.go(0)
  }

  function clearAndLogout() {
    localStorage.clear()
    sessionStorage.clear()
    reloadPage()
  }

  function checkAuth() {
    return (
      new Promise((function (resolve, reject) {
        if (auth.checkToken) {
          setIsAuthorized(true)
          resolve(true)
          return
        }
        if (!auth.checkToken) {
          setIsAuthorized(false)  // if false - setLoading  
          resolve(false)
          return
        }
      }))
    )
  }

  function checkNeedLoading() {
    if (isAuthenticated && isFetchSuccessful) {
      return false
    } else {
      return true
    }
  }


  function getTotalLength(len) {
    let i = 0
    logger.console('Total Logs count ' + len)
    i = ((len - (len % ItemsPerPage)) / ItemsPerPage) + 1
    logger.console('Total Pages Count ' + i)

    return i
  }

  // function getSearchedPage(len) {
  //   logger.console("Calculating slice: " + len)
  //   if (len <= 0) {
  //     setLoading(false)
  //     setSubmitted(false)
  //     return
  //   }
  //   let i = ((len - (len % ItemsPerPage)) / ItemsPerPage) + 1
  //   logger.console("Calculated page: " + i)
  //   setSearchedPage(i)
  //   // setCurrentLogsPage(i)
  //   // setCurrentSearchPage(i)
  //   setLoading(false)
  //   setSubmitted(false)
  // }





  function SendRegisterRequest(input) {
    const path = serverUrl + '/v2/registrationAttempt'
    auth.secureFetch(path, {
      method: 'POST',
      body: JSON.stringify(input)
    })
      .then(res => {
        if (res.status === 200) {
          logger.info('Registeration Successful')
          goToHomepage()
        } else {
          throw Error(res.statusText)
        }
      })
      .catch(error => {
        logger.console('Catching error')
        logger.err(error)
        goToHomepage()
        return
      })
  }

  function handleLoginSubmit() {
    setSubmitted(true)
    const path = serverUrl + '/v2/loginAttempt'
    const input = { login: user, password: pass, rememberMe: canRememberMe };
    logger.console(process.env.NODE_ENV)
    auth.secureAuth(path, {
      method: 'POST',
      body: JSON.stringify(input)
    })
      .then(res => {
        return handleErrorAlert(res) //TODO Check why continues to LoginSuccess
      })
      .then(data => {
        logger.info('LoginSuccess')
        logger.console(data)
        if (canRememberMe) {
          localStorage.setItem("isAuthenticated", true)
        }
        if (!canRememberMe) {
          sessionStorage.setItem("isAuthenticated", true)
        }
        goToHomepage()
      })
      .catch(error => {
        logger.console('Catching error')
        logger.err(error)
        goToHomepage()
      })
  }

  function handleLogoutSubmit() {
    //setSubmitted(true)
    const path = serverUrl + '/v2/logoutAttempt'
    auth.secureFetch(path)
      .then(res => {
        if (res.status === 200) {
          logger.info('LogoutSuccess')
          clearAndLogout()
          // localStorage.clear()
          // sessionStorage.clear()

          // reloadPage()
          //goToHomepage()
        } else {
          throw Error(res.statusText)
        }
      })
      .catch(error => {
        logger.console('Catching error')
        logger.err(error)
        goToHomepage()
      })
  }

  // function handleAuthSubmit() {
  //   setSubmitted(true)
  //   logger.console("auth attempt func")
  //   const path = serverUrl + '/v2/authAttempt'
  //   auth.secureAuth(path)
  //     .then(res => {
  //       logger.console(res)
  //       if (res.status == 200) {
  //         return res.json()
  //       }
  //       if (res.status === 401) {
  //         alert(res.statusText)
  //         goToLoginPage()
  //       }
  //     })
  //     .then(data => {
  //       logger.console('AuthSuccess')
  //       logger.console(data)
  //       auth.setToken(data.accessToken)
  //       setSubmitted(false)
  //       setIsAuthorized(true)
  //     })
  //     .catch(error => {
  //       logger.console('Catching error')
  //       logger.console(error)
  //       logger.err(error)
  //       goToHomepage()
  //     })
  // }



  // function handleUnblockSubmit() {
  //   //setSubmitted(true)
  //   const path = serverUrl + '/v2/private/unblock/'
  //   const unblockInput = { login: "test5" };
  //   auth.secureFetch(path, {
  //     method: 'POST', 
  //     body: JSON.stringify(unblockInput)
  //   })
  //   .then(res => {
  //     if (res.status == 200) {
  //       return res.json()
  //     } else {
  //       throw Error(res.statusText)
  //     }
  //   })
  //   .then(data => {
  //       logger.console('UnblockSuccess')
  //       alert("User is successfuly unblocked")
  //     })
  //     .catch(error => {
  //       logger.console('Catching error')
  //       logger.err(error)
  //       goToHomepage()
  //     })
  // }

  // function handleUNSECUREDGetLogsSubmit() {
  //   setSubmitted(true)
  //   const path = serverUrl + '/v2/getLogs'
  //   //   const config = {
  //   //     headers: { Authorization: `Bearer ${authToken}` }
  //   // };
  //   //const input = {login: user, password: pass, fingerprint:fingerprint};
  //   fetch(path)
  //     .then(res => {
  //       if (res.status == 200) {
  //         return res.json()
  //       } else {
  //         throw Error(res.statusText)
  //       }
  //     })
  //     .then(data => {
  //       //cookies.set("rToken", res.data.refreshToken, {path: '/'})
  //       setLogs(data.logs)
  //       //setSubmitted(false)
  //     })
  //         .catch(error => {
  //           logger.console('Catching error')
  //           logger.err(error)
  //           goToHomepage()
  //         })
  //     }

  // function handleGetLogsSubmit() {
  //   setSubmitted(true)
  //   const path = serverUrl + '/v2/private/getLogs'
  //   auth.secureFetch(path)
  //     .then(res => {
  //       return handleErrorAlert(res)
  //     })
  //     .then(data => {
  //       logger.console('GetLogsSuccess')
  //       setLogs(data.logs)
  //       setSubmitted(false)
  //     })
  //     .catch(error => {
  //       logger.console('Catching error')
  //       logger.err(error)
  //     })
  // }



  function handleGetLogsSessions() {
    return (
      new Promise((function (resolve, reject) {

        logger.console(selectedFilename)
        if (selectedFilename === "") {
          return
        }

        //const logFilename = localStorage.getItem("lastSelectedFilename")
        const path = serverUrl + '/v2/private/getLogsSessions'
        var logsInputs = { value: selectedFilename }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(logsInputs)
        })
          .then(res => {
            return handleErrorSilent(res)
          })
          .then(data => {
            if (data === null) {
              return
            }
            logger.console('GetLogSessionsSuccess')
            const sessions = data.sessions
            //var sessionId = sessions[sessions.length - 1].id
            // if (selectedSessionId === 0) {
            //   setSelectedSessionId(sessionId)
            // }
            setSessions(sessions)
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.console(error)
            logger.err(error)
            reject(error)
            goToHomepage()
          })

      })
      )
    )
  }

  function handleGetLogsFilenames() {
    logger.console("Getting log filenames")
    const path = serverUrl + '/v2/private/getLogsFilenames'
    auth.secureFetch(path)
      .then(res => {
        return handleErrorAlert(res)
      })
      .then(data => {
        logger.console(data)
        logger.console(data.results)
        setFilenames(data.results)
      })
      .catch(error => {
        logger.console('Catching error')
        logger.console(error)
        logger.err(error)
        goToHomepage()
      })
  }


  // function handleGetLogsFilenames() {
  //   setSubmitted(true)
  //   const path = serverUrl + '/v2/private/getLogsFilenames'
  //   const config = {
  //     headers: { Authorization: `Bearer ${authToken}` }
  //   };
  //   axios.get(path, config)
  //     .then(res => {
  //       logger.console('GetLogsFilenamesSuccess')
  //       logger.console(res.data)
  //       setFilenames(res.data.results)
  //       setSubmitted(false)
  //     })
  //     .catch(error => {
  //       logger.console('Catching error')
  //       if (error.response !== undefined && error.response.status === 401) {
  //         logger.err(error)
  //         setAuthToken("")
  //         history.replace('/')
  //         history.go(0)
  //         return
  //       }
  //       logger.err(error)
  //       goToHomepage()
  //     })
  // }

  function handleGetLogServiceInfo() {
    if (selectedFilename === "") {
      return
    }
    const path = serverUrl + '/v2/private/getLogsServiceInfo'
    var logsInputs = { value: selectedFilename }
    auth.secureFetch(path, {
      method: 'POST',
      body: JSON.stringify(logsInputs)
    })
      .then(res => {
        return handleErrorAlert(res)
      })
      .then(data => {
        if (data === null) {
          return
        }
        logger.console('GetLogServiceInfoSuccess')
        setServices(data.services)
      })
      .catch(error => {
        logger.console('Catching error')
        logger.err(error)
        goToHomepage()
      })
  }

  function handleGetLogsByDate() {

    return (
      new Promise((resolve, reject) => {
        logger.console(selectedFilename)
        if (selectedFilename === "") {
          logger.console("No filename selected")
          return
        }
        setLoading(true)
        // if (logs[0] !== undefined) {
        //   currentSessionId = logs[0].sessionId
        // }

        const path = serverUrl + '/v2/private/getLogsByDate'
        var sessionInputs = { startDate: startDate, endDate: endDate, filename: selectedFilename }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(sessionInputs)
        })
          .then(res => {
            return handleErrorAlert(res)
          })
          .then(data => {
            logger.console('GetLogsByDateSuccess')
            var errors = data.errors
            var logs = data.logs
            setLogsPageCount(getTotalLength(logs.length))
            setErrorsPageCount(getTotalLength(errors.length))
            //setLogs(logs)
            //setErrors(errors)
            //setLoading(false)
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.err(error)
            history.replace('/') // make it for unautorized/ //TODO
            history.go(0)
            reject(error)
          })
      })
    )
  }

  function handleGetLogsSliceByDate() {

    return (
      new Promise((resolve, reject) => {
        logger.console(selectedFilename)
        if (selectedFilename === "") {
          logger.console("No filename selected")
          return
        }
        if (!isFetchSuccessful) {
          setLoading(checkNeedLoading())
          logger.console("Preventing Request Dupe")
          return
        }
        setIsFetchSuccessful(false)
        setLoading(!isLoadingSilent)

        const path = serverUrl + '/v2/private/getLogsByDateWithLimit'
        var sessionInputs = { startDate: startDate, endDate: endDate, filename: selectedFilename, limit: PreloadPageCount * ItemsPerPage, offset: offset * ItemsPerPage }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(sessionInputs)
        })
          .then(res => {
            return handleErrorAlert(res)
          })
          .then(data => {
            logger.console('GetLogsByDateSuccess')
            setLogsPageCount(getTotalLength(data.logsCount))
            //setErrorsPageCount(getTotalLength(data.errorsCount))
            setFetchedLogsLength(data.logsCount)
            //setFetchedErrorsLength(data.errorsCount)
            //setLogs(data.currentLogs)
            //setCurrentLogSlice(data.currentLogs)
            setFetchedLogsSlice(data.currentLogs)
            setBackwardPreloadLogs(data.backLogs)
            setForwardPreloadLogs(data.forwardLogs)
            //setFetchedErrorsSlice(data.errorsSlice) //TODO GetErrors slice
            setIsFetchSuccessful(true)
            //setLoading(false)
            if (currentLogsPage === 1) {
              setFetchedFirstLogRowId(data.currentLogs[0].id)
            }
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.err(error)
            history.replace('/') // make it for unautorized/ //TODO
            history.go(0)
            reject(error)
          })
      })
    )
  }

  function handleGetLogsSliceBySession() {
    // TODO Prevent request duplicate
    return (
      new Promise((resolve, reject) => {
        logger.console(selectedFilename)
        if (selectedFilename === "") {
          logger.console("No filename selected")
          return
        }
        if (!isFetchSuccessful) {
          setLoading(checkNeedLoading())
          logger.console("Preventing Request Dupe")
          return
        }
        setIsFetchSuccessful(false)
        setLoading(!isLoadingSilent)
        var currentSessionId = -1
        if (fetchedLogsSlice[0] !== undefined) {
          currentSessionId = fetchedLogsSlice[0].sessionId
        }

        const path = serverUrl + '/v2/private/getLogsBySessionWithLimit'
        var sessionInputs = { sessionId: selectedSessionId, filename: selectedFilename, limit: PreloadPageCount * ItemsPerPage, offset: offset * ItemsPerPage }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(sessionInputs)
        })
          .then(res => {
            return handleErrorAlert(res)
          })
          .then(data => {
            logger.console('GetLogsBySessionSuccess')
            // var errors = data.errors
            // var logs = data.logs
            var sessionId = data.currentLogs[0].sessionId
            logger.console(data)
            setLogsPageCount(getTotalLength(data.logsCount))
            //setErrorsPageCount(getTotalLength(data.errorsCount))
            setFetchedLogsLength(data.logsCount)
            // setFetchedErrorsLength(data.errorsCount)
            // setLogs(data.currentLogs)
            //setCurrentLogSlice(data.currentLogs)
            setFetchedLogsSlice(data.currentLogs)
            if (currentTab === 0) {
              setBackwardPreloadLogs(data.backLogs)
              setForwardPreloadLogs(data.forwardLogs)
            }
            if (sessionId !== selectedSessionId) {
              setSelectedSessionId(sessionId)
            }
            //setFetchedErrorsSlice(data.errorsSlice) //TODO GetErrors slice
            setIsFetchSuccessful(true)
            //setLoading(false)
            if (currentLogsPage === 1) {
              setFetchedFirstLogRowId(data.currentLogs[0].id)
            }
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.err(error)
            reject(error)
            history.replace('/') // make it for unautorized/ //TODO
            history.go(0)
          })
      })
    )
  }

  function handleGetErrorsSliceBySession() {
    // TODO Prevent request duplicate
    return (
      new Promise((resolve, reject) => {
        logger.console(selectedFilename)
        if (selectedFilename === "") {
          logger.console("No filename selected")
          return
        }
        if (!isFetchSuccessful) {
          setLoading(checkNeedLoading())
          logger.console("Preventing Request Dupe")
          return
        }
        setIsFetchSuccessful(false)
        setLoading(!isLoadingSilent)
        var currentSessionId = -1
        if (fetchedLogsSlice[0] !== undefined) {
          currentSessionId = fetchedLogsSlice[0].sessionId
        }

        const path = serverUrl + '/v2/private/getErrorsBySessionWithLimit'
        var sessionInputs = { sessionId: selectedSessionId, filename: selectedFilename, limit: PreloadPageCount * ItemsPerPage, offset: offset * ItemsPerPage }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(sessionInputs)
        })
          .then(res => {
            return handleErrorAlert(res)
          })
          .then(data => {
            logger.console('GetErrorsBySessionSuccess')
            // var errors = data.errors
            // var logs = data.logs
            logger.console(data)
            setErrorsPageCount(getTotalLength(data.errorsCount))
            setFetchedErrorsLength(data.errorsCount)
            setFetchedErrorsSlice(data.currentErrors) //TODO GetErrors slice
            if (currentTab === 1) {
              setBackwardPreloadLogs(data.backErrors)
              setForwardPreloadLogs(data.forwardErrors)
            }
            setIsFetchSuccessful(true)
            //setLoading(false)
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.err(error)
            reject(error)
            history.replace('/') // make it for unautorized/ //TODO
            history.go(0)
          })
      })
    )
  }

  function handleGetErrorsSliceByDate() {

    return (
      new Promise((resolve, reject) => {
        logger.console(selectedFilename)
        if (selectedFilename === "") {
          logger.console("No filename selected")
          return
        }
        if (!isFetchSuccessful) {
          setLoading(checkNeedLoading())
          logger.console("Preventing Request Dupe")
          return
        }
        setIsFetchSuccessful(false)
        setLoading(!isLoadingSilent)

        const path = serverUrl + '/v2/private/getErrorsByDateWithLimit'
        var sessionInputs = { startDate: startDate, endDate: endDate, filename: selectedFilename, limit: PreloadPageCount * ItemsPerPage, offset: offset * ItemsPerPage }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(sessionInputs)
        })
          .then(res => {
            return handleErrorAlert(res)
          })
          .then(data => {
            logger.console('GetErrorsByDateSuccess')
            logger.console(data)
            setErrorsPageCount(getTotalLength(data.errorsCount))
            setFetchedErrorsLength(data.errorsCount)
            setFetchedErrorsSlice(data.currentErrors) //TODO GetErrors slice
            //setFetchedLogsSlice(data.currentErrors)
            setBackwardPreloadLogs(data.backErrors)
            setForwardPreloadLogs(data.forwardErrors)
            setIsFetchSuccessful(true)
            //setLoading(false)
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.err(error)
            history.replace('/') // make it for unautorized/ //TODO
            history.go(0)
            reject(error)
          })
      })
    )
  }




  function handleGetLogsBySession() {
    // let sessionId = localStorage.getItem("lastSelectedSession")
    // sessionId = parseInt(sessionId)
    // const logFilename = localStorage.getItem("lastSelectedFilename")
    // if (logFilename === null) {
    //   return
    // }

    return (
      new Promise((resolve, reject) => {
        logger.console(selectedFilename)
        if (selectedFilename === "") {
          logger.console("No filename selected")
          return
        }
        setLoading(true)
        var currentSessionId = -1
        if (fetchedLogsSlice[0] !== undefined) {
          currentSessionId = fetchedLogsSlice[0].sessionId
        }

        // if (searchedSessionId === currentSessionId && searchedFilename === selectedFilename) { //Check if needed
        //   logger.console("Searching in the SAME session")
        //   logger.console(currentSessionId)
        //   setSubmitted(false)
        //   setLoading(false)
        //   return
        // }
        // if (searchedSessionId !== 0) {
        //   logger.console("Searching in ANOTHER same session")
        //   sessionId = searchedSessionId
        // }
        //setSearchedPage(0)
        const path = serverUrl + '/v2/private/getLogsBySession'
        var sessionInputs = { sessionId: selectedSessionId, filename: selectedFilename }
        auth.secureFetch(path, {
          method: 'POST',
          body: JSON.stringify(sessionInputs)
        })
          .then(res => {
            return handleErrorAlert(res)
          })
          .then(data => {
            logger.console('GetLogsBySessionSuccess')
            var errors = data.errors
            var logs = data.logs
            setLogsPageCount(getTotalLength(data.logsCount))
            setErrorsPageCount(getTotalLength(data.errorsCount))
            //setLogs(logs)
            //setErrors(errors)
            setLoading(false)
            resolve(1)
          })
          .catch(error => {
            logger.console('Catching error')
            logger.err(error)
            reject(error)
            history.replace('/') // make it for unautorized/ //TODO
            history.go(0)
          })
      })
    )
  }

  const handleLinkClick = (event) => {
    event.preventDefault()
    setSubmitted(true)
    setLoading(true)
    var logId = event.target.dataset.logId
    logger.console(event.target)
    logger.console(logId)
    logId = parseInt(logId)
    const path = serverUrl + '/v2/private/getLogById'
    var linkInputs = { id: logId, value: selectedFilename }
    auth.secureFetch(path, {
      method: 'POST',
      body: JSON.stringify(linkInputs)
    })
      .then(res => {
        return handleErrorAlert(res)
      })
      .then(data => {
        const logId = data.log.id
        const sessionId = data.log.sessionId
        logger.console('GetLogSuccess')
        //setSelectedSession(sessions[sessionId - 1].dates)
        // setIsLogsUpdateNeeded(true)
        // setIsErrorsUpdateNeeded(true)
        setSearchedSessionId(sessionId)
        setSelectedSessionId(sessionId)
        //setSearchedFilename(selectedFilename)
        //setSearchResults([])
        setDateRange([null, null])
        setCompleteSearchRequest('')
        setFetchedSearchSlice([])
        setSelectedLogId(logId)
        setCurrentTab(0)
        //setDateRange([null, null])
  
        // setSubmitted(false)
        // setLoading(false)
        //setFilterText(res.data.log.date)
        //handleGetLogsBySession() // Check if session id is set up
      })
      .catch(error => {
        logger.console('Catching error')
        // if (error.response === undefined) {
        //   alert(error)
        //   goToHomepage()
        //   return
        // }
        logger.err(error)
        goToHomepage()
      })
  }

  const handleSearchSubmitWithLimit = (event) => {
    setSubmitted(true)
    setLoading(true)
    event.preventDefault()
    const path = serverUrl + '/v2/private/findLogsWithLimit'
    const searchInput = { text: searchText, filename: selectedFilename, limit: PreloadPageCount * ItemsPerPage, offset: offset * ItemsPerPage }
    auth.secureFetch(path, {
      method: 'POST',
      body: JSON.stringify(searchInput)
    })
      .then(res => {
        return handleErrorAlert(res)
      })
      .then(data => {
        logger.console('SearchSuccess')
        setLogsPageCount(getTotalLength(data.logsCount))
        setFetchedSearchLength(data.logsCount)
        //setCurrentLogSlice(data.currentLogs)
        setFetchedSearchSlice(data.currentLogs)
        setBackwardPreloadLogs(data.backLogs)
        setForwardPreloadLogs(data.forwardLogs)
        setCompleteSearchRequest(searchText)
        setSubmitted(false)
        setLoading(false)
      })
      .catch(error => {
        logger.console('Catching error')
        logger.console(error.response)
        if (error.response !== undefined && error.response.status === 401) {
          cookies.remove("rToken", { path: '/' })
        }
        logger.err(error)
        goToHomepage()
      })
  }

  const handleSearchSubmit = (event) => {
    setSubmitted(true)
    setLoading(true)
    event.preventDefault()
    const path = serverUrl + '/v2/private/findLogs'
    const searchInput = { text: searchText, filename: selectedFilename };
    auth.secureFetch(path, {
      method: 'POST',
      body: JSON.stringify(searchInput)
    })
      .then(res => {
        return handleErrorAlert(res)
      })
      .then(data => {
        logger.console('SearchSuccess')
        var logs = data.logs
        var lenLogs = logs.length
        setSearchPageCount(getTotalLength(lenLogs))
        //setSearchResults(logs)
        setSubmitted(false)
        setLoading(false)
      })
      .catch(error => {
        logger.console('Catching error')
        logger.console(error.response)
        if (error.response !== undefined && error.response.status === 401) {
          cookies.remove("rToken", { path: '/' })
        }
        logger.err(error)
        goToHomepage()
      })
  }

  const handleRegisterSubmit = event => {
    const password = pass;
    const confirmPassword = confirmPass;
    event.preventDefault();

    if (password !== confirmPassword) {
      setPasswordMatch(false)
    }
    else {
      setPasswordMatch(true)
      const input = { login: user, password: pass }//, fingerprint: fingerprint }
      SendRegisterRequest(input)
    }
  }

  const countOffsetJump = (page) => {
    let i = 0
    i = (page / 10) | 0
    logger.console("OFFSET: " + i)
    logger.console("ROUND OFFSET: " + i * PreloadPageCount)
    return i * PreloadPageCount
  }

  const isNumOutOfSequence = (num, start, end) => {
    if (num < start || num > end) {
      return true
    }
    return false
  }

  const isNumBetweenSequence = (num, start, end) => {
    if (num >= start && num <= end) {
      return true
    }
    return false
  }

  const handleLogsPageChange = (page) => {
    logger.console("PAGE CHANGED? OFFSET " + offset)
    let dif = page - offset
    logger.console("COUNT currentLogsPage - offset = " + dif)
    logger.console("REQUEST Pending " + !isFetchSuccessful)
    if (!isFetchSuccessful && isNumOutOfSequence(dif, 1, PreloadPageCount)) {
      logger.console("PAGE NOT CHanged")
      setLoading(true)
      return
    }
    logger.console("Setting State Logs Page")
    //setIsLogsUpdateNeeded(true)
    setCurrentLogsPage(page)
    setCurrentSlicePage(getCurrentSliceIndex(page))
    handleOffset(page)
  }

  // if (page > 5 + offset && page < 10 + offset) {
  //   setOffset(offset + 5)
  // }

  // if (page < offset && offset - page < 5) {
  //   setOffset(offset - 5)
  // }

  // if (page < offset && offset - page > 5 && page < 5) {
  //   setOffset(page - 5)
  // }

  // if (page > 10 + offset) {
  //   setOffset(page - 5)
  // }

  //}

  const handleSetOffsetSilent = (dif) => {

    logger.console("Setting Offset ")
    logger.console("dif=" + dif)
    if (dif > 0) {
      logger.console("Offset Forward")
      setOffset(dif)
      return
    }
    if (dif <= -10) {
      logger.console("Offset Backward")
      setOffset(-dif)
    }
    if (dif > -10 || dif === 0) {
      logger.console("Offset 0")
      setOffset(0)
    }

  }

  const countSliceOffset = (page) => {
    let n = page % 10
    if (n === 0) {
      n = 10
    }
    return n
  }

  const getCurrentSliceIndex = (page) => {
    return countSliceOffset(page) - 1
  }

  const handleErrorsPageChange = (page) => {
    logger.console("ERR PAGE CHANGED? OFFSET " + offset)
    let dif = page - offset
    logger.console("COUNT currentLogsPage - offset = " + dif)
    logger.console("REQUEST Pending " + !isFetchSuccessful)
    if (!isFetchSuccessful && isNumOutOfSequence(dif, 1, PreloadPageCount)) {
      logger.console("PAGE NOT CHanged")
      setLoading(true)
      return
    }
    logger.console("Setting State Errors Page")
    //setIsErrorsUpdateNeeded(true)
    setCurrentErrorsPage(page)
    setCurrentSlicePage(getCurrentSliceIndex(page))
    handleOffset(page)
  }
  //   logger.console("Setting State Errors Page")
  //   setCurrentErrorsPage(page)
  //   setCurrentSlicePage(getCurrentSliceIndex(page))
  //   handleOffset(page)
  // }

  const handleSearchPageChange = (page) => {
    logger.console("Setting State Search Page")
    setCurrentSearchPage(page)
    setCurrentSlicePage(getCurrentSliceIndex(page))
    handleOffset(page)
  }

  const getLogsSessions = () => {
    if (isAuthorized && selectedFilename !== "") {
      logger.console('GettingLogsSessions')
      localStorage.setItem("lastSelectedFilename", selectedFilename)
      sessionStorage.removeItem("lastSelectedSessionId")
      //resetPagination()
      handleGetLogsSessions()
    }
  }

  function toggleUpdate() {
    logger.console("Offset changed!")
    if (currentTab === 0 && offset !== logsOffset) {
      logger.console("SETTING LOGS UPDATE NEEDED")
      setLogsOffset(offset)
      setIsLogsUpdateNeeded(true)
    }
    if (currentTab === 1 && offset !== errorsOffset) {
      logger.console("SETTING ERRORS UPDATE NEEDED")
      setErrorsOffset(offset)
      setIsErrorsUpdateNeeded(true)
    }
  }

  // useEffect(() => {
  //   let pathname = window.location.hash
  //   const checkValues = ["#/login", "#/register", "#/logout"] 
  //   logger.console(pathname)

  //   if (checkValues.includes(pathname)) {
  //     logger.console("Auth NOT needed!")
  //     setIsAuthNeeded(false)
  //   } else {
  //     logger.console("Auth needed!")
  //     setIsAuthNeeded(true)
  //   }
  // })

  useEffect(() => {
    let isLoggedIn = sessionStorage.getItem("isAuthenticated")
    if (isLoggedIn === null) {
      isLoggedIn = localStorage.getItem("isAuthenticated")
      setCanRememberMe(isLoggedIn !== null)
    }

    if (isLoggedIn) {
      logger.console("Setting isAuth")
      setIsAuthenticated(true)
      return
    }
    logger.console("Login page from Effect App")
    goToLoginPage()
  }, [])


  // useEffect(() => {
  //   // let isLoggedIn = localStorage.getItem("isAuthenticated")
  //   logger.console("Try Auth Hook")
  //   if (isAuthenticated && isFingerprintReady) {
  //     logger.console('AuthAttempt hook')
  //     handleAuthSubmit()
  //   }
  //   // else {
  //   //   logger.console("Login page from App")
  //   // //history.push("/login")
  //   // goToLoginPage()
  //   // }
  //   // if () {
  //   //   history.replace('/')
  //   //   history.go(0)
  //   // }
  // }, [isAuthenticated, isFingerprintReady]); // used for setting hook only on Mount event; [authToken] listens to variable change



  useEffect(() => {
    if (startDate === null && endDate === null && selectedFilename !== '') {
      resetPagination()
      setIsErrorsUpdateNeeded(true)
      setIsLogsUpdateNeeded(true)
    }
  }, [startDate, endDate])

  function handleOffset(currentPage) {
    var dif = currentPage - offset
    logger.console("PAge changed " + currentPage)
    logger.console("dif " + (dif))
    //logger.console(dif + PreloadPageCount <= -PreloadPageCount)
    if (dif > 0 && dif > PreloadPageCount && dif - PreloadPageCount < PreloadPageCount) { // Set await for logs update (new state?)
      logger.console("Forward silent")
      setCurrentLogSlice(forwardPreloadLogs)
      setIsGoingForward(true)
      setIsGoingBackwards(false)
      setIsLoadingSilent(true)
      handleSetOffsetSilent(PreloadPageCount + offset)
      return
    }

    if (dif <= 0 && dif >= -PreloadPageCount && dif + PreloadPageCount <= PreloadPageCount) {
      logger.console("Backward silent")
      setCurrentLogSlice(backwardPreloadLogs)
      setIsGoingForward(false)
      setIsGoingBackwards(true)
      setIsLoadingSilent(true)
      handleSetOffsetSilent(PreloadPageCount - offset)
      return
    }

    if (dif > 0 && dif > PreloadPageCount && dif - PreloadPageCount >= PreloadPageCount) {
      logger.console("Forward LOUD")
      setIsGoingForward(false)
      setIsGoingBackwards(false)
      setIsLoadingSilent(false)
      handleSetOffsetSilent(countOffsetJump(currentPage))
      return
    }

    if (dif < 0 && dif < -PreloadPageCount && dif + PreloadPageCount <= 0) {
      logger.console("Backward LOUD")
      setIsGoingForward(false)
      setIsGoingBackwards(false)
      setIsLoadingSilent(false)
      handleSetOffsetSilent(countOffsetJump(currentPage))//- offset)
      return
    }
  }

  useEffect(() => {
    if (currentTab === 1 && fetchedErrorsSlice[0] !== undefined) {
      logger.console("Showing Errors")
      setCurrentLogSlice(fetchedErrorsSlice)
      setCurrentSlicePage(getCurrentSliceIndex(currentErrorsPage))
      setOffset(errorsOffset)
    }
  }, [currentTab, fetchedErrorsSlice])

  useEffect(() => {
    if (completeSearchRequest !== '' && fetchedSearchSlice[0] !== undefined) {
      setCurrentLogSlice(fetchedSearchSlice)
      setCurrentSlicePage(getCurrentSliceIndex(currentSearchPage))
      setOffset(countOffsetJump(currentSearchPage))
    }
  }, [currentTab, completeSearchRequest, fetchedSearchSlice])

  useEffect(() => {
    if (currentTab === 0 && fetchedLogsSlice[0] !== undefined && completeSearchRequest === '' ) {
      logger.console("Setting Logs Offset")
      setCurrentLogSlice(fetchedLogsSlice)
      setCurrentSlicePage(getCurrentSliceIndex(currentLogsPage))
      setOffset(logsOffset)
    }
  }, [currentTab, fetchedLogsSlice])

  useEffect(() => {
    if (searchedPage > 0) { //&& currentTab === 0 && completeSearchRequest === '') {
      setCurrentLogSlice(fetchedLogsSlice)
      setCurrentLogsPage(searchedPage)
      setCurrentSlicePage(getCurrentSliceIndex(searchedPage))
      handleOffset(searchedPage)
      setSearchedPage(0)
      logger.console("Set searched Page: " + searchedPage)
    }
  }, [searchedPage])

  // useEffect(() => {
  //   if (searchedPage !== 0) {
  //     setCurrentLogsPage(searchedPage)
  //     handleOffset(searchedPage)
  //     logger.console("Set searched Page: " + searchedPage)
  //   }
  // }, [searchedPage])

  useEffect(() => {
    logger.console("offset changed?")
    if (offset >= 0 && isAuthorized && selectedFilename !== '') {
      toggleUpdate()
    }
  }, [offset, isAuthorized, selectedFilename])

  useEffect(() => {
    if (searchText === '') {
      setCompleteSearchRequest('')
    }
  }, [searchText])

  // useEffect(() => {
  //   if (offset >= 0 && filenames[0] !== undefined) {
  //     setLoading(!isFetchSuccessful)
  //   }
  // }, [offset])

  useEffect(() => {
    let filename = localStorage.getItem("lastSelectedFilename")
    let sessionId = sessionStorage.getItem("lastSelectedSessionId")
    sessionId = parseInt(sessionId)

    logger.console("Selected SessionID is:")
    logger.console(sessionId)
    logger.console(filename)

    if (filename === null) {
      filename = ""
    }
    if (isNaN(sessionId)) {
      sessionId = 0
    }

    setSelectedFilename(filename)
    setSelectedSessionId(sessionId)
  }, [])

  // useEffect(() => {
  //   if(selectedSessionId === 0 && sessions[0] !== undefined) {
  //     setSelectedSessionId(sessions[sessions.length - 1].id)
  //   }
  // },[])

  const usePathname = () => {
    const location = useLocation();
    return location.pathname;
  }


  //let pathname = usePathname()

  useEffect(() => {
    logger.console("loadFingerPrint")
    let promise = auth.loadFingerprint()

    promise.then(
      res => setIsFingerprintReady(res)
    )
  }, [])

  useEffect(() => {
    logger.console("check token hook")
    setIsAuthorized(auth.checkToken())
  }, [])




  useEffect(() => {
    if (isAuthorized) {
      logger.console("Getting Logs Filenames")
      handleGetLogsFilenames()
    }
  }, [isAuthorized])

  // useEffect(() => {
  //   async function fetchLogs() {
  //     if (isAuthorized && selectedFilename !== '' && startDate === null && endDate === null) { // && logs[0] === undefined && selectedSessionId !== 0
  //       logger.console('GettingLogsBySession')
  //       await handleGetLogsBySession()
  //         .then(res => {
  //           logger.console(res + " GOT LOGS")
  //           getLogsSessions()
  //         })
  //     }
  //   }
  //   fetchLogs()
  // }, [selectedFilename, isAuthorized, startDate, endDate])

  // useEffect(() => {
  //   async function fetchLogs() {
  //     if (isAuthorized && startDate !== null && endDate !== null) { // && logs[0] === undefined && selectedSessionId !== 0
  //       logger.console('GettingLogsByDate')
  //       //setSelectedSession("")
  //       if (isErrorsUpdateNeeded) {
  //         // await handleGetLogsSliceByDate()
  //         //   .then(res => {
  //         //     logger.console(res + " GOT LOGS")
  //         //     setIsLogsUpdateNeeded(false)
  //         //   })
  //         await handleGetErrorsSliceByDate()
  //           .then(res => {
  //             logger.console(res + " GOT Errors")
  //             setIsErrorsUpdateNeeded(false)
  //           })
  //       } else {
  //         logger.console("NO ErrorsUpdate NEEDED")
  //       }
  //       if (isLogsUpdateNeeded) {
  //         // await handleGetErrorsSliceByDate()
  //         //   .then(res => {
  //         //     logger.console(res + " GOT Errors")
  //         //     setIsErrorsUpdateNeeded(false)
  //         //   })
  //         await handleGetLogsSliceByDate()
  //           .then(res => {
  //             logger.console(res + " GOT LOGS")
  //             setIsLogsUpdateNeeded(false)
  //             getLogsSessions()
  //           })
  //       } else {
  //         logger.console("NO LogsUpdate NEEDED")
  //       }
  //     }
  //   }
  //   fetchLogs()
  // }, [isAuthorized, startDate, endDate])

  useEffect(() => {
    async function fetchLogs() {
      if (isAuthorized && startDate !== null && endDate !== null && isErrorsUpdateNeeded) { // && logs[0] === undefined && selectedSessionId !== 0
        logger.console('GettingLogsByDate')
          await handleGetErrorsSliceByDate()
            .then(res => {
              logger.console(res + " GOT Errors")
              setIsErrorsUpdateNeeded(false)
            })       
      }
    }
    fetchLogs()
  }, [isAuthorized, startDate, endDate, isErrorsUpdateNeeded])

  useEffect(() => {
    async function fetchLogs() {
      if (isAuthorized && startDate !== null && endDate !== null && isLogsUpdateNeeded) { // && logs[0] === undefined && selectedSessionId !== 0
        logger.console('GettingLogsByDate')
          await handleGetLogsSliceByDate()
            .then(res => {
              logger.console(res + " GOT LOGS")
              setIsLogsUpdateNeeded(false)
              getLogsSessions()
            })
      }
    }
    fetchLogs()
  }, [isAuthorized, startDate, endDate, isLogsUpdateNeeded])

  useEffect(() => {
    async function fetchLogs() {
      if (isAuthorized && startDate === null && endDate === null && isErrorsUpdateNeeded) { // && logs[0] === undefined && selectedSessionId !== 0
        logger.console('GettingErrorsSLICEBySession')
          // await handleGetLogsSliceBySession()
          //   .then(res => {
          //     logger.console(res + " GOT LOGS")
          //     setIsLogsUpdateNeeded(false)
          //   })
          await handleGetErrorsSliceBySession()
            .then(res => {
              logger.console(res + " GOT Errors")
              setIsErrorsUpdateNeeded(false)
            })

      //TODO Make Errors Request
    }
  }
    logger.console("Try fetch Errors")
    fetchLogs()
  }, [isAuthorized, startDate, endDate, isErrorsUpdateNeeded])


  useEffect(() => {
    async function fetchLogs() {
      if (isAuthorized && startDate === null && endDate === null && isLogsUpdateNeeded) { // && logs[0] === undefined && selectedSessionId !== 0
        logger.console('GettingLogsSLICEBySession')
                  // await handleGetErrorsSliceBySession()
          //   .then(res => {
          //     logger.console(res + " GOT Errors")
          //     setIsErrorsUpdateNeeded(false)
          //   })
          await handleGetLogsSliceBySession()
            .then(res => {
              logger.console(res + " GOT LOGS")
              setIsLogsUpdateNeeded(false)
              getLogsSessions()
            })
      }
    }
    logger.console("Try fetch Logs")
    fetchLogs()
  }, [isAuthorized, startDate, endDate, isLogsUpdateNeeded])

  
  // useEffect(() => {
  //   async function fetchLogs() {
  //     if (isAuthorized && offset >= 0 && startDate === null && endDate === null) { // && logs[0] === undefined && selectedSessionId !== 0
  //       logger.console('GettingLogsSLICEBySession')
  //       if (isErrorsUpdateNeeded) {
  //         // await handleGetLogsSliceBySession()
  //         //   .then(res => {
  //         //     logger.console(res + " GOT LOGS")
  //         //     setIsLogsUpdateNeeded(false)
  //         //   })
  //         await handleGetErrorsSliceBySession()
  //           .then(res => {
  //             logger.console(res + " GOT Errors")
  //             setIsErrorsUpdateNeeded(false)
  //           })
  //       } else {
  //         logger.console("NO ErrorsUpdate NEEDED")
  //       }
  //       if (isLogsUpdateNeeded) {
  //         // await handleGetErrorsSliceBySession()
  //         //   .then(res => {
  //         //     logger.console(res + " GOT Errors")
  //         //     setIsErrorsUpdateNeeded(false)
  //         //   })
  //         await handleGetLogsSliceBySession()
  //           .then(res => {
  //             logger.console(res + " GOT LOGS")
  //             setIsLogsUpdateNeeded(false)
  //             getLogsSessions()
  //           })
  //       } else {
  //         logger.console("NO LogsUpdate NEEDED")
  //       }
  //     }

  //     //TODO Make Errors Request
  //   }
  //   fetchLogs()
  // }, [selectedFilename, offset, isAuthorized, startDate, endDate])

  // useEffect(() => {
  //   if (currentTab === 0)
  // })

  // useEffect(() => {
  //   async function fetchLogs() {
  //     if (isAuthorized && selectedFilename !== '' && offset >= 0 && startDate !== null && endDate !== null) { // && logs[0] === undefined && selectedSessionId !== 0
  //       logger.console('GettingLogsSliceByDate')
  //       setSelectedSession("")
  //       await handleGetLogsSliceByDate()
  //         .then(res => {
  //           logger.console(res + " GOT LOGS")
  //           getLogsSessions()
  //         })
  //     }
  //   }
  //   fetchLogs()
  // }, [selectedFilename, offset, isAuthorized, startDate, endDate])

  // useEffect(() => {
  //   if (isAuthorized && selectedLogId !== 0) {
  //     logger.console('GettingLogsBySearchRequest')
  //     handleGetLogsBySession()
  //   }
  // }, [selectedLogId, isAuthorized])

  useEffect(() => {
    async function authorize() {
      if (isAuthenticated && !isAuthorized && isFingerprintReady) {
        logger.console("AuthPromiseHook")
        //var prom = auth.authPromise(serverUrl + '/v2/authAttempt')
        //await prom.then(res => {
        await auth.authPromise(serverUrl + '/v2/authAttempt')
          .then(res => {
            setIsAuthorized(true)
          })
          .catch(error => {
            alert(error)
            goToLoginPage()
          })
      }
      return
    }
    authorize()
  }, [isAuthorized, isAuthenticated, isFingerprintReady])

useEffect(() => {
  if (selectedSessionId >= 0 && selectedFilename !== '' && selectedLogId === 0) {
    logger.console("Setting IS Logs UPDATE needed")
    setIsLogsUpdateNeeded(true)
    setIsErrorsUpdateNeeded(true)
    resetPagination()
  }
}, [selectedFilename, selectedSessionId])

useEffect(() => {
  if (isFetchSuccessful && !isLogsUpdateNeeded && !isErrorsUpdateNeeded) {
    setLoading(false)
  }
}, [isFetchSuccessful, isLogsUpdateNeeded, isErrorsUpdateNeeded])


  useEffect(() => {
    if (isAuthorized && selectedSession !== "" && sessions[0] !== undefined) { // && selectedSessionId !== 0
      logger.console('GettingLogsServices')
      // localStorage.setItem("lastSelectedFilename", selectedFilename)
      // sessionStorage.removeItem("lastSelectedSessionId")
      // resetPagination()
      handleGetLogServiceInfo()
    }
  }, [selectedSession, isAuthorized, sessions])

  useEffect(() => {
    var sessionId = sessionStorage.getItem("lastSelectedSessionId")
    logger.console("Reading SessionId fromSession " + sessionId)
    sessionId = parseInt(sessionId)
    if (isNaN(sessionId)) {
      return
    }
    setSelectedSessionId(sessionId)

    if (sessions[0] !== undefined && selectedSessionId > 0) {
      logger.console("Setting Session Name")
      setSelectedSession(sessions[selectedSessionId - 1].value)
    }
  }, [])

  useEffect(() => {
    if (sessions[0] !== undefined && selectedSessionId > 0) {
      logger.console("Setting Session Name " + selectedSessionId)
      setSelectedSession(sessions[selectedSessionId - 1].value)
    }
  }, [sessions, selectedSessionId])

  useEffect(() => {
    if (searchedSessionId !== 0) {
      sessionStorage.setItem("lastSelectedSessionId", searchedSessionId)
    }
  }, [searchedSessionId])

  useEffect(() => {
    if (selectedSessionId !== 0) {
      sessionStorage.setItem("lastSelectedSessionId", selectedSessionId)
    }
  }, [selectedSessionId])



  // useEffect(() => {
  //   logger.console("Try getting search page")
  //   if (logs[0] !== undefined && selectedLogId !== 0) { //|| selectedLogId !== 0 && searchedSessionIdId === logs[0].sessionId) {
  //     if (logs[0].sessionId !== searchedSessionId) {
  //       return
  //     }
  //     logger.console("Getting Search Page")
  //     var slice = logs.slice(0, selectedLogId - logs[0].id)
  //     var lenSlice = slice.length
  //     getSearchedPage(lenSlice)
  //   }
  // }, [selectedLogId, logs])

  useEffect(() => {
    if (selectedSessionId > 0 && selectedLogId >= fetchedFirstLogRowId) { //|| selectedLogId !== 0 && searchedSessionIdId === logs[0].sessionId) {
      logger.console("Try getting search page")
      if (selectedSessionId !== searchedSessionId) {
        setSelectedLogId(0)
        return
      }
      logger.console("Getting Search Page")
      getSearchedPage()
      // setIsLogsUpdateNeeded(true)
      // setIsErrorsUpdateNeeded(true)
    }
  }, [selectedLogId, fetchedFirstLogRowId, selectedSessionId])

  function getSearchedPage() {
    logger.console("Calculating slice")
    // if (len <= 0) {
    //   setLoading(false)
    //   setSubmitted(false)
    //   return
    // }
    logger.console("LOG ID" + fetchedLogsSlice[0].id)
    let i = countSearchedPage(selectedLogId)
    logger.console("Calculated page: " + i)
    setSearchedPage(i)
    // setCurrentLogsPage(i)
    // setCurrentSearchPage(i)
    setLoading(false)
    setSubmitted(false)
  }

  const countSearchedPage = (logId) => {
    let searchedPage = Math.trunc((logId - fetchedFirstLogRowId) / ItemsPerPage)
    logger.console("TRUNCATE: " + searchedPage)
    searchedPage = searchedPage + 1
    return searchedPage
  }



  // useEffect(() => {

  //   if (filterText !== "" || selectedService !== "") {
  //     var logRows = []
  //     logs.forEach((log) => {
  //       const logDate = log.dateUtc
  //       const logDesc = log.description.toLowerCase()
  //       if (logDate !== filterText && logDesc.indexOf(filterText) === -1 || log.service.indexOf(selectedService) === -1) {
  //         return;
  //       }
  //       logRows.push(log)
  //     });
  //     var errorRows = []
  //     errors.forEach((error) => {
  //       const errorDate = error.dateUtc
  //       const errorDesc = error.description.toLowerCase()
  //       if (errorDate !== filterText && errorDesc.indexOf(filterText) === -1 || error.service.indexOf(selectedService) === -1) {
  //         return;
  //       }
  //       errorRows.push(error)
  //     });
  //     var searchRows = []
  //     searchResults.forEach((searchResult) => {
  //       const searchDate = searchResult.dateUtc
  //       const searchDesc = searchResult.description.toLowerCase()
  //       if (searchDate !== filterText && searchDesc.indexOf(filterText) === -1 || searchResult.service.indexOf(selectedService) === -1) {
  //         return;
  //       }
  //       searchRows.push(searchResult)
  //     });
  //     var lenLogs = logRows.length
  //     var lenErrs = errorRows.length
  //     var lenSearch = searchRows.length
  //     setLogsPageCount(getTotalLength(fetchedLogsLength))
  //     setErrorsPageCount(getTotalLength(fetchedErrorsLength))
  //     setSearchPageCount(getTotalLength(fetchedSearchLength))
  //     setFilteredLogs(logRows)
  //     setFilteredErrors(errorRows)
  //     setFilteredSearch(searchRows)
  //     //setCurrentPage(1)
  //     logger.console("filter complete")

  //   }
  // }, [filterText, selectedService, logs, errors, searchResults])

  // useEffect(() => {
  //   if (authToken !== "" && logs[0] !== undefined) {
  //     if (selectedSessionId !== undefined){
  //       logger.console('CalculatingLogsLength')
  //       setLogsPageCount(getTotalLength())
  //     }
  //   }
  // }, [selectedSessionId])

  return (<HashRouter history={history}>
    <div className="App">
      <AuthContext.Provider value={{ isAuthenticated: isAuthenticated }}>
        <Navbar searchText={searchText} handleSearchInput={handleSearchInput}
          handleSearchSubmit={handleSearchSubmitWithLimit}// handleUnblockSubmit={handleUnblockSubmit}
          handleFilterTextInput={handleFilterTextInput} filterText={filterText}
          handleSelectServiceChange={handleSelectServiceChange} services={services}
          value={selectedService} />
        <ScrollButton />
        <Switch>
          <PrivateRoute exact path='/'
            data={
              {
                //errors: errors, searchResults: searchResults, 
                logs: currentLogSlice, submitted: submitted, logsDropdownValue: selectedFilename,
                filteredLogs: filteredLogs, filteredErrors: filteredErrors, filteredSearch: filteredSearch, filenames: filenames,
                sessions: sessions, sessionsDropdownValue: selectedSession, filterText: filterText, filterService: selectedService,
                loading: loading, logsPageCount: logsPageCount, errorsPageCount: errorsPageCount, searchPageCount: searchPageCount,
                currentTab: currentTab, itemsPerPage: ItemsPerPage, logId: selectedLogId, searchedPage: searchedPage, currentSlicePage: currentSlicePage,
                currentLogsPage: currentLogsPage, currentErrorsPage: currentErrorsPage, currentSearchPage: currentSearchPage,
                datePicker: datePicker, needScrollDown: needScrollDown, dateRange: dateRange, completeSearchRequest: completeSearchRequest,
                handleLogsPageChange: handleLogsPageChange, handleErrorsPageChange: handleErrorsPageChange, handleSearchPageChange: handleSearchPageChange,
                handleSelectSessionChange: handleSelectSessionChange, handleSelectFilenameChange: handleSelectFilenameChange,
                handleTabClick: handleTabClick, handleLinkClick: handleLinkClick
              }
            }
            component={Board} />
          <Route exact path="/login"
            render={() =>
              <LoginContainer user={user} pass={pass} handleUsernameInput={handleUsernameInput}
                handlePassInput={handlePassInput} handleLoginSubmit={handleLoginSubmit}
                handleRememberMeCheckbox={handleRememberMeCheckbox}
                canRememberMe={canRememberMe} submitted={submitted} />
            } />
          <Route exact path="/logout"
            render={() =>
              <Logout handleLogoutSubmit={handleLogoutSubmit}
                submitted={submitted} isAuthenticated={isAuthenticated} />
            } />
          <Route exact path="/registration"
            render={() =>
              <Register user={user} pass={pass}
                confirmPass={confirmPass} passwordMatch={passwordMatch}
                handleUsernameInput={handleUsernameInput} handlePassInput={handlePassInput}
                handleConfirmPassInput={handleConfirmPassInput} handleRegisterSubmit={handleRegisterSubmit} />
            } />
          {/* <Route exact path='/board'
            render={() => <BoardContainer handleSearchInput={handleSearchInput} handleSearchSubmit={handleAuthSubmit}
            searchText={searchText} 
            />} /> */}
        </Switch>
      </AuthContext.Provider>
    </div>
  </HashRouter>
  );
}



const getParams = (pathname) => {
  const matchProfile = matchPath(pathname, {
    path: `/login`,
  });
  logger.console(matchProfile)
  logger.console(matchProfile.params)
  return (matchProfile && matchProfile.params) || {};
};



// // if (process.env.NODE_ENV !== 'production') {
// //   logger.add(new transports.Console({
// //     format: format.combine(
// //       format.colorize(),
// //       format.simple()
// //     )
// //   }));
// // }




function UsePathname() {

  const [pathname, setPathname] = useState("")

  let location = useLocation()
  logger.console(location.pathname)
  // let pathname = ""

  useEffect(() => {
    if (location.pathname === "/login") {
      setPathname("Login")
    }
    if (location.pathname === "/logout") {
      setPathname("Logout")
    }
    if (location.pathname === "/registration") {
      setPathname("Registration")
    }
    if (location.pathname === "/") {
      setPathname("Root")
    }
  }, [])
  return pathname
}







