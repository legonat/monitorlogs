import * as Sentry from '@sentry/react';


export default function Logger() {
    function err(err) {
      if (process.env.NODE_ENV !== 'production') {
        console.log(err)
      }
      if (err.response !== undefined) {
          alert(err.response.status + ': ' + err.response.data.error)
      }
      return Sentry.captureException(err)
    }
    function info(msg) {
      if (process.env.NODE_ENV !== 'production') {
        console.log(msg)
      }
      return Sentry.captureMessage(msg)
    }
    function local(msg) {
      if (process.env.NODE_ENV === 'production') {
        return
      }
      console.log(msg)
      return
    }
    return {
      err: err,
      info: info,
      console: local
    }
  }