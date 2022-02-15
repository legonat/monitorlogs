import FingerprintJS from '@fingerprintjs/fingerprintjs';
import Logger from './logger';

const logger = new Logger()

export default function AuthModule() {
  const fetch = window.fetch
  const authOrigins = ["http://127.0.0.1:5000", "http://localhost:5000"];
  const serverUrl = authOrigins[1]
  let token = '';
  let fingerprint = '';
  let isLoadingComplete = false;
  let isAuthPending = false;
  let resStream = new ReadableStream()

  async function loadFingerprint() {

    const fp = await FingerprintJS.load();

    const result = await fp.get();

    fingerprint = result.visitorId;
    logger.console(fingerprint);
    isLoadingComplete = true


    return isLoadingComplete
  }

  function getLoadingStatus() {
    logger.console("Fingerprint Loaded " + isLoadingComplete)
    return isLoadingComplete
  }

  function checkAuthPending() {
    logger.console("Auth Pending " + isAuthPending)
    return isAuthPending
  }


  function setToken(value) {
    token = value
  }

  function setAuthPending(value) {
    logger.console("Setting auth pending " + value)
    isAuthPending = value
  }

  function checkToken(token) {
    logger.console("token " + token)
    if (token === undefined || token === '') {
      logger.console("Empty Token")
      return false
    }
    let parsedToken = parseJwt(token)
    if (parsedToken === undefined) {
      logger.console("Parsing failed")
      return false
    }
    // logger.console(parsedToken)
    // logger.console(parsedToken.exp)
    // logger.console(Math.floor(Date.now()/1000))
    if (parsedToken.exp <= Math.floor(Date.now() / 1000)) {
      logger.console("Token expired")
      return false
    }
    logger.console("Token Is Valid")
    return true
  }

  function delay(ms) {
    return new Promise((resolve, reject) => {
      setTimeout(resolve, ms);
    });
  }

  async function authPromise(path) {
    if (isAuthPending) {
      logger.console("Auth Pending in Promise")
      await delay(5000).then(() =>
        logger.console("Auth Timout Finished")
      )
    }
    return new Promise((resolve, reject) => {
      // if (isAuthPending) {
      //   logger.console("Auth Pending in Promise")
      //   reject("Auth failed")
      // }
      try {
        if (!isAuthPending) {
          logger.console("Start Auth Promise")
          setAuthPending(true)
          isAuthPending = true
          secureAuth(path)
            .then(res => {
              logger.console(res)
              if (res.status === 200) {
                logger.console("Auth Promise")
                resStream = res
                return res.json()
              }
              if (res.status === 205) {
                logger.console("ReAuth Prevented")
                return new Promise((resolve, reject) => {
                  resolve(token)
                })
              }
              if (res.status === 401) {
                logger.console("401: Unauthorized")
                alert("401: Unauthorized")
                localStorage.clear()
                sessionStorage.clear()
                window.location.reload()
                return
              }
              if (!res.ok) {
                throw Error(res.status + ": " + res.statusText)
              }
            })
            .then(data => {
              //newToken = data.accessToken
              if (data === token) {
                logger.console("No Reauth Success")
                setToken(data)
                setAuthPending(false)
                resolve(true)
                return
              }
              setToken(data.accessToken)
              setAuthPending(false)
              logger.console("propmise authSuccess")
              resolve(true)
              //return new Promise.resolve("done")
              // setIsAuthorized(true) //TODO make zamykanie
            })
            .catch(error => {
              logger.console('Catching error')
              logger.console(error)
              logger.err(error)
              reject(error)
            })
        }
      }
      catch (e) {
        logger.console(e)
      }
      
      //return promise

    })
  }

  async function secureFetch(resource, options) {
    //TODO: To bypass auth race make functions with Promisification
    var newToken = ""
    //var done = await new Promise.resolve("done")
    logger.console("secure fetch")
    let isAuthorized = checkToken(token)
    if (!isAuthorized) {
      logger.console("try Auth Fetch")
      //const path = serverUrl + '/v2/authAttempt'
      await authPromise(serverUrl + '/v2/authAttempt')
        .then(res => {
          logger.console("Reauth success " + res)
        })
        .catch(error => {
          alert(error)
        })

    }
    //TODO Place Auth Here or Check Token Validity only
    logger.console("building fetch req")
    let req = new Request(resource, options)
    let destOrigin = new URL(req.url).origin
    // if (newToken !== '') {
    //   token = newToken
    // }
    if (authOrigins.includes(destOrigin)) {
      req.headers.set('Authorization', `Bearer ${token}`)
      req.headers.set('Content-Type', 'application/json')
    } else {
      logger.console("Fetch failed. No auth token")
      return
    }
    if (process.env.NODE_ENV !== 'production') {
      return fetch(req, {
        credentials: 'include'
      })
    }
    logger.console("sec fetch ready")
    return fetch(req)
  }

  function secureAuth(resource, options) {
    try {
      let isAuthorized = checkToken(token)
    if (!isAuthorized) {
      logger.console("secAuth attempt")
      let req = new Request(resource, options)
      let destOrigin = new URL(req.url).origin
      logger.console(req.url)
      if (fingerprint !== '' && authOrigins.includes(destOrigin)) {
        req.headers.set('Content-Type', 'application/json')
        req.headers.set('Fingerprint', fingerprint)
      }
      if (process.env.NODE_ENV !== 'production') {
        return fetch(req, {
          credentials: 'include'
        })
      }
      logger.console("auth request ready")
      return fetch(req)
    }
    // const obj = {hello: 'world'};
    // const blob = new Blob([JSON.stringify(obj, null, 2)], {type : 'application/json'});
    var init = {
      status: 205
    }
    var response = new Response(null, init)
      return new Promise((resolve, reject) => {
        resolve (response)
      })
  }
  catch (e) {
    logger.console(e)
  }
  }

  function parseJwt(token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
  };
  return {
    loadFingerprint: loadFingerprint,
    getLoadingStatus: getLoadingStatus,
    checkAuthPending: checkAuthPending,
    checkToken: checkToken,
    authPromise: authPromise,
    setToken: setToken,
    secureFetch: secureFetch,
    secureAuth: secureAuth
  }
}